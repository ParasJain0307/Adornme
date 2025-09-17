package database

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func ConnectMinio(endpoint, accessKey, secretKey string) (*minio.Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, // use true if SSL enabled
	})
	if err != nil {
		return nil, err
	}

	log.Println("MinIO connected")
	return minioClient, nil
}
