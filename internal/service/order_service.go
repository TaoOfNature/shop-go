package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dawnstack/shop-go/internal/cache"
	"github.com/dawnstack/shop-go/internal/model"
	"github.com/dawnstack/shop-go/internal/pkg/snowflake"
	"github.com/dawnstack/shop-go/internal/repository"
)

type OrderService struct {
	orders   *repository.OrderRepository
	carts    *repository.CartRepository
	products *repository.ProductRepository
	tx       *repository.TxManager
	idGen    *snowflake.Generator
	locks    *cache.RedisCache
}

type CreateOrderRequest struct {
	AddressID      int64 `json:"address_id" binding:"required"`
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
}

func NewOrderService(
	orders *repository.OrderRepository,
	carts *repository.CartRepository,
	products *repository.ProductRepository,
	tx *repository.TxManager,
	idGen *snowflake.Generator,
	locks *cache.RedisCache,
) *OrderService {
	return &OrderService{orders: orders, carts: carts, products: products, tx: tx, idGen: idGen, locks: locks}
}

func (s *OrderService) List(ctx context.Context, userID int64) ([]model.Order, error) {
	return s.orders.ListByUserID(ctx, userID)
}

func (s *OrderService) Detail(ctx context.Context, userID int64, id int64) (*model.Order, error) {
	order, err := s.orders.GetByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (s *OrderService) Create(ctx context.Context, userID int64, req CreateOrderRequest) (*model.Order, error) {
	lockKey := fmt.Sprintf("order:lock:%d:%s", userID, req.IdempotencyKey)
	locked, err := s.locks.SetNX(ctx, lockKey, "1", 5*time.Minute)
	if err != nil {
		return nil, err
	}
	if !locked {
		return nil, errors.New("duplicate order submission")
	}

	cartItems, err := s.carts.ListByUserID(ctx, userID)
	if err != nil {
		_ = s.locks.Delete(ctx, lockKey)
		return nil, err
	}
	if len(cartItems) == 0 {
		_ = s.locks.Delete(ctx, lockKey)
		return nil, errors.New("cart is empty")
	}

	productIDs := make([]int64, 0, len(cartItems))
	for _, item := range cartItems {
		if item.Checked {
			productIDs = append(productIDs, item.ProductID)
		}
	}
	if len(productIDs) == 0 {
		_ = s.locks.Delete(ctx, lockKey)
		return nil, errors.New("no checked cart items")
	}

	products, err := s.products.GetByIDs(ctx, productIDs)
	if err != nil {
		_ = s.locks.Delete(ctx, lockKey)
		return nil, err
	}

	productMap := make(map[int64]model.Product, len(products))
	for _, product := range products {
		productMap[product.ID] = product
	}

	order := &model.Order{
		ID:             s.idGen.NextID(),
		UserID:         userID,
		AddressID:      req.AddressID,
		OrderNo:        fmt.Sprintf("M%s%d", time.Now().Format("20060102150405"), userID),
		Status:         "pending",
		IdempotencyKey: req.IdempotencyKey,
	}

	for _, cartItem := range cartItems {
		if !cartItem.Checked {
			continue
		}
		product, ok := productMap[cartItem.ProductID]
		if !ok {
			_ = s.locks.Delete(ctx, lockKey)
			return nil, errors.New("product not found in cart")
		}
		if product.Stock < cartItem.Quantity {
			_ = s.locks.Delete(ctx, lockKey)
			return nil, errors.New("insufficient stock")
		}
		order.TotalAmount += product.Price * int64(cartItem.Quantity)
		order.Items = append(order.Items, model.OrderItem{
			ID:           s.idGen.NextID(),
			OrderID:      order.ID,
			ProductID:    product.ID,
			ProductName:  product.Name,
			ProductPrice: product.Price,
			Quantity:     cartItem.Quantity,
		})
	}

	err = s.tx.WithTx(ctx, func(txCtx context.Context) error {
		for _, item := range order.Items {
			ok, err := s.products.DecreaseStock(txCtx, item.ProductID, item.Quantity)
			if err != nil {
				return err
			}
			if !ok {
				return errors.New("insufficient stock")
			}
		}

		if err := s.orders.Create(txCtx, order); err != nil {
			return err
		}

		return s.carts.DeleteCheckedByUserID(txCtx, userID)
	})
	if err != nil {
		_ = s.locks.Delete(ctx, lockKey)
		return nil, err
	}
	return order, nil
}
