package utils

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"had-service/config"
)

func ValidateToken(tokenString string) (interface{}, error) {
	envConfig := config.EnvLoad()
	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(envConfig.JwtScretKey), nil
	})
	if err != nil {
		return nil, errors.New("failed to parse token: " + err.Error())
	}

	// Check if token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
