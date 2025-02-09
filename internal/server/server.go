package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"github.com/Prototype-1/admin_routes_service/internal/handler"
	"github.com/Prototype-1/admin_routes_service/internal/repository"
	"github.com/Prototype-1/admin_routes_service/internal/usecase"
	pb "github.com/Prototype-1/admin_routes_service/proto"
	"github.com/Prototype-1/admin_routes_service/utils"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func StartServer() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
    if jwtSecretKey == "" {
        utils.Log.Fatal("JWT_SECRET_KEY environment variable is not set")
    }
    fmt.Println("JWT_SECRET_KEY successfully loaded")

    utils.InitLogger() 
	utils.Log.Info("Logger initialized successfully")

    database := utils.InitDB()


    routeRepo := repository.NewRouteRepository(database)
	routeUsecase := usecase.NewRouteUsecase(routeRepo)
	routeServer := handlers.NewRouteServer(routeUsecase)

    grpcServer := grpc.NewServer()
	pb.RegisterRouteServiceServer(grpcServer, routeServer)

    listener, err := net.Listen("tcp", "0.0.0.0:50053")
    if err != nil {
        log.Fatal("Failed to listen:", err)
    }

    fmt.Println("Starting gRPC server on port 50053...")
    if err := grpcServer.Serve(listener); err != nil {
        log.Fatal("Failed to serve gRPC server:", err)
    }
}
