package users

import (
	Auth "Adornme/Auth"
	db "Adornme/databases"
	"Adornme/logging"
	"Adornme/models"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"golang.org/x/crypto/bcrypt"
)

var logs = logging.Component("users")

func (u *User) RegisterUser(ctx context.Context, params *models.RegisterRequest) (*models.AuthResponse, *models.ErrorResponse) {
	logs.Infof(ctx, "Register User called with requestID: %s", u.RequestID)

	// Extract and trim request fields
	email := ""
	name := ""
	phone := ""

	if params.Email != nil {
		email = strings.TrimSpace(params.Email.String())
	}
	if params.Name != nil {
		name = strings.TrimSpace(*params.Name)
	}
	if params.Phone != "" {
		phone = strings.TrimSpace(params.Phone)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*params.Password), bcrypt.DefaultCost)
	if err != nil {
		msg := "failed to hash password"
		return nil, &models.ErrorResponse{Error: &msg}
	}

	// Build DB user
	dbUser := &db.User{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Email:       email,
		Name:        name,
		PhoneNumber: phone,
		Password:    string(hashedPassword),
	}

	// Create user in DB
	id, err := u.DB.CreateUser(ctx, dbUser)
	if err != nil {
		msg := err.Error()
		return nil, &models.ErrorResponse{Error: &msg}
	}

	userID := fmt.Sprintf("%d", id) // convert to string for JWT

	// Generate JWT tokens
	accessToken, err := Auth.GenerateToken(userID, 24) // 24h expiry
	if err != nil {
		msg := "failed to generate access token"
		return nil, &models.ErrorResponse{Error: &msg}
	}

	refreshToken, err := Auth.GenerateRefreshToken(userID, 24*7) // 7-day expiry
	if err != nil {
		msg := "failed to generate refresh token"
		return nil, &models.ErrorResponse{Error: &msg}
	}

	dbUser.RefreshToken = refreshToken
	if err := u.DB.UpdateUserTokens(ctx, id, refreshToken); err != nil {
		logs.Errorf(ctx, "failed to save tokens for user %d: %v", id, err)
	}

	// Build API response
	Id := int64(id)
	emailStr := strfmt.Email(email)
	return &models.AuthResponse{
		User: &models.User{
			ID:    &Id,
			Name:  &name,
			Phone: phone,
			Email: &emailStr,
		},
		Token:        &accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GetUser retrieves a user by ID
func (u *User) GetUser(ctx context.Context, userID int64) (*models.User, *models.ErrorResponse) {
	logs.Infof(ctx, "GetUser called with requestID: %s, userID: %d", u.RequestID, userID)

	// Fetch user from DB
	dbUser, err := u.DB.GetUser(ctx, int(userID))
	if err != nil {
		msg := err.Error()
		return nil, &models.ErrorResponse{Error: &msg}
	}

	var email *strfmt.Email
	if dbUser.Email != "" {
		e := strfmt.Email(dbUser.Email)
		email = &e
	}
	if dbUser == nil {
		msg := "user not found"
		return nil, &models.ErrorResponse{Error: &msg}
	}

	// Map DB user → API model
	return &models.User{
		ID:    &dbUser.ID,
		Name:  &dbUser.Name,
		Email: email,
		Phone: dbUser.PhoneNumber,
	}, nil
}

func (u *User) GetUserByEmail(ctx context.Context, Email strfmt.Email) (*db.User, *models.ErrorResponse) {
	logs.Infof(ctx, "GetUser called with requestID: %s, Email: %d", u.RequestID, Email)

	// Fetch user from DB
	dbUser, err := u.DB.GetUserByEmail(ctx, string(Email))
	if err != nil {
		msg := err.Error()
		return nil, &models.ErrorResponse{Error: &msg}
	}
	return dbUser, nil
}
