package service

import (
	"context"

	"github.com/dawnstack/shop-go/internal/cache"
	"github.com/dawnstack/shop-go/internal/model"
	"github.com/dawnstack/shop-go/internal/repository"
)

type HomeService struct {
	repo  *repository.HomeRepository
	cache *cache.HomeCache
}

func NewHomeService(repo *repository.HomeRepository, cache *cache.HomeCache) *HomeService {
	return &HomeService{repo: repo, cache: cache}
}

func (s *HomeService) Banners(ctx context.Context) ([]model.Banner, error) {
	return s.cache.GetOrLoadBanners(ctx, s.repo.ListBanners)
}

func (s *HomeService) Categories(ctx context.Context) ([]model.Category, error) {
	return s.cache.GetOrLoadCategories(ctx, s.repo.ListCategories)
}

func (s *HomeService) Recommend(ctx context.Context) ([]model.Product, error) {
	return s.cache.GetOrLoadRecommend(ctx, s.repo.ListRecommendations)
}
