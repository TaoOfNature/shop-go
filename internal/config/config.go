package config

import (
	"os"
	"strings"
	"strconv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Kafka    KafkaConfig
}

type AppConfig struct {
	Name           string
	Env            string
	Port           string
	RateLimitRPS   int
	RateLimitBurst int
	PrewarmSeconds int
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	Database int
}

type JWTConfig struct {
	Secret            string
	Issuer            string
	AccessExpireMin   int
	RefreshExpireHour int
}

type KafkaConfig struct {
	Enabled bool
	Brokers []string
	Topic   string
	GroupID string
}

func Load() Config {
	return Config{
		App: AppConfig{
			Name:           getenv("APP_NAME", "shop-go"),
			Env:            getenv("APP_ENV", "dev"),
			Port:           getenv("APP_PORT", "8080"),
			RateLimitRPS:   getenvInt("RATE_LIMIT_RPS", 10),
			RateLimitBurst: getenvInt("RATE_LIMIT_BURST", 20),
			PrewarmSeconds: getenvInt("CACHE_PREWARM_SECONDS", 300),
		},
		Database: DatabaseConfig{
			Host:     getenv("DB_HOST", "127.0.0.1"),
			Port:     getenv("DB_PORT", "5432"),
			User:     getenv("DB_USER", "postgres"),
			Password: getenv("DB_PASS", "postgres"),
			Name:     getenv("DB_NAME", "mall"),
			SSLMode:  getenv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     getenv("REDIS_ADDR", "127.0.0.1:6379"),
			Password: getenv("REDIS_PASS", ""),
			Database: getenvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:            getenv("JWT_SECRET", "change-me"),
			Issuer:            getenv("JWT_ISSUER", "shop-go"),
			AccessExpireMin:   getenvInt("JWT_ACCESS_EXPIRE_MIN", 120),
			RefreshExpireHour: getenvInt("JWT_REFRESH_EXPIRE_HOUR", 168),
		},
		Kafka: KafkaConfig{
			Enabled: getenv("KAFKA_ENABLED", "false") == "true",
			Brokers: splitCSV(getenv("KAFKA_BROKERS", "127.0.0.1:9092")),
			Topic:   getenv("KAFKA_TOPIC", "shop-go-seckill"),
			GroupID: getenv("KAFKA_GROUP_ID", "shop-go-seckill-consumer"),
		},
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return number
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	if len(result) == 0 {
		return []string{"127.0.0.1:9092"}
	}
	return result
}
