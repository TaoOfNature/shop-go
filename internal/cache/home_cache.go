package cache

import (
	"context"
	"time"

	"github.com/dawnstack/shop-go/internal/model"
	"golang.org/x/sync/singleflight"
)

type HomeCache struct {
	store *RedisCache
	group singleflight.Group
}

func NewHomeCache(store *RedisCache) *HomeCache {
	return &HomeCache{store: store}
}

func (c *HomeCache) GetOrLoadBanners(ctx context.Context, loader func(context.Context) ([]model.Banner, error)) ([]model.Banner, error) {
	return getOrLoad(c, ctx, "banners", 2*time.Minute, loader)
}

func getOrLoad[T any](c *HomeCache, ctx context.Context, key string, ttl time.Duration, loader func(context.Context) (T, error)) (T, error) {
	cacheKey := "home:" + key

	var cached T
	ok, err := c.store.GetJSON(ctx, cacheKey, &cached)
	if err == nil && ok {
		return cached, nil
	}

	value, err, _ := c.group.Do(cacheKey, func() (interface{}, error) {
		result, loadErr := loader(ctx)
		if loadErr != nil {
			return nil, loadErr
		}
		if setErr := c.store.SetJSON(ctx, cacheKey, result, ttl); setErr != nil {
			return nil, setErr
		}
		return result, nil
	})
	if err != nil {
		var zero T
		return zero, err
	}
	return value.(T), nil
}

func (c *HomeCache) GetOrLoadCategories(ctx context.Context, loader func(context.Context) ([]model.Category, error)) ([]model.Category, error) {
	return getOrLoad(c, ctx, "categories", 2*time.Minute, loader)
}

func (c *HomeCache) GetOrLoadRecommend(ctx context.Context, loader func(context.Context) ([]model.Product, error)) ([]model.Product, error) {
	return getOrLoad(c, ctx, "recommend", 2*time.Minute, loader)
}

func (c *HomeCache) GetOrLoadVideos(ctx context.Context, loader func(context.Context) ([]model.Video, error)) ([]model.Video, error) {
	return getOrLoad(c, ctx, "videos", 2*time.Minute, loader)
}
