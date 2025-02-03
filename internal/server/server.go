package server

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	db "github.com/Prototype-1/admin_routes_service/utils"
	"github.com/Prototype-1/admin_routes_service/internal/handler"
	"github.com/Prototype-1/admin_routes_service/internal/repository"
	"github.com/Prototype-1/admin_routes_service/internal/usecase"
	pb "github.com/Prototype-1/admin_routes_service/proto"
	"github.com/Prototype-1/admin_routes_service/internal/middleware"
	"github.com/joho/godotenv"
)

func StartServer() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database := db.InitDB()

	routeRepo := repository.NewRouteRepository(database)
	routeUsecase := usecase.NewRouteUsecase(routeRepo)
	routeHandler := handler.NewRouteHandler(routeUsecase)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.JWTAuthInterceptor()), 
	)

	pb.RegisterRouteServiceServer(grpcServer, routeHandler)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	fmt.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Failed to serve gRPC server:", err)
	}
}
