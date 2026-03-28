package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// DBs
	PostgresDSN   string
	MongoURI      string
	RedisAddr     string
	RedisPass     string
	OpenSearchURL string

	// Storage
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string

	// 🔐 Auth (NEW)
	JWTSecret              string
	RefreshSecret          string
	AccessTokenExpiryHours int
	RefreshTokenExpiryDays int
}

func LoadConfig() *Config {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := &Config{
		// DB
		PostgresDSN:   getEnv("POSTGRES_DSN", ""),
		MongoURI:      getEnv("MONGO_URI", ""),
		RedisAddr:     getEnv("REDIS_ADDR", ""),
		RedisPass:     getEnv("REDIS_PASS", ""),
		OpenSearchURL: getEnv("OPENSEARCH_URL", ""),

		// Storage
		MinioEndpoint:  getEnv("MINIO_ENDPOINT", ""),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", ""),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", ""),

		// 🔐 Auth
		JWTSecret:              getEnv("JWT_SECRET", "dev-secret"), // fallback for dev
		RefreshSecret:          getEnv("REFRESH_SECRET", "dev-refresh-secret"),
		AccessTokenExpiryHours: getEnvAsInt("ACCESS_TOKEN_EXPIRY_HOURS", 1),
		RefreshTokenExpiryDays: getEnvAsInt("REFRESH_TOKEN_EXPIRY_DAYS", 1),
	}

	validateConfig(cfg)

	return cfg
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Printf("Invalid int for %s, using default %d\n", key, defaultVal)
		return defaultVal
	}
	return val
}

// 🔥 Validate critical configs (production safety)
func validateConfig(cfg *Config) {
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	if cfg.RefreshSecret == "" {
		log.Fatal("REFRESH_SECRET is required")
	}
}
