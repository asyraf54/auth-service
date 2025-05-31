package utils

import (
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokens(userID int) (accessToken string, refreshToken string, err error) {
	// Access token claims (expires in 24 hours)
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"type":    "access", 
	}
	accessJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessJwt.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", "", err
	}

	// Refresh token claims (expires in 7 days)
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"type":    "refresh", 
	}
	refreshJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshJwt.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ValidateToken(tokenStr string, expectedType string) (map[string]interface{}, error) {
	// Parse the token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// Token parsing or signature validation failed
	if err != nil || !token.Valid {
		return nil, err
	}

	// Extract and validate claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKeyType
	}

	// Check for correct token type
	if tokenType, ok := claims["type"].(string); !ok || tokenType != expectedType {
		return nil, jwt.ErrInvalidType
	}

	// Check for expiration
	if exp, ok := claims["exp"].(float64); !ok || int64(exp) < time.Now().Unix() {
		return nil, jwt.ErrTokenExpired
	}

	return claims, nil
}
