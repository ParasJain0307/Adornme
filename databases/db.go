package database

import (
	log "Adornme/logging"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	Do   map[string]DatabaseProvider // registry of all DBs
	once sync.Once
	ctx  = log.WithRequestID(context.Background(), "req-12345")
	logs = log.Component("database")
)

type DatabaseProvider interface {
	HealthCheck(ctx context.Context) error
	Close() error
}

// Init DB registry
func init() {
	once.Do(func() {
		Do = make(map[string]DatabaseProvider)
		if err := setupDatabases("config/db-config.json"); err != nil {
			logs.Fatalf(ctx, "DB initialization failed: %v", err)
		}
	})
}

// Setup databases from JSON config
func setupDatabases(cfgPath string) error {
	logs.Noticef(ctx, "Setup Databases Called!")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return fmt.Errorf("cannot read db-config.json: %w", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("cannot parse db-config.json: %w", err)
	}

	// Postgres
	if r, ok := raw["postgres"]; ok {
		var pgRaw map[string]json.RawMessage
		if err := json.Unmarshal(r, &pgRaw); err != nil {
			return fmt.Errorf("cannot parse postgres config: %w", err)
		}

		pgClients, err := ConnectAllPostgres(pgRaw) // ✅ now only postgres section
		if err != nil {
			return err
		}

		Do["postgres"] = pgClients
	}

	// Mongo
	if r, ok := raw["mongo"]; ok {
		mg, err := connectMongo(r)
		if err != nil {
			return err
		}
		Do["mongo"] = mg
	}

	// Redis
	if r, ok := raw["redis"]; ok {
		rd, err := connectRedis(r)
		if err != nil {
			return err
		}
		Do["redis"] = rd
	}

	// OpenSearch
	if r, ok := raw["opensearch"]; ok {
		es, err := connectOpenSearch(r)
		if err != nil {
			return err
		}
		Do["opensearch"] = es
	}

	// MinIO
	if r, ok := raw["minio"]; ok {
		minioClient, err := connectMinio(r)
		if err != nil {
			return err
		}
		Do["minio"] = minioClient
	}

	// Start periodic health check
	stopCh := make(chan struct{})
	go StartDBHealthTicker(5*time.Minute, stopCh)

	return nil
}

// Graceful shutdown
func CloseAll() {
	for name, db := range Do {
		if err := db.Close(); err != nil {
			logs.Errorf(ctx, "Error closing %s: %v\n", name, err)
		}
	}
}

// Health check for all DBs
func CheckAll(ctx context.Context) {
	for name, db := range Do {
		if err := db.HealthCheck(ctx); err != nil {
			logs.Errorf(ctx, "%s unhealthy: %v\n", name, err)
		} else {
			logs.Infof(ctx, "%s healthy ✅\n", name)
		}
	}
}
