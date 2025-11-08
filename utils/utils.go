package utils

import (
	"context"
	"log"
	"os"
	"time"
)

func Retry(attempts int, delay time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		log.Printf("Attempt %d failed: %v. Retrying in %s...", i+1, err, delay)
		time.Sleep(delay)
	}
	return err
}

func GetPodID() string {
	if podID := os.Getenv("POD_ID"); podID != "" {
		return podID
	}
	if hostname, err := os.Hostname(); err == nil {
		return hostname
	}
	return "pod-unknown"
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, "requestID", requestID)
}

func RequestIDFromContext(ctx context.Context) string {
	if v := ctx.Value("requestID"); v != nil {
		return v.(string)
	}
	return ""
}
