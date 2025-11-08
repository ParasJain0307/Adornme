package database

import (
	"Adornme/utils"
	"context"
	"log"
	"time"
)

func StartDBHealthTicker(interval time.Duration, stopCh <-chan struct{}) {

	
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for name, db := range Do {
				ctx, cancel := context.WithTimeout(context.Background(), interval/2)
				err := utils.Retry(3, 5*time.Second, func() error {
					return db.HealthCheck(ctx)
				})
				cancel()
				if err != nil {
					logs.Warningf(Ctx, "%s health check failed: %v", name, err)
				} else {
					logs.Infof(Ctx, "%s healthy ✅", name)
				}
			}
		case <-stopCh:
			log.Println("Stopping DB health checker")
			return
		}
	}
}
