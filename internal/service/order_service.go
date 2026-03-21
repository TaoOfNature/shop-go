package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TaoOfNature/shop-go/internal/model"
	"github.com/TaoOfNature/shop-go/internal/pkg/snowflake"
	"github.com/TaoOfNature/shop-go/internal/repository"
)

type OrderService struct {
	orders   *repository.OrderRepository
	carts    *repository.CartRepository
	products *repository.ProductRepository
	tx       *repository.TxManager
	idGen    *snowflake.Generator
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
) *OrderService {
	return &OrderService{orders: orders, carts: carts, products: products, tx: tx, idGen: idGen}
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
	cartItems, err := s.carts.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	productIDs := make([]int64, 0, len(cartItems))
	for _, item := range cartItems {
		if item.Checked {
			productIDs = append(productIDs, item.ProductID)
		}
	}
	if len(productIDs) == 0 {
		return nil, errors.New("no checked cart items")
	}

	products, err := s.products.GetByIDs(ctx, productIDs)
	if err != nil {
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
			return nil, errors.New("product not found in cart")
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
		return s.orders.Create(txCtx, order)
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}
