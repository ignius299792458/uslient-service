package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"had-service/config"
	"time"
)

func GenerateToken(username string, id string) (string, error) {
	envConfig := config.EnvLoad()
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = id
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix() // week expircy token

	tokenString, err := token.SignedString([]byte(envConfig.JwtScretKey))
	if err != nil {
		return "", errors.New("failed to sign token: " + err.Error())
	}

	return tokenString, nil
}
