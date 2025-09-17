package health

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/opensearch-project/opensearch-go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthResponse struct {
	Postgres   string `json:"postgres"`
	Mongo      string `json:"mongo"`
	Redis      string `json:"redis"`
	OpenSearch string `json:"opensearch"`
	MinIO      string `json:"minio"`
}

func HealthCheckHandler(pg *pgxpool.Pool, mongo *mongo.Client, redis *redis.Client, os *opensearch.Client, minioClient *minio.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := HealthResponse{
			Postgres:   checkPostgres(pg),
			Mongo:      checkMongo(mongo),
			Redis:      checkRedis(redis),
			OpenSearch: checkOpenSearch(os),
			MinIO:      checkMinio(minioClient),
		}

		w.Header().Set("Content-Type", "application/json")
		statusCode := http.StatusOK
		if resp.Postgres != "ok" || resp.Mongo != "ok" || resp.Redis != "ok" || resp.OpenSearch != "ok" || resp.MinIO != "ok" {
			statusCode = http.StatusServiceUnavailable
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(resp)
	}
}

func checkPostgres(pg *pgxpool.Pool) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := pg.Ping(ctx); err != nil {
		log.Println("Postgres health error:", err)
		return "down"
	}
	return "ok"
}

func checkMongo(client *mongo.Client) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Println("Mongo health error:", err)
		return "down"
	}
	return "ok"
}

func checkRedis(client *redis.Client) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Println("Redis health error:", err)
		return "down"
	}
	return "ok"
}

func checkOpenSearch(client *opensearch.Client) string {
	res, err := client.Cluster.Health()
	if err != nil {
		log.Println("OpenSearch health error:", err)
		return "down"
	}
	res.Body.Close()
	return "ok"
}

func checkMinio(client *minio.Client) string {
	_, err := client.ListBuckets(context.Background())
	if err != nil {
		log.Println("MinIO health error:", err)
		return "down"
	}
	return "ok"
}
