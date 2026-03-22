package service

import (
	"context"

	"github.com/dawnstack/shop-go/internal/cache"
	"github.com/dawnstack/shop-go/internal/model"
	"github.com/dawnstack/shop-go/internal/repository"
)

type VideoService struct {
	repo  *repository.HomeRepository
	cache *cache.HomeCache
}

func NewVideoService(repo *repository.HomeRepository, cache *cache.HomeCache) *VideoService {
	return &VideoService{repo: repo, cache: cache}
}

func (s *VideoService) Recommend(ctx context.Context) ([]model.Video, error) {
	if cached, ok := s.cache.Get("videos"); ok {
		if items, ok := cached.([]model.Video); ok {
			return items, nil
		}
	}
	items, err := s.repo.ListVideos(ctx)
	if err != nil {
		return nil, err
	}
	s.cache.Set("videos", items)
	return items, nil
}
