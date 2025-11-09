package handlers

import (
	ad "Adornme/cmd/app"
	database "Adornme/databases"
	"Adornme/logging"
	"Adornme/restapi/operations/system"
	"Adornme/utils"
	"context"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
)

// GetHealth handles the GET /health endpoint
func GetHealth(params system.GetHealthParams) middleware.Responder {
	requestID := uuid.New().String()
	ctx := logging.WithRequestID(context.Background(), requestID)
	startCheck := time.Now()
	logs.Infof(ctx, "Starting system health check at %v", startCheck)

	dependencies := make(map[string]system.GetHealthOKBodyDependenciesAnon)
	overallHealthy := true

	for name, db := range database.Do {
		depCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

		var (
			uptime     string
			latency    float32
			checkError error
		)

		// Check if dependency supports detailed health info
		if detailed, ok := db.(interface {
			HealthDetails(context.Context) (string, float64, error)
		}); ok {
			up, lat, err := detailed.HealthDetails(depCtx)
			uptime = up
			latency = float32(lat)
			checkError = err
		} else {
			checkError = utils.Retry(3, 2*time.Second, func() error {
				return db.HealthCheck(depCtx)
			})
		}

		cancel() // cancel context immediately

		lastChecked := strfmt.DateTime(time.Now().UTC())

		if checkError != nil {
			logs.Warningf(ctx, "%s health check failed: %v", name, checkError)
			dependencies[name] = system.GetHealthOKBodyDependenciesAnon{
				Status:       "unhealthy ❌",
				ErrorMessage: checkError.Error(),
				LastChecked:  lastChecked,
			}
			overallHealthy = false
		} else {
			logs.Infof(ctx, "%s is healthy ✅ (uptime: %s, latency: %.2fms)", name, uptime, latency)
			dependencies[name] = system.GetHealthOKBodyDependenciesAnon{
				Status:      "healthy ✅",
				Uptime:      uptime,
				LatencyMs:   latency,
				LastChecked: lastChecked,
			}
		}
	}

	// Overall system status
	status := "ok"
	description := "All dependencies are healthy and running smoothly"
	if !overallHealthy {
		status = "degraded"
		description = "Some dependencies are unhealthy — check logs for details"
	}

	// Calculate real application uptime from AppStartTime
	appUptime := time.Since(ad.AppStartTime).Truncate(time.Second).String()

	resp := &system.GetHealthOKBody{
		Status:      status,
		Description: description,
		Timestamp:   strfmt.DateTime(time.Now().UTC()),
		Uptime:      appUptime,
		Dependencies: func() map[string]system.GetHealthOKBodyDependenciesAnon {
			depMap := make(map[string]system.GetHealthOKBodyDependenciesAnon)
			for name, info := range dependencies {
				depMap[name] = info
			}
			return depMap
		}(),
	}

	logs.Infof(ctx, "Health check completed: %+v", resp)

	return system.NewGetHealthOK().WithPayload(resp)
}
