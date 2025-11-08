package users

import (
	db "Adornme/databases"
	"Adornme/models"
	model "Adornme/models"
	"context"

	"github.com/go-openapi/strfmt"
)

// User struct holds request-related metadata for tracking
type User struct {
	RequestID   string
	InstanceID  string
	ServiceName string
	AcceptLang  string
	DB          db.PostgresProvider
}

// Users interface defines user-related operations
type Users interface {
	RegisterUser(ctx context.Context, params *models.RegisterRequest) (*model.AuthResponse, *model.ErrorResponse)
	GetUser(ctx context.Context, userID int64) (*models.User, *models.ErrorResponse)
	GetUserByEmail(ctx context.Context, Email strfmt.Email) (*db.User, *models.ErrorResponse)
}

// NewUser initializes a User instance with request metadata
func NewUser(reqID, acceptLang, instanceID, serviceName string) Users {
	// Get the postgres client from registry
	pgAny := db.Do["postgres"]

	// Type assert to PostgresClients
	pgClients, ok := pgAny.(*db.PostgresClients)
	if !ok {
		panic("postgres client not initialized properly")
	}

	return &User{
		RequestID:   reqID,
		InstanceID:  instanceID,
		ServiceName: serviceName,
		AcceptLang:  acceptLang,
		DB:          *pgClients.UsersDB, // ✅ inject UsersDB
	}
}
