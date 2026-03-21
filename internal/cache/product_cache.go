package cache

import (
	"fmt"
	"time"
)

type ProductCache struct {
	store *RedisCache
}

func NewProductCache(store *RedisCache) *ProductCache {
	return &ProductCache{store: store}
}

func (c *ProductCache) Get(productID int64) (interface{}, bool) {
	return c.store.Get(fmt.Sprintf("product:%d", productID))
}

func (c *ProductCache) Set(productID int64, value interface{}) {
	c.store.Set(fmt.Sprintf("product:%d", productID), value, 5*time.Minute)
}

func (c *ProductCache) Invalidate(productID int64) {
	c.store.Delete(fmt.Sprintf("product:%d", productID))
}
