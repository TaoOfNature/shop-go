package cache

import "time"

type HomeCache struct {
	store *RedisCache
}

func NewHomeCache(store *RedisCache) *HomeCache {
	return &HomeCache{store: store}
}

func (c *HomeCache) Get(key string) (interface{}, bool) {
	return c.store.Get("home:" + key)
}

func (c *HomeCache) Set(key string, value interface{}) {
	c.store.Set("home:"+key, value, 2*time.Minute)
}
