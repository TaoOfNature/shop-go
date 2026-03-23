package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/dawnstack/shop-go/internal/model"
	"golang.org/x/sync/singleflight"
)

type ProductCache struct {
	store *RedisCache
	group singleflight.Group
}

func NewProductCache(store *RedisCache) *ProductCache {
	return &ProductCache{store: store}
}

func (c *ProductCache) GetOrLoad(ctx context.Context, productID int64, loader func(context.Context) (*model.Product, error)) (*model.Product, error) {
	key := fmt.Sprintf("product:%d", productID)

	var cached model.Product
	ok, err := c.store.GetJSON(ctx, key, &cached)
	if err == nil && ok {
		return &cached, nil
	}

	value, err, _ := c.group.Do(key, func() (interface{}, error) {
		product, loadErr := loader(ctx)
		if loadErr != nil {
			return nil, loadErr
		}
		if setErr := c.store.SetJSON(ctx, key, product, 5*time.Minute); setErr != nil {
			return nil, setErr
		}
		return product, nil
	})
	if err != nil {
		return nil, err
	}
	return value.(*model.Product), nil
}

func (c *ProductCache) Invalidate(ctx context.Context, productID int64) error {
	return c.store.Delete(ctx, fmt.Sprintf("product:%d", productID))
}
