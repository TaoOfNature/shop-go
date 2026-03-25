package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dawnstack/shop-go/internal/cache"
	"github.com/dawnstack/shop-go/internal/config"
	"github.com/dawnstack/shop-go/internal/model"
	"github.com/dawnstack/shop-go/internal/mq"
	"github.com/dawnstack/shop-go/internal/pkg/snowflake"
	"github.com/dawnstack/shop-go/internal/repository"
)

type SeckillService struct {
	cfg      config.KafkaConfig
	products *repository.ProductRepository
	orders   *repository.OrderRepository
	tx       *repository.TxManager
	idGen    *snowflake.Generator
	locks    *cache.RedisCache
	publish  mq.SeckillPublisher
}

type CreateSeckillOrderRequest struct {
	ProductID      int64  `json:"product_id" binding:"required"`
	Quantity       int    `json:"quantity" binding:"required,min=1"`
	AddressID      int64  `json:"address_id" binding:"required"`
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
}

type SeckillSubmitResult struct {
	Queued    bool         `json:"queued"`
	MessageID string       `json:"message_id"`
	Order     *model.Order `json:"order,omitempty"`
}

func NewSeckillService(
	cfg config.KafkaConfig,
	products *repository.ProductRepository,
	orders *repository.OrderRepository,
	tx *repository.TxManager,
	idGen *snowflake.Generator,
	locks *cache.RedisCache,
) *SeckillService {
	service := &SeckillService{
		cfg:      cfg,
		products: products,
		orders:   orders,
		tx:       tx,
		idGen:    idGen,
		locks:    locks,
	}
	service.publish = mq.NewInlinePublisher(service.HandleMessage)
	return service
}

func (s *SeckillService) SetPublisher(publisher mq.SeckillPublisher) {
	s.publish = publisher
}

func (s *SeckillService) Submit(ctx context.Context, userID int64, req CreateSeckillOrderRequest) (*SeckillSubmitResult, error) {
	message := mq.SeckillMessage{
		MessageID:      fmt.Sprintf("SQ%d", s.idGen.NextID()),
		UserID:         userID,
		ProductID:      req.ProductID,
		Quantity:       req.Quantity,
		AddressID:      req.AddressID,
		IdempotencyKey: req.IdempotencyKey,
	}

	if s.cfg.Enabled {
		if err := s.publish.Publish(ctx, message); err != nil {
			return nil, err
		}
		return &SeckillSubmitResult{
			Queued:    true,
			MessageID: message.MessageID,
		}, nil
	}

	order, err := s.process(ctx, message)
	if err != nil {
		return nil, err
	}
	return &SeckillSubmitResult{
		Queued:    false,
		MessageID: message.MessageID,
		Order:     order,
	}, nil
}

func (s *SeckillService) HandleMessage(ctx context.Context, msg mq.SeckillMessage) error {
	_, err := s.process(ctx, msg)
	return err
}

func (s *SeckillService) process(ctx context.Context, msg mq.SeckillMessage) (*model.Order, error) {
	lockKey := fmt.Sprintf("seckill:lock:%d:%s", msg.UserID, msg.IdempotencyKey)
	locked, err := s.locks.SetNX(ctx, lockKey, "1", 5*time.Minute)
	if err != nil {
		return nil, err
	}
	if !locked {
		return nil, errors.New("duplicate seckill submission")
	}

	product, err := s.products.GetByID(ctx, msg.ProductID)
	if err != nil {
		_ = s.locks.Delete(ctx, lockKey)
		return nil, err
	}
	if product == nil {
		_ = s.locks.Delete(ctx, lockKey)
		return nil, errors.New("product not found")
	}
	if product.Stock < msg.Quantity {
		_ = s.locks.Delete(ctx, lockKey)
		return nil, errors.New("insufficient stock")
	}

	order := &model.Order{
		ID:             s.idGen.NextID(),
		UserID:         msg.UserID,
		AddressID:      msg.AddressID,
		OrderNo:        fmt.Sprintf("SK%s%d", time.Now().Format("20060102150405"), msg.UserID),
		Status:         "pending",
		IdempotencyKey: msg.IdempotencyKey,
		TotalAmount:    product.Price * int64(msg.Quantity),
		Items: []model.OrderItem{
			{
				ID:           s.idGen.NextID(),
				ProductID:    product.ID,
				ProductName:  product.Name,
				ProductPrice: product.Price,
				Quantity:     msg.Quantity,
			},
		},
	}

	err = s.tx.WithTx(ctx, func(txCtx context.Context) error {
		ok, err := s.products.DecreaseStock(txCtx, msg.ProductID, msg.Quantity)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("insufficient stock")
		}
		return s.orders.Create(txCtx, order)
	})
	if err != nil {
		_ = s.locks.Delete(ctx, lockKey)
		return nil, err
	}

	return order, nil
}
