package cache

import (
	"testing"
	"time"

	"github.com/TaoOfNature/shop-go/internal/config"
)

func TestRedisCacheSetGetDelete(t *testing.T) {
	cache := NewRedisCache(config.RedisConfig{})
	cache.Set("k", "v", 0)

	value, ok := cache.Get("k")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if value.(string) != "v" {
		t.Fatalf("unexpected value: %#v", value)
	}

	cache.Delete("k")
	if _, ok := cache.Get("k"); ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestRedisCacheTTL(t *testing.T) {
	cache := NewRedisCache(config.RedisConfig{})
	cache.Set("k", "v", 10*time.Millisecond)

	time.Sleep(30 * time.Millisecond)

	if _, ok := cache.Get("k"); ok {
		t.Fatal("expected key to expire")
	}
}
