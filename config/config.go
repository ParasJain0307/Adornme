package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresDSN    string
	MongoURI       string
	RedisAddr      string
	RedisPass      string
	OpenSearchURL  string
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
}

func LoadConfig() *Config {
	// Load .env from project root
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		PostgresDSN:    os.Getenv("POSTGRES_DSN"),   // e.g., "host=localhost port=5432 user=postgres password=postgres dbname=ecommerce sslmode=disable"
		MongoURI:       os.Getenv("MONGO_URI"),      // e.g., "mongodb://localhost:27017"
		RedisAddr:      os.Getenv("REDIS_ADDR"),     // e.g., "localhost:6380"
		RedisPass:      os.Getenv("REDIS_PASS"),     // e.g., ""
		OpenSearchURL:  os.Getenv("OPENSEARCH_URL"), // e.g., "http://localhost:9200"
		MinioEndpoint:  os.Getenv("MINIO_ENDPOINT"), // e.g., "localhost:9000"
		MinioAccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		MinioSecretKey: os.Getenv("MINIO_SECRET_KEY"),
	}
}
