package server

import (
	"fmt"
	"log"
	"net"
	"os"

	authpb "github.com/Prototype-1/admin-auth-service/proto"
	"github.com/Prototype-1/admin_routes_service/internal/handler"
	"github.com/Prototype-1/admin_routes_service/internal/middleware"
	"github.com/Prototype-1/admin_routes_service/internal/repository"
	"github.com/Prototype-1/admin_routes_service/internal/usecase"
	pb "github.com/Prototype-1/admin_routes_service/proto"
	db "github.com/Prototype-1/admin_routes_service/utils"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"context"
)

func StartServer() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
    if jwtSecretKey == "" {
        log.Fatal("JWT_SECRET_KEY environment variable is not set")
    }
    fmt.Println("JWT_SECRET_KEY successfully loaded")

    database := db.InitDB()

conn, err := grpc.DialContext(
    context.Background(),
    "localhost:50051",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
)
if err != nil {
    log.Fatal("Failed to create gRPC client:", err)
}
defer conn.Close()

authClient := authpb.NewAdminServiceClient(conn)

    routeRepo := repository.NewRouteRepository(database)
    routeUsecase := usecase.NewRouteUsecase(routeRepo, authClient) 
    routeHandler := handler.NewRouteHandler(routeUsecase)

    grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(middleware.JWTAuthInterceptor()), 
    )

    pb.RegisterRouteServiceServer(grpcServer, routeHandler)
    reflection.Register(grpcServer)

    listener, err := net.Listen("tcp", ":50052")
    if err != nil {
        log.Fatal("Failed to listen:", err)
    }

    fmt.Println("Starting gRPC server on port 50052...")
    if err := grpcServer.Serve(listener); err != nil {
        log.Fatal("Failed to serve gRPC server:", err)
    }
}
