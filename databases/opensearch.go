package database

import (
	"log"
	"net/http"
	"time"

	"github.com/opensearch-project/opensearch-go"
)

func ConnectOpenSearch(url string) (*opensearch.Client, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{url},
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     30 * time.Second,
		},
	})
	if err != nil {
		return nil, err
	}

	// Retry health check up to 5 times
	for i := 0; i < 5; i++ {
		res, err := client.Cluster.Health()
		if err == nil {
			defer res.Body.Close()
			log.Println("OpenSearch connected")
			return client, nil
		}

		log.Println("OpenSearch not ready, retrying in 3s...", err)
		time.Sleep(3 * time.Second)
	}

	return nil, err
}
