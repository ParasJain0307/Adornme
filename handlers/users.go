package handlers

import (
	auth "Adornme/Auth"
	user "Adornme/controllers/users"
	"Adornme/logging"
	"Adornme/models"
	"Adornme/restapi/operations/users"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
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

func GetUserProfile(params users.GetUserProfileParams, principal *models.Principal) middleware.Responder {
	// Step 1: Extract Authorization header
	authHeader := params.HTTPRequest.Header.Get("Authorization")
	if authHeader == "" {
		msg := "missing authorization header"
		logs.Errorf(context.Background(), "AUTH ERROR: %s", msg)

		return users.NewGetUserProfileUnauthorized().WithPayload(&models.ErrorResponse{
			Error: &msg,
		})
	}

	// Step 2: Validate JWT token and extract userID
	userID, err := auth.ValidateAccessToken(authHeader)
	if err != nil {
		// Log actual error for debugging (IMPORTANT)
		logs.Errorf(context.Background(), "JWT VALIDATION FAILED: %v, header: %s", err, authHeader)

		msg := "invalid or expired token"
		return users.NewGetUserProfileUnauthorized().WithPayload(&models.ErrorResponse{
			Error: &msg,
		})
	}

	// Step 3: Generate request ID for tracing
	requestID := uuid.New().String()
	ctx := logging.WithRequestID(context.Background(), requestID)

	logs.Infof(ctx, "GetUserProfile started for userID: %s", userID)

	// Step 4: Initialize user controller/service
	u := user.NewUser(requestID, "en", requestID, "My-Service")

	// Step 5: Convert userID (string → int64)
	userIDI, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		logs.Errorf(ctx, "USER ID PARSE FAILED: userID=%s, error=%v", userID, err)

		msg := "invalid user ID"
		return users.NewGetUserProfileUnauthorized().WithPayload(&models.ErrorResponse{
			Error: &msg,
		})
	}

	// Step 6: Fetch user details from DB
	userData, errResp := u.GetUser(ctx, userIDI)
	if errResp != nil {
		// Log DB/service layer error
		logs.Errorf(ctx, "GET USER FAILED: userID=%d, error=%v", userIDI, errResp)

		return users.NewGetUserProfileUnauthorized().WithPayload(errResp)
	}

	// Step 7: Success response
	logs.Infof(ctx, "GetUserProfile success for userID: %d", userIDI)

	return users.NewGetUserProfileOK().WithPayload(userData)
}

func LoginUser(params users.LoginUserParams) middleware.Responder {

	requestID := uuid.New().String()
	ctx := logging.WithRequestID(context.Background(), requestID)

	u := user.NewUser(requestID, "en", requestID, "My-Service")

	resp, err := u.Login(ctx, params.Body.Email, *params.Body.Password)
	if err != nil {
		msg := err.Error()
		return users.NewLoginUserUnauthorized().
			WithPayload(&models.ErrorResponse{Error: &msg})
	}

	return users.NewLoginUserOK().WithPayload(resp)
}

func RefreshToken(params users.RefreshTokenParams) middleware.Responder {

	// 1️⃣ Request ID + context
	requestID := uuid.New().String()
	ctx := logging.WithRequestID(context.Background(), requestID)

	logs.Infof(ctx, "RefreshToken called")

	// 2️⃣ Validate request body
	if params.Body == nil || params.Body.RefreshToken == nil {
		msg := "refresh token is required"
		return users.NewRefreshTokenBadRequest().WithPayload(&models.ErrorResponse{
			Error: &msg,
		})
	}

	// 3️⃣ Initialize user service
	u := user.NewUser(requestID, "en", requestID, "My-Service")

	// 4️⃣ Call service layer (ALL logic handled there)
	tokenResp, errResp := u.RefreshToken(ctx, *params.Body.RefreshToken)
	if errResp != nil {
		logs.Errorf(ctx, "REFRESH TOKEN FAILED: %+v", errResp)

		return users.NewGetUserProfileUnauthorized().WithPayload(errResp)
	}

	// 5️⃣ Success response
	logs.Infof(ctx, "RefreshToken success")

	return users.NewRefreshTokenOK().WithPayload(tokenResp)
}

func LogoutUser(params users.LogoutUserParams) middleware.Responder {

	// 🔹 Request context + logging
	requestID := uuid.New().String()
	ctx := logging.WithRequestID(context.Background(), requestID)

	u := user.NewUser(requestID, "en", requestID, "My-Service")

	// 🔹 1. Extract refresh token
	authHeader := params.HTTPRequest.Header.Get("Authorization")
	if authHeader == "" {
		msg := "missing authorization header"
		return users.NewLogoutUserUnauthorized().
			WithPayload(&models.ErrorResponse{Error: &msg})
	}

	// 🔹 2. Remove "Bearer "
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		msg := "invalid authorization format"
		return users.NewLoginUserUnauthorized().
			WithPayload(&models.ErrorResponse{Error: &msg})
	}

	refreshToken := parts[1]
	fmt.Printf("Extracted refresh token: %s\n", refreshToken)

	// 🔹 3. Call service layer
	err := u.Logout(ctx, refreshToken)
	if err != nil {
		msg := err.Error()
		return users.NewLogoutUserUnauthorized().
			WithPayload(&models.ErrorResponse{Error: &msg})
	}

	// 🔹 4. Success response
	success := "logged out successfully"
	return users.NewLogoutUserOK().WithPayload(&models.SuccessResponse{Message: &success})
}

func ForgetPassword(params users.ForgetPasswordParams) middleware.Responder {

	// 🔹 Request context + logging
	requestID := uuid.New().String()
	ctx := logging.WithRequestID(context.Background(), requestID)

	u := user.NewUser(requestID, "en", requestID, "My-Service")

	// 🔹 Log request start
	logs.Info(ctx, "ForgetPassword request received")

	// 🔹 Validate input
	if params.Body.Email == nil {
		msg := "email is required"
		logs.Error(ctx, "missing email in request")
		return users.NewForgetPasswordBadRequest().
			WithPayload(&models.ErrorResponse{Error: &msg})
	}

	email := *params.Body.Email
	logs.Info(ctx, "processing forgot password", "email", email)

	// 🔹 Call service layer
	err := u.ForgetPassword(ctx, email)
	if err != nil {
		msg := "failed to process request"

		logs.Error(ctx, "ForgetPassword service failed",
			"email", email,
			"error", err.Error(),
		)

		return users.NewGetUserProfileUnauthorized().
			WithPayload(&models.ErrorResponse{Error: &msg})
	}

	// 🔹 Success log
	logs.Info(ctx, "ForgetPassword email sent (if user exists)", "email", email)

	success := "If the email exists, a reset link has been sent"
	return users.NewForgetPasswordOK().
		WithPayload(&users.ForgetPasswordOKBody{Message: success})
}

func IdentifyUser(params users.IdentifyUserParams) middleware.Responder {
	// 🔹 Request context + logging
	requestID := uuid.New().String()
	ctx := logging.WithRequestID(context.Background(), requestID)

	u := user.NewUser(requestID, "en", requestID, "My-Service")

	// 🔹 Log request start
	logs.Info(ctx, "IdentifyUser request received")

	// 🔹 Validate input

	if params.Body.Identifier == nil || strings.TrimSpace(*params.Body.Identifier) == "" {
		msg := "invalid identifier format"
		return users.NewIdentifyUserBadRequest().WithPayload(&models.ErrorResponse{
			Error: &msg,
		})
	}

	identifier := strings.TrimSpace(*params.Body.Identifier)

	// ✅ Call controller
	err := u.IdentifyUser(params.HTTPRequest.Context(), identifier)
	if err != nil {
		log.Printf("Identify error: %v", err)

		msg := "invalid identifier"
		return users.NewIdentifyUserBadRequest().WithPayload(&models.ErrorResponse{
			Error: &msg,
		})
	}

	return users.NewIdentifyUserOK()

}

func SendOTP(params users.OTPLoginParams) middleware.Responder {
	// 🔹 Request context + logging
	requestID := uuid.New().String()
	ctx := logging.WithRequestID(context.Background(), requestID)

	u := user.NewUser(requestID, "en", requestID, "My-Service")

	// 🔹 Log request start
	logs.Info(ctx, "SendOTP request received")

	// 🔹 Validate input
	if params.Body.Identifier == nil || strings.TrimSpace(*params.Body.Identifier) == "" {
		msg := "invalid identifier format"
		return users.NewOTPLoginBadRequest().WithPayload(&models.ErrorResponse{
			Error: &msg,
		})
	}

	identifier := strings.TrimSpace(*params.Body.Identifier)

	// ✅ Call controller
	err := u.SendOTP(params.HTTPRequest.Context(), identifier)
	if err != nil {
		log.Printf("Send OTP error: %v", err)

		msg := "failed to send OTP"
		return users.NewOTPLoginBadRequest().WithPayload(&models.ErrorResponse{
			Error: &msg,
		})
	}

	return users.NewOTPLoginOK().WithPayload(&models.GenericResponse{Message: "OTP sent successfully"})
}
