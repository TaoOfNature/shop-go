package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/dawnstack/shop-go/internal/config"
	"github.com/jackc/pgx/v5/pgconn"
)

func StartCacheWarmup(ctx context.Context, services *Services, cfg config.Config) {
	interval := time.Duration(cfg.App.PrewarmSeconds) * time.Second
	if interval <= 0 {
		interval = 5 * time.Minute
	}

	run := func() {
		if err := services.Home.Prewarm(ctx); err != nil {
			if !isMissingTable(err) {
				log.Printf("home cache prewarm failed: %v", err)
			}
		}
		if err := services.Video.Prewarm(ctx); err != nil {
			if !isMissingTable(err) {
				log.Printf("video cache prewarm failed: %v", err)
			}
		}
	}

	run()

	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				run()
			}
		}
	}()
}

func isMissingTable(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "42P01"
	}
	return false
}
