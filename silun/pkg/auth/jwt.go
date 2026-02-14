package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret  = []byte("your-access-secret-key")
	refreshSecret = []byte("your-refresh-secret-key")
)

func GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp": time.Now().Add(2 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(claims)
	return token.SignedString(accessSecret)
}

func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(claims)
	return token.SignedString(refreshSecret)
}

func ParseAccessToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", err
	}
	
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}
	
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("invalid user_id in token")
	}
	
	return userID, nil
}

func ParseRefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", err
	}
	
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}
	
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("invalid user_id in token")
	}
	
	return userID, nil
}
