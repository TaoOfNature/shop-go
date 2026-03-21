package config

import "testing"

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("APP_PORT", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("JWT_SECRET", "")

	cfg := Load()

	if cfg.App.Port != "8080" {
		t.Fatalf("unexpected default port: %s", cfg.App.Port)
	}
	if cfg.Database.Host != "127.0.0.1" {
		t.Fatalf("unexpected default db host: %s", cfg.Database.Host)
	}
	if cfg.JWT.Secret != "change-me" {
		t.Fatalf("unexpected default jwt secret: %s", cfg.JWT.Secret)
	}
}

func TestLoadUsesEnvironmentOverrides(t *testing.T) {
	t.Setenv("APP_PORT", "9090")
	t.Setenv("RATE_LIMIT_RPS", "99")
	t.Setenv("DB_HOST", "db")
	t.Setenv("JWT_SECRET", "secret")

	cfg := Load()

	if cfg.App.Port != "9090" {
		t.Fatalf("unexpected app port: %s", cfg.App.Port)
	}
	if cfg.App.RateLimitRPS != 99 {
		t.Fatalf("unexpected rate limit rps: %d", cfg.App.RateLimitRPS)
	}
	if cfg.Database.Host != "db" {
		t.Fatalf("unexpected db host: %s", cfg.Database.Host)
	}
	if cfg.JWT.Secret != "secret" {
		t.Fatalf("unexpected jwt secret: %s", cfg.JWT.Secret)
	}
}
