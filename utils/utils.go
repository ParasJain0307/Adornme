package utils

import (
	"log"
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
