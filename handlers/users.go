package handlers

import (
	auth "Adornme/Auth"
	user "Adornme/controllers/users"
	"Adornme/logging"
	"Adornme/models"
	"Adornme/restapi/operations/users"
	"context"
	"fmt"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var logs = logging.Component("restapi")

func RegisterUser(params users.RegisterUserParams) middleware.Responder {
	// Generate a unique request ID
	requestID := uuid.New().String()
	ctx := logging.WithRequestID(context.Background(), requestID)

	logs.Infof(ctx, "Register User Called: %v", params.Body)

	// TODO: Save user to DB and get real userID
	userID := uuid.New().String()

	u := user.NewUser(userID, "en", userID, "My-Service")
	authResponse, err := u.RegisterUser(ctx, params.Body)

	if err != nil {
		return users.NewRegisterUserBadRequest().WithPayload(err)
	}

	return users.NewRegisterUserCreated().WithPayload(authResponse)
}

func GetUserProfile(params users.GetUserProfileParams, principle *models.Principal) middleware.Responder {
	// Extract token from Authorization header
	authHeader := params.HTTPRequest.Header.Get("Authorization")
	if authHeader == "" {
		msg := "missing authorization header"
		return users.NewGetUserProfileUnauthorized().WithPayload(&models.ErrorResponse{
			Error: &msg,
		})
	}

	userID, err := auth.ValidateAccessToken(authHeader)
	if err != nil {
		msg := "invalid or expired token"
		return users.NewGetUserProfileUnauthorized().WithPayload(&models.ErrorResponse{
			Error: &msg,
		})
	}

	// Generate a request ID for logging/tracing
	requestID := uuid.New().String()
	ctx := logging.WithRequestID(context.Background(), requestID)

	logs.Infof(ctx, "GetUserProfile called for userID: %s", userID)

	// Initialize user controller
	u := user.NewUser(requestID, "en", requestID, "My-Service")

	// Convert string userID to int64
	userIDI, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		msg := "invalid user ID"
		return users.NewGetUserProfileUnauthorized().WithPayload(&models.ErrorResponse{
			Error: &msg,
		})
	}
	// Fetch user details from DB
	userData, errResp := u.GetUser(ctx, userIDI)
	if errResp != nil {
		return users.NewGetUserProfileUnauthorized().WithPayload(errResp)
	}

	return users.NewGetUserProfileOK().WithPayload(userData)
}

func LoginUser(params users.LoginUserParams, principal *models.Principal) middleware.Responder {
	// Generate a request ID for logging/tracing
	requestID := uuid.New().String()
	ctx := logging.WithRequestID(context.Background(), requestID)

	// Initialize user controller
	u := user.NewUser(requestID, "en", requestID, "My-Service")

	// 1️⃣ Fetch user from DB
	dbUser, err := u.GetUserByEmail(ctx, *params.Body.Email)
	if err != nil {
		msg := "invalid email or password"
		return users.NewLoginUserUnauthorized().WithPayload(&models.ErrorResponse{Error: &msg})
	}

	// 2️⃣ Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(*params.Body.Password)); err != nil {
		msg := "invalid email or password"
		return users.NewLoginUserUnauthorized().WithPayload(&models.ErrorResponse{Error: &msg})
	}

	userID := fmt.Sprintf("%d", dbUser.ID)

	// 3️⃣ Generate JWT tokens
	accessToken, err1 := auth.GenerateToken(userID, 24) // 24h access token
	if err1 != nil {
		msg := "failed to generate access token"
		return users.NewLoginUserUnauthorized().WithPayload(&models.ErrorResponse{Error: &msg})
	}

	refreshToken, err2 := auth.GenerateRefreshToken(userID, 24*7) // 7 days refresh token
	if err2 != nil {
		msg := "failed to generate refresh token"
		return users.NewLoginUserUnauthorized().WithPayload(&models.ErrorResponse{Error: &msg})
	}

	// // 4️⃣ Save refresh token in DB
	// if err := u.DB.UpdateRefreshToken(ctx, int(dbUser.ID), refreshToken); err != nil {
	// 	logs.Errorf(ctx, "failed to save refresh token for user %d: %v", dbUser.ID, err)
	// 	// Not blocking login; token returned to user
	// }

	// 5️⃣ Build response
	emailStr := strfmt.Email(dbUser.Email)
	return users.NewLoginUserOK().WithPayload(&models.AuthResponse{
		User: &models.User{
			ID:    &dbUser.ID,
			Name:  &dbUser.Name,
			Email: &emailStr,
			Phone: dbUser.PhoneNumber,
		},
		Token:        &accessToken,
		RefreshToken: refreshToken,
	})
}
