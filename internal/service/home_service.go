package service

import (
	"context"

	"github.com/TaoOfNature/shop-go/internal/cache"
	"github.com/TaoOfNature/shop-go/internal/model"
	"github.com/TaoOfNature/shop-go/internal/repository"
)

type HomeService struct {
	repo  *repository.HomeRepository
	cache *cache.HomeCache
}

func NewHomeService(repo *repository.HomeRepository, cache *cache.HomeCache) *HomeService {
	return &HomeService{repo: repo, cache: cache}
}

func (s *HomeService) Banners(ctx context.Context) ([]model.Banner, error) {
	if cached, ok := s.cache.Get("banners"); ok {
		if items, ok := cached.([]model.Banner); ok {
			return items, nil
		}
	}
	items, err := s.repo.ListBanners(ctx)
	if err != nil {
		return nil, err
	}
	s.cache.Set("banners", items)
	return items, nil
}

func (s *HomeService) Categories(ctx context.Context) ([]model.Category, error) {
	if cached, ok := s.cache.Get("categories"); ok {
		if items, ok := cached.([]model.Category); ok {
			return items, nil
		}
	}
	items, err := s.repo.ListCategories(ctx)
	if err != nil {
		return nil, err
	}
	s.cache.Set("categories", items)
	return items, nil
}

func (s *HomeService) Recommend(ctx context.Context) ([]model.Product, error) {
	if cached, ok := s.cache.Get("recommend"); ok {
		if items, ok := cached.([]model.Product); ok {
			return items, nil
		}
	}
	items, err := s.repo.ListRecommendations(ctx)
	if err != nil {
		return nil, err
	}
	s.cache.Set("recommend", items)
	return items, nil
}
