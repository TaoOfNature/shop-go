package service

import (
	"context"

	"github.com/TaoOfNature/shop-go/internal/model"
	"github.com/TaoOfNature/shop-go/internal/pkg/snowflake"
	"github.com/TaoOfNature/shop-go/internal/repository"
)

type CartService struct {
	carts    *repository.CartRepository
	products *repository.ProductRepository
	idGen    *snowflake.Generator
}

type UpsertCartRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=1"`
	Checked   bool  `json:"checked"`
}

type UpdateCartRequest struct {
	Quantity int  `json:"quantity" binding:"required,min=1"`
	Checked  bool `json:"checked"`
}

func NewCartService(carts *repository.CartRepository, products *repository.ProductRepository, idGen *snowflake.Generator) *CartService {
	return &CartService{carts: carts, products: products, idGen: idGen}
}

func (s *CartService) List(ctx context.Context, userID int64) ([]model.CartItem, error) {
	return s.carts.ListByUserID(ctx, userID)
}

func (s *CartService) Add(ctx context.Context, userID int64, req UpsertCartRequest) error {
	return s.carts.Upsert(ctx, &model.CartItem{
		ID:        s.idGen.NextID(),
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Checked:   req.Checked,
	})
}

func (s *CartService) Update(ctx context.Context, userID int64, id int64, req UpdateCartRequest) error {
	return s.carts.Update(ctx, &model.CartItem{
		ID:       id,
		UserID:   userID,
		Quantity: req.Quantity,
		Checked:  req.Checked,
	})
}

func (s *CartService) Delete(ctx context.Context, userID int64, id int64) error {
	return s.carts.Delete(ctx, userID, id)
}
