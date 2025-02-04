package middleware

import (
    "context"
    "errors"
    "fmt"
    "strings"
    "github.com/golang-jwt/jwt/v4"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
    "github.com/Prototype-1/admin_routes_service/utils"
    "log"
    "time"
)

type contextKey string

const adminIDKey contextKey = "admin_id"

func VerifyJWT(tokenString string, secretKey string) (jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secretKey), nil
    })

    if err != nil {
        return nil, fmt.Errorf("invalid JWT token: %w", err)
    }
    if !token.Valid {
        return nil, fmt.Errorf("invalid JWT token: signature is invalid")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, errors.New("invalid token claims")
    }

    if exp, ok := claims["exp"].(float64); ok {
        if time.Now().Unix() > int64(exp) {
            return nil, errors.New("token has expired")
        }
    } else {
        return nil, errors.New("invalid expiration claim")
    }

    return claims, nil
}

func JWTAuthInterceptor() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            log.Println("Missing metadata in context")
            return nil, fmt.Errorf("missing metadata")
        }

        authHeader := md["authorization"]
        if len(authHeader) == 0 {
            log.Println("Missing Authorization token")
            return nil, fmt.Errorf("missing authorization token")
        }

        token := authHeader[0]
        if !strings.HasPrefix(token, "Bearer ") {
            log.Println("Invalid token format:", token)
            return nil, fmt.Errorf("invalid token format")
        }

        token = token[7:] 
        log.Println("Extracted Token:", token)

        adminID, err := utils.ParseJWT(token)
        if err != nil {
            log.Println("Invalid JWT token:", err)
            return nil, fmt.Errorf("invalid token")
        }

        log.Println("Parsed admin ID:", adminID)

 
        ctx = context.WithValue(ctx, adminIDKey, adminID)

        return handler(ctx, req)
    }
}


func GetAdminID(ctx context.Context) (int, error) {
    adminID, ok := ctx.Value(adminIDKey).(uint)
    if !ok {
        return 0, fmt.Errorf("missing admin ID in context")
    }
    return int(adminID), nil
}
