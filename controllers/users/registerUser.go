package users

import (
	Auth "Adornme/Auth"
	auth "Adornme/Auth"
	db "Adornme/databases"
	"Adornme/logging"
	"Adornme/models"
	"Adornme/utils"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
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
		logs.Errorf(ctx, "failed to create user: %v", err)
		return nil, &models.ErrorResponse{Error: &msg}
	}

	userID := fmt.Sprintf("%d", id) // convert to string for JWT

	// Generate JWT tokens
	accessToken, err := Auth.GenerateToken(userID) // 24h expiry
	if err != nil {
		msg := "failed to generate access token"
		return nil, &models.ErrorResponse{Error: &msg}
	}

	refreshToken, err := Auth.GenerateRefreshToken(userID) // 7-day expiry
	if err != nil {
		msg := "failed to generate refresh token"
		return nil, &models.ErrorResponse{Error: &msg}
	}

	dbUser.RefreshToken = refreshToken
	if err := u.DB.UpdateUserTokens(ctx, id, refreshToken, accessToken); err != nil {
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
	logs.Infof(ctx, "GetUser called with requestID: %s, Email: %s", u.RequestID, Email)

	// Fetch user from DB
	dbUser, err := u.DB.GetUserByEmail(ctx, string(Email))
	if err != nil {
		msg := err.Error()
		return nil, &models.ErrorResponse{Error: &msg}
	}
	return dbUser, nil
}

func (u *User) RefreshToken(ctx context.Context, token string) (*models.AuthResponse, *models.ErrorResponse) {
	logs.Infof(ctx, "RefreshToken called with requestID: %s:---->%s", u.RequestID, token)

	// 1️⃣ Validate input
	token = strings.TrimSpace(token)
	if token == "" {
		msg := "refresh token is required"
		return nil, &models.ErrorResponse{Error: &msg}
	}

	// 2️⃣ Validate JWT
	userId, err := auth.ValidateRefreshToken(token)
	if err != nil {
		logs.Errorf(ctx, "INVALID REFRESH TOKEN: %v", err)

		msg := "invalid or expired token"
		return nil, &models.ErrorResponse{Error: &msg}
	}

	logs.Infof(ctx, "RefreshToken validated for userID: %s", userId)

	// 3️⃣ Convert userID → int
	userIDI, err := strconv.Atoi(userId)
	if err != nil {
		logs.Errorf(ctx, "INVALID USER ID FORMAT: %v", err)

		msg := "invalid user id"
		return nil, &models.ErrorResponse{Error: &msg}
	}

	// 4️⃣ Fetch token from DB
	storedToken, err := u.DB.GetRefreshToken(ctx, userId)
	if err != nil {
		logs.Errorf(ctx, "FAILED TO FETCH REFRESH TOKEN: userID=%d, err=%v", userIDI, err)

		msg := "token not found"
		return nil, &models.ErrorResponse{Error: &msg}
	}

	// 5️⃣ Compare tokens
	if storedToken != token {
		logs.Errorf(ctx, "TOKEN MISMATCH: userID=%d", userIDI)

		msg := "invalid refresh token"
		return nil, &models.ErrorResponse{Error: &msg}
	}

	// 6️⃣ Generate new tokens
	newAccessToken, err := auth.GenerateToken(userId)
	if err != nil {
		logs.Errorf(ctx, "FAILED TO GENERATE ACCESS TOKEN: userID=%s, err=%v", userId, err)

		msg := "failed to generate access token"
		return nil, &models.ErrorResponse{Error: &msg}
	}

	newRefreshToken, err := auth.GenerateRefreshToken(userId)
	if err != nil {
		logs.Errorf(ctx, "FAILED TO GENERATE REFRESH TOKEN: userID=%s, err=%v", userId, err)

		msg := "failed to generate refresh token"
		return nil, &models.ErrorResponse{Error: &msg}
	}

	// 7️⃣ Update DB (rotation)
	if err := u.DB.UpdateUserTokens(ctx, userIDI, newRefreshToken, newAccessToken); err != nil {
		logs.Errorf(ctx, "FAILED TO UPDATE REFRESH TOKEN: userID=%d, err=%v", userIDI, err)

		msg := "failed to update refresh token"
		return nil, &models.ErrorResponse{Error: &msg}
	}

	logs.Infof(ctx, "RefreshToken success for userID: %s", userId)

	// 8️⃣ Response (match your login/register format)
	return &models.AuthResponse{
		Token:        &newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (u *User) Login(ctx context.Context, email *strfmt.Email, password string) (*models.AuthResponse, error) {

	// 1. Fetch user
	dbUser, err := u.DB.GetUserByEmail(ctx, string(*email))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// 2. Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}
	userID := fmt.Sprintf("%d", dbUser.ID)
	// 3. Generate tokens
	accessToken, err := auth.GenerateToken(userID)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := auth.GenerateRefreshToken(userID)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// 4. Save refresh token (IMPORTANT)
	err = u.DB.UpdateUserTokens(ctx, int(dbUser.ID), refreshToken, accessToken)
	if err != nil {
		return nil, errors.New("failed to store refresh token")
	}

	// 5. Build response
	emailStr := strfmt.Email(dbUser.Email)

	return &models.AuthResponse{
		User: &models.User{
			ID:    &dbUser.ID,
			Name:  &dbUser.Name,
			Email: &emailStr,
			Phone: dbUser.PhoneNumber,
		},
		Token:        &accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *User) Logout(ctx context.Context, refreshToken string) error {

	// 🔹 1. Validate refresh token (JWT)
	userID, err := auth.ValidateRefreshToken(refreshToken)
	if err != nil {
		logs.Errorf(ctx, "logout failed: invalid refresh token | err=%v", err)
		return errors.New("invalid or expired refresh token")
	}

	logs.Infof(ctx, "logout initiated for user_id=%s", userID)

	// 🔹 2. Get stored token from DB
	storedToken, err := u.DB.GetRefreshToken(ctx, userID)
	if err != nil {
		logs.Errorf(ctx, "logout failed: session not found | user_id=%s err=%v", userID, err)
		return errors.New("session not found")
	}

	// 🔹 3. Match token (CRITICAL SECURITY CHECK)
	if storedToken != refreshToken {
		logs.Warning(ctx, "logout failed: token mismatch | user_id=%s", userID)
		return errors.New("invalid session or token mismatch")
	}

	// 🔹 4. Delete / invalidate refresh token
	err = u.DB.DeleteRefreshToken(ctx, userID)
	if err != nil {
		logs.Errorf(ctx, "logout failed: DB error while deleting token | user_id=%s err=%v", userID, err)
		return errors.New("failed to logout")
	}

	logs.Infof(ctx, "logout successful | user_id=%s", userID)

	return nil
}

func (u *User) ForgetPassword(ctx context.Context, email string) error {

	// 🔹 Log start
	logs.Info(ctx, "ForgetPassword service started", "email", email)

	// 🔹 1. Get user by email
	user, err := u.DB.GetUserByEmail(ctx, email)
	if err != nil {
		// ⚠️ Do NOT expose user existence
		logs.Info(ctx, "user not found for email (safe ignore)", "email", email)
		return nil
	}

	logs.Info(ctx, "user found", "user_id", user.ID)

	// 🔹 2. Generate raw token
	rawToken := uuid.New().String()

	// 🔹 3. Hash token
	hash := sha256.Sum256([]byte(rawToken))
	hashedToken := hex.EncodeToString(hash[:])

	// 🔹 4. Expiry (15 min)
	expiry := time.Now().Add(15 * time.Minute)

	userId := fmt.Sprintf("%d", user.ID)

	// 🔹 5. Save token in DB
	err = u.DB.SaveResetToken(ctx, userId, hashedToken, expiry)
	if err != nil {
		logs.Error(ctx, "failed to save reset token",
			"user_id", user.ID,
			"error", err.Error(),
		)
		return err
	}

	logs.Info(ctx, "reset token saved",
		"user_id", user.ID,
		"expiry", expiry.String(),
	)

	// 🔹 6. Create reset link
	resetLink := fmt.Sprintf(
		"http://localhost:3000/reset-password?token=%s",
		rawToken,
	)

	logs.Info(ctx, "reset link generated")

	// 🔹 7. Send email (async)
	go func() {
		err := utils.SendResetEmail(user.Email, user.Name, resetLink)
		if err != nil {
			logs.Error(ctx, "failed to send reset email",
				"email", user.Email,
				"error", err.Error(),
			)
			return
		}
		logs.Info(ctx, "reset email sent", "email", user.Email)
	}()

	return nil
}
