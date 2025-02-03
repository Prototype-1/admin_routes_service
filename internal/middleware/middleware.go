package middleware

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Prototype-1/admin_routes_service/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func JWTAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("missing metadata")
		}
		authHeader := md["authorization"]
		if len(authHeader) == 0 {
			return nil, fmt.Errorf("missing authorization token")
		}

		token := authHeader[0]
		if !strings.HasPrefix(token, "Bearer ") {
			return nil, fmt.Errorf("invalid token format")
		}

		token = token[7:]
		_, err := utils.ParseJWT(token)
		if err != nil {
			log.Println("Invalid JWT token:", err)
			return nil, fmt.Errorf("invalid token")
		}
		return handler(ctx, req)
	}
}
