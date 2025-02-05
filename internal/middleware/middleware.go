// package middleware

// import (
//     "context"
//     "errors"
//     "fmt"
//     "strings"
//     "github.com/golang-jwt/jwt/v4"
//     "google.golang.org/grpc"
//     "google.golang.org/grpc/metadata"
//     "github.com/Prototype-1/admin_routes_service/utils"
//     "log"
//     "time"
//     "strconv"
// )

// type contextKey string

// const adminIDKey contextKey = "admin_id"

// func VerifyJWT(tokenString string, secretKey string) (jwt.MapClaims, error) {
//     token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//         if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//             return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
//         }
//         return []byte(secretKey), nil
//     })

//     if err != nil {
//         return nil, fmt.Errorf("invalid JWT token: %w", err)
//     }
//     if !token.Valid {
//         return nil, fmt.Errorf("invalid JWT token: signature is invalid")
//     }

//     claims, ok := token.Claims.(jwt.MapClaims)
//     if !ok {
//         return nil, errors.New("invalid token claims")
//     }
//     log.Println("Extracted Claims:", claims)

//     if exp, ok := claims["exp"].(float64); ok {
//         if time.Now().Unix() > int64(exp) {
//             return nil, errors.New("token has expired")
//         }
//     } else {
//         return nil, errors.New("invalid expiration claim")
//     }

//     return claims, nil
// }

// func JWTAuthInterceptor() grpc.UnaryServerInterceptor {
//     return func(
//         ctx context.Context,
//         req interface{},
//         info *grpc.UnaryServerInfo,
//         handler grpc.UnaryHandler,
//     ) (interface{}, error) {
//         log.Println("JWTAuthInterceptor triggered for method:", info.FullMethod)

//         md, ok := metadata.FromIncomingContext(ctx)
//         if !ok {
//             log.Println("Missing metadata in context")
//             return nil, fmt.Errorf("missing metadata")
//         }

//         authHeader := md["authorization"]
//         if len(authHeader) == 0 {
//             log.Println("Missing Authorization token")
//             return nil, fmt.Errorf("missing authorization token")
//         }

//         token := authHeader[0]
//         if !strings.HasPrefix(token, "Bearer ") {
//             log.Println("Invalid token format:", token)
//             return nil, fmt.Errorf("invalid token format")
//         }

//         token = token[7:] 
//         log.Println("Extracted Token:", token)

//         adminID, err := utils.ParseJWT(token)
//         if err != nil {
//             log.Println("Invalid JWT token:", err)
//             return nil, fmt.Errorf("invalid token")
//         }

//         log.Println("Parsed admin ID:", adminID)

//         md = metadata.Pairs("admin_id", fmt.Sprintf("%d", adminID))
//         ctx = metadata.NewIncomingContext(ctx, md)        

//         log.Println("Context after adding admin ID:", ctx)

//         return handler(ctx, req)
//     }
// }

// func GetAdminID(ctx context.Context) (int, error) {
//     md, ok := metadata.FromIncomingContext(ctx)
//     if !ok {
//         log.Println("Metadata not found in context")
//         return 0, fmt.Errorf("metadata not found")
//     }

//     adminIDStrs := md.Get("admin_id")
//     if len(adminIDStrs) == 0 {
//         log.Println("Admin ID not found in metadata")
//         return 0, fmt.Errorf("admin ID not found in metadata")
//     }

//     adminID, err := strconv.Atoi(adminIDStrs[0])
//     if err != nil {
//         log.Println("Invalid admin ID format:", err)
//         return 0, fmt.Errorf("invalid admin ID format")
//     }

//     log.Println("Extracted Admin ID from metadata:", adminID)
//     return adminID, nil
// }


package middleware

import (
    "context"
    "errors"
    "fmt"
    "strings"
    "time"
    "log"

    "github.com/golang-jwt/jwt/v4"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"

    "github.com/Prototype-1/admin_routes_service/utils"
)

type contextKey string

const AdminIDKey contextKey = "admin_id"

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
        return nil, errors.New("invalid JWT token: signature is invalid")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, errors.New("invalid token claims")
    }
    log.Println("Extracted Claims:", claims)

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
        log.Println("JWTAuthInterceptor triggered for method:", info.FullMethod)

        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            log.Println("Metadata missing from context")
            return nil, fmt.Errorf("metadata missing")
        }

        authHeader, found := md["authorization"]
        if !found || len(authHeader) == 0 {
            log.Println("Authorization token not found in metadata")
            return nil, fmt.Errorf("authorization token not found")
        }

        tokenString := authHeader[0]
        if !strings.HasPrefix(tokenString, "Bearer ") {
            log.Println("Invalid token format received:", tokenString)
            return nil, fmt.Errorf("invalid token format")
        }

        tokenString = strings.TrimPrefix(tokenString, "Bearer ")
        log.Println("Extracted Token:", tokenString)

        secretKey := utils.GetJWTSecret() 
        claims, err := VerifyJWT(tokenString, string(secretKey))
        if err != nil {
            log.Println("Invalid JWT token:", err)
            return nil, fmt.Errorf("invalid JWT token")
        }

        adminIDFloat, ok := claims["admin_id"].(float64)
        if !ok {
            log.Println("Admin ID missing or invalid in token claims")
            return nil, fmt.Errorf("invalid admin ID in token claims")
        }

        adminID := int(adminIDFloat)
        log.Println("Extracted admin_id:", adminID)

        // Add admin_id to context
        ctx = context.WithValue(ctx, AdminIDKey, adminID)
        log.Println("Context after adding admin ID:", ctx)

        return handler(ctx, req)
    }
}

func GetAdminID(ctx context.Context) (int, error) {
    adminID, ok := ctx.Value(AdminIDKey).(int)
    if !ok {
        log.Println("Admin ID not found in context")
        return 0, fmt.Errorf("admin ID not found")
    }

    log.Println("Retrieved Admin ID from context:", adminID)
    return adminID, nil
}
