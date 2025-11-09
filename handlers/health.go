package handlers

import (
	"Adornme/logging"
	"Adornme/restapi/operations/system"
	"context"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
)

// Health handles the GET /health endpoint
func GetHealth(params system.GetHealthParams) middleware.Responder {
	// Generate a unique request ID
	requestID := uuid.New().String()
	ctx := logging.WithRequestID(context.Background(), requestID)

	logs.Infof(ctx, "Checking system health at %v", time.Now())

	resp := &system.GetHealthOKBody{
		Description: "System is running and healthy",
		Status:      "ok",
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	return system.NewGetHealthOK().WithPayload(resp)
}
