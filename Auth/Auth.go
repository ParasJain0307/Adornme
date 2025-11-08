package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-very-secure-secret-key")
var refreshSecret = []byte("your-refresh-token-secret-key")

type AuthClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken generates an access JWT token
func GenerateToken(userID string, expiryHours int) (string, error) {
	if userID == "" {
		return "", errors.New("userID cannot be empty")
	}

	claims := AuthClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Adornme",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

// GenerateRefreshToken generates a long-lived refresh token
func GenerateRefreshToken(userID string, expiryDays int) (string, error) {
	claims := AuthClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryDays) * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Adornme",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

// ValidateAccessToken validates an access JWT token
func ValidateAccessToken(tokenString string) (string, error) {
	return validateToken(tokenString, jwtSecret)
}

// ValidateRefreshToken validates a refresh token
func ValidateRefreshToken(tokenString string) (string, error) {
	return validateToken(tokenString, refreshSecret)
}

// validateToken validates a JWT and returns the user ID if it's valid.
// Returns detailed errors for easier debugging.
func validateToken(tokenString string, secret []byte) (string, error) {
	// Clean and normalize token string
	tokenString = strings.TrimPrefix(strings.TrimSpace(tokenString), "Bearer ")

	if tokenString == "" {
		return "", errors.New("token is empty")
	}

	// Parse the JWT using the provided secret key
	parsedToken, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token uses HMAC signing
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	// Handle parsing errors
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return "", errors.New("invalid token signature")
		}
		if strings.Contains(err.Error(), "token is expired") {
			return "", errors.New("token has expired")
		}
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	// Validate claims
	if claims, ok := parsedToken.Claims.(*AuthClaims); ok {
		if !parsedToken.Valid {
			return "", errors.New("token is invalid")
		}
		if claims.UserID == "" {
			return "", errors.New("token does not contain a user ID")
		}
		return claims.UserID, nil
	}

	return "", errors.New("token claims are invalid or malformed")
}
