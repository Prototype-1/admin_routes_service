package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
)

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY")) 

func ParseJWT(tokenStr string) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecretKey, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	adminID := uint(claims["admin_id"].(float64))
	return adminID, nil
}
