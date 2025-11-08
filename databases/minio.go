package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioConfig struct {
	Enabled    bool   `json:"enabled"`
	Endpoint   string `json:"endpoint"`
	AccessKey  string `json:"accessKey"`
	SecretKey  string `json:"secretKey"`
	Secure     bool   `json:"secure"`
	MaxRetries int    `json:"maxRetries"`
}

type MinioProvider struct {
	Client *minio.Client
}

// connectMinio initializes MinIO client with retries
func connectMinio(raw json.RawMessage) (*MinioProvider, error) {
	var cfg minioConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, nil
	}

	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 5
	}

	var client *minio.Client
	var err error

	for i := 0; i < cfg.MaxRetries; i++ {
		client, err = minio.New(cfg.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: cfg.Secure,
		})
		if err == nil {
			// quick health check: list buckets
			_, err = client.ListBuckets(context.Background())
			if err == nil {
				logs.Info(Ctx, "MinIO connected ✅")
				return &MinioProvider{Client: client}, nil
			}
		}

		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("unable to connect to MinIO at %s: %v", cfg.Endpoint, err)
}

// HealthCheck implements DatabaseProvider
func (m *MinioProvider) HealthCheck(ctx context.Context) error {
	_, err := m.Client.ListBuckets(ctx)
	return err
}

// Close implements DatabaseProvider
func (m *MinioProvider) Close() error {
	// MinIO client does not require explicit close
	return nil
}
