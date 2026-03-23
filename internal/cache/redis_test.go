package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/dawnstack/shop-go/internal/config"
)

type testPayload struct {
	Value string `json:"value"`
}

func TestRedisCacheSetGetDelete(t *testing.T) {
	server := miniredis.RunT(t)
	cache := NewRedisCache(config.RedisConfig{Addr: server.Addr()})
	ctx := context.Background()

	if err := cache.SetJSON(ctx, "k", testPayload{Value: "v"}, 0); err != nil {
		t.Fatalf("SetJSON() error = %v", err)
	}

	var payload testPayload
	ok, err := cache.GetJSON(ctx, "k", &payload)
	if err != nil {
		t.Fatalf("GetJSON() error = %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}
	if payload.Value != "v" {
		t.Fatalf("unexpected value: %#v", payload)
	}

	if err := cache.Delete(ctx, "k"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	ok, err = cache.GetJSON(ctx, "k", &payload)
	if err != nil {
		t.Fatalf("GetJSON() after delete error = %v", err)
	}
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestRedisCacheTTL(t *testing.T) {
	server := miniredis.RunT(t)
	cache := NewRedisCache(config.RedisConfig{Addr: server.Addr()})
	ctx := context.Background()

	if err := cache.SetJSON(ctx, "k", testPayload{Value: "v"}, 10*time.Millisecond); err != nil {
		t.Fatalf("SetJSON() error = %v", err)
	}

	server.FastForward(100 * time.Millisecond)

	var payload testPayload
	ok, err := cache.GetJSON(ctx, "k", &payload)
	if err != nil {
		t.Fatalf("GetJSON() error = %v", err)
	}
	if ok {
		t.Fatal("expected key to expire")
	}
}
