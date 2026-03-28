package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"Adornme/config"

	"github.com/golang-jwt/jwt/v5"
)

var cfg = config.LoadConfig()

var jwtSecret = []byte(cfg.JWTSecret)
var refreshSecret = []byte(cfg.RefreshSecret)

type AuthClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// 🔹 Generate Access Token
func GenerateToken(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("userID cannot be empty")
	}

	claims := AuthClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.AccessTokenExpiryHours) * time.Hour)), // ✅ FIXED
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Adornme.in",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// 🔹 Generate Refresh Token
func GenerateRefreshToken(userID string) (string, error) {
	claims := AuthClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.RefreshTokenExpiryDays) * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Adornme",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

// 🔹 Validate Access Token
func ValidateAccessToken(tokenString string) (string, error) {
	token := extractToken(tokenString)
	return validateToken(token, jwtSecret)
}

// 🔹 Validate Refresh Token
func ValidateRefreshToken(tokenString string) (string, error) {
	return validateToken(strings.TrimSpace(tokenString), refreshSecret)
}

// 🔹 Core Validation
func validateToken(tokenString string, secret []byte) (string, error) {
	if tokenString == "" {
		return "", errors.New("token is empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return "", errors.New("invalid token signature")
		}
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", errors.New("token expired")
		}
		return "", err
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	if claims.UserID == "" {
		return "", errors.New("user_id missing")
	}

	fmt.Println(claims.UserID)
	return claims.UserID, nil
}

// 🔹 Extract Bearer Token
func extractToken(header string) string {
	parts := strings.Fields(header)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}
