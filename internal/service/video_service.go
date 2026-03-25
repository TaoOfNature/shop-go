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
	return s.cache.GetOrLoadVideos(ctx, s.repo.ListVideos)
}

func (s *VideoService) Prewarm(ctx context.Context) error {
	_, err := s.Recommend(ctx)
	return err
}
