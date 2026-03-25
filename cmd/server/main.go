package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/dawnstack/shop-go/internal/api"
	"github.com/dawnstack/shop-go/internal/cache"
	"github.com/dawnstack/shop-go/internal/config"
	"github.com/dawnstack/shop-go/internal/repository"
	"github.com/dawnstack/shop-go/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := repository.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	txManager := repository.NewTxManager(db)

	redisCache := cache.NewRedisCache(cfg.Redis)
	if err := redisCache.Ping(context.Background()); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	productCache := cache.NewProductCache(redisCache)
	homeCache := cache.NewHomeCache(redisCache)

	repos := repository.NewRepositories(db, txManager)
	services := service.NewServices(repos, txManager, redisCache, productCache, homeCache, cfg)
	router := api.NewRouter(services, cfg)

	server := &http.Server{
		Addr:              ":" + cfg.App.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	service.StartCacheWarmup(ctx, services, cfg)

	go func() {
		log.Printf("mall backend listening on :%s", cfg.App.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to listen and serve: %v", err)
		}
	}()

	<-ctx.Done()

	gracefulCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(gracefulCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
