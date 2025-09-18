package database

import (
	"context"
	"encoding/json"

	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoConfig struct {
	Enabled     bool   `json:"enabled"`
	DSN         string `json:"dsn"`
	MaxPoolSize uint64 `json:"maxPoolSize"`
}

type MongoProvider struct {
	Client *mongo.Client
}

func connectMongo(raw json.RawMessage) (*MongoProvider, error) {
	var cfg mongoConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DSN).SetMaxPoolSize(cfg.MaxPoolSize))
	if err != nil {
		return nil, err
	}

	logs.Info(ctx, "MongoDB connected ✅")
	return &MongoProvider{Client: client}, nil
}

func (m *MongoProvider) HealthCheck(ctx context.Context) error {
	return m.Client.Ping(ctx, nil)
}

func (m *MongoProvider) Close() error {
	return m.Client.Disconnect(context.Background())
}
