package main

import (
	"Adornme/config"
	db "Adornme/databases"
	health "Adornme/internal/health"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	postgres, err := db.ConnectPostgres(cfg.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}
	mongo, err := db.ConnectMongo(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}

	redis, err := db.ConnectRedis(cfg.RedisAddr, cfg.RedisPass)
	if err != nil {
		log.Fatal(err)
	}

	openSearch, err := db.ConnectOpenSearch(cfg.OpenSearchURL)
	if err != nil {
		log.Fatal(err)
	}

	minioClient, err := db.ConnectMinio(cfg.MinioEndpoint, cfg.MinioAccessKey, cfg.MinioSecretKey)
	if err != nil {
		log.Fatal(err)
	}

	// Health API
	http.HandleFunc("/health", health.HealthCheckHandler(postgres, mongo, redis, openSearch, minioClient))

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
