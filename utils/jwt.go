package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

var jwtSecretKey = []byte("RsaGPCEv4Oy1rRzVq8rT9vyPV29DZ4yQ4X_xd0QRK") 

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
