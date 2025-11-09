package handlers

import (
	database "Adornme/databases"
	"Adornme/logging"
	"Adornme/restapi/operations/system"
	"Adornme/utils"
	"context"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
)

// GetHealth handles the GET /health endpoint
func GetHealth(params system.GetHealthParams) middleware.Responder {
	requestID := uuid.New().String()
	ctx := logging.WithRequestID(context.Background(), requestID)
	logs.Infof(ctx, "Starting system health check at %v", time.Now())

	dependencies := make(map[string]string)
	overallHealthy := true

	// Iterate through all registered dependencies in database.Do
	for name, db := range database.Do {
		depCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := utils.Retry(3, 2*time.Second, func() error {
			return db.HealthCheck(depCtx)
		})
		cancel()

		if err != nil {
			logs.Warningf(ctx, "%s health check failed: %v", name, err)
			dependencies[name] = "unhealthy ❌"
			overallHealthy = false
		} else {
			logs.Infof(ctx, "%s is healthy ✅", name)
			dependencies[name] = "healthy ✅"
		}
	}

	// Determine overall system status
	status := "ok"
	description := "All dependencies are healthy and running smoothly"
	if !overallHealthy {
		status = "degraded"
		description = "Some dependencies are unhealthy — check logs for details"
	}

	// Build final JSON response
	resp := &system.GetHealthOKBody{
		Status:       status,
		Description:  description,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		Dependencies: dependencies,
	}

	logs.Infof(ctx, "Health check completed: %+v", resp)

	return system.NewGetHealthOK().WithPayload(resp)
}
