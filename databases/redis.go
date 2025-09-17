package database

import (
	"context"
	"log"
	"time"

	logs "Adornme/logging"

	"github.com/redis/go-redis/v9"
)

var logger = logs.InitLogger("database")

func ConnectRedis(addr, password string) (*redis.Client, error) {
	logger.Info("Connecting to database...")

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Println("Redis connected")
	return rdb, nil
}
