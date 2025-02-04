 package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
)

var jwtSecretKey []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
    key := os.Getenv("JWT_SECRET_KEY")
    if key == "" {
        log.Fatal("JWT_SECRET_KEY is not set")
    }
    jwtSecretKey = []byte(key)
    fmt.Println("Loaded JWT Secret Key:", jwtSecretKey)
}

func ParseJWT(tokenStr string) (uint, error) {
	if len(jwtSecretKey) == 0 {
		return 0, errors.New("JWT_SECRET_KEY is not set")
	}

	log.Println("Using JWT Secret Key:", string(jwtSecretKey)) 

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("Invalid signing method:", token.Header["alg"])
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})
	if err != nil {
		log.Println("JWT Parsing Error:", err)
		return 0, fmt.Errorf("failed to parse JWT: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Println("Invalid token claims:", claims)
		return 0, errors.New("invalid token claims")
	}

	log.Println("Extracted Claims:", claims)

	adminID, err := extractAdminID(claims)
	if err != nil {
		log.Println("Error extracting admin_id:", err)
		return 0, err
	}

	log.Println("Extracted admin_id:", adminID)
	return adminID, nil
}

func extractAdminID(claims jwt.MapClaims) (uint, error) {
	rawAdminID, exists := claims["admin_id"]
	if !exists {
		return 0, errors.New("missing admin_id in token claims")
	}

	switch v := rawAdminID.(type) {
	case float64:
		return uint(v), nil
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("invalid admin_id format: %v", v)
		}
		return uint(id), nil
	default:
		return 0, fmt.Errorf("unsupported admin_id type: %T", v)
	}
}