package service

import (
	"context"
	"errors"

	"github.com/dawnstack/shop-go/internal/cache"
	"github.com/dawnstack/shop-go/internal/model"
	"github.com/dawnstack/shop-go/internal/repository"
)

type ProductService struct {
	products *repository.ProductRepository
	cache    *cache.ProductCache
}

func NewProductService(products *repository.ProductRepository, cache *cache.ProductCache) *ProductService {
	return &ProductService{products: products, cache: cache}
}

func (s *ProductService) Search(ctx context.Context, keyword string, sort string) ([]model.Product, error) {
	return s.products.Search(ctx, keyword, sort)
}

func (s *ProductService) Detail(ctx context.Context, id int64) (*model.Product, error) {
	return s.cache.GetOrLoad(ctx, id, func(ctx context.Context) (*model.Product, error) {
		product, err := s.products.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, errors.New("product not found")
		}
		return product, nil
	})
}
