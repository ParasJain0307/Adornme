package database

import (
	"context"
	"encoding/json"
	"fmt"

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

// ConnectMongo initializes MongoDB connection with retry in background
func connectMongo(raw json.RawMessage) (*MongoProvider, error) {
	var cfg mongoConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, nil
	}

	provider := &MongoProvider{}

	// start background retry loop
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			client, err := mongo.Connect(ctx,
				options.Client().
					ApplyURI(cfg.DSN).
					SetMaxPoolSize(cfg.MaxPoolSize),
			)

			if err != nil {
				logs.Error(Ctx, "MongoDB connection failed ❌", "error", err)
				time.Sleep(5 * time.Second) // wait before retry
				continue
			}

			// verify connection with ping
			if err := client.Ping(ctx, nil); err != nil {
				logs.Error(Ctx, "MongoDB ping failed ❌", "error", err)
				time.Sleep(5 * time.Second)
				continue
			}

			// success
			provider.Client = client
			logs.Info(Ctx, "MongoDB connected ✅")
			break
		}
	}()

	return provider, nil
}

func (m *MongoProvider) HealthDetails(ctx context.Context) (uptime string, latencyMs float64, err error) {
	if m == nil || m.Client == nil {
		return "", 0, fmt.Errorf("mongo client is nil")
	}

	start := time.Now()
	if err := m.Client.Ping(ctx, nil); err != nil {
		return "", 0, err
	}
	latencyMs = time.Since(start).Seconds() * 1000

	var result struct {
		Uptime float64 `bson:"uptime"`
	}

	if err := m.Client.Database("admin").RunCommand(ctx, map[string]int{"serverStatus": 1}).Decode(&result); err != nil {
		return "", latencyMs, fmt.Errorf("failed to get serverStatus: %w", err)
	}

	uptime = time.Duration(result.Uptime * float64(time.Second)).Truncate(time.Second).String()
	return uptime, latencyMs, nil
}

func (m *MongoProvider) HealthCheck(ctx context.Context) error {
	return m.Client.Ping(ctx, nil)
}

func (m *MongoProvider) Close() error {
	return m.Client.Disconnect(context.Background())
}
