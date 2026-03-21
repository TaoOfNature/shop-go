package service

import (
	"context"
	"errors"

	"github.com/TaoOfNature/shop-go/internal/cache"
	"github.com/TaoOfNature/shop-go/internal/model"
	"github.com/TaoOfNature/shop-go/internal/repository"
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
	if cached, ok := s.cache.Get(id); ok {
		if product, ok := cached.(*model.Product); ok {
			return product, nil
		}
	}

	product, err := s.products.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	s.cache.Set(id, product)
	return product, nil
}
