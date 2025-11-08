package database

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/opensearch-project/opensearch-go"
)

type opensearchConfig struct {
	Enabled bool   `json:"enabled"`
	URL     string `json:"url"`
	MaxIdle int    `json:"maxIdleConnsPerHost"`
	Timeout int    `json:"idleConnTimeoutSeconds"`
}

type OpenSearchProvider struct {
	Client *opensearch.Client
}

// Connect with retries and integrate into global registry
func connectOpenSearch(raw json.RawMessage) (*OpenSearchProvider, error) {
	var cfg opensearchConfig
	if err := json.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, nil
	}

	var client *opensearch.Client
	var err error

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		client, err = opensearch.NewClient(opensearch.Config{
			Addresses: []string{cfg.URL},
			Transport: &http.Transport{
				MaxIdleConnsPerHost: cfg.MaxIdle,
				IdleConnTimeout:     time.Duration(cfg.Timeout) * time.Second,
			},
		})
		if err != nil {
			logs.Errorf(Ctx, "OpenSearch client creation failed: %v. Retrying...\n", err)
			time.Sleep(2 * time.Second)
			continue
		}

		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, pingErr := client.Cluster.Health()
		if pingErr == nil {
			defer res.Body.Close()
			logs.Info(Ctx, "OpenSearch connected ✅")
			return &OpenSearchProvider{Client: client}, nil
		}

		logs.Warningf(Ctx, "OpenSearch not ready: %v. Retrying...\n", pingErr)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("unable to connect to OpenSearch at %s: %v", cfg.URL, err)
}

// HealthCheck implements DatabaseProvider interface
func (o *OpenSearchProvider) HealthCheck(ctx context.Context) error {
	_, err := o.Client.Cluster.Health()
	return err
}

// Close implements DatabaseProvider interface
func (o *OpenSearchProvider) Close() error {
	// OpenSearch client does not require close, but implement if needed
	return nil
}
