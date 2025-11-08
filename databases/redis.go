package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisConfig struct {
	Enabled    bool   `json:"enabled"`
	Addr       string `json:"addr"`
	Password   string `json:"password"`
	DB         int    `json:"db"`
	MaxRetries int    `json:"maxRetries"`
	Timeout    int    `json:"timeoutSeconds"`
}

type RedisProvider struct {
	Client *redis.Client
}

// connectRedis initializes Redis client with retries
func connectRedis(raw json.RawMessage) (*RedisProvider, error) {
	var cfg redisConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, nil
	}

	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 5
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 5
	}

	var client *redis.Client
	var err error

	for i := 0; i < cfg.MaxRetries; i++ {
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		})

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
		defer cancel()

		err = client.Ping(ctx).Err()
		if err == nil {
			logs.Info(Ctx, "Redis connected ✅")
			return &RedisProvider{Client: client}, nil
		}

		logs.Warningf(Ctx, "Redis not ready (%v), retrying in 2s...\n", err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("unable to connect to Redis at %s: %v", cfg.Addr, err)
}

// HealthCheck implements DatabaseProvider interface
func (r *RedisProvider) HealthCheck(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

// Close implements DatabaseProvider interface
func (r *RedisProvider) Close() error {
	return r.Client.Close()
}
