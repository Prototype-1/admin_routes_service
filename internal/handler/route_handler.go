package handlers

import (
	"context"
	"log"

	"github.com/Prototype-1/admin_routes_service/internal/models"
	"github.com/Prototype-1/admin_routes_service/internal/usecase"
	pb "github.com/Prototype-1/admin_routes_service/proto"
	"github.com/Prototype-1/admin_routes_service/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/metadata"
)

type RouteServer struct {
	pb.UnimplementedRouteServiceServer
	usecase usecase.RouteUsecase
}

func NewRouteServer(usecase usecase.RouteUsecase) *RouteServer {
	return &RouteServer{usecase: usecase}
}

func (s *RouteServer) AddRoute(ctx context.Context, req *pb.AddRouteRequest) (*pb.AddRouteResponse, error) {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
    }

    authHeader, exists := md["authorization"]
    if !exists || len(authHeader) == 0 {
        return nil, status.Errorf(codes.Unauthenticated, "authorization token missing")
    }

    token := authHeader[0]

    _, role, err := utils.ParseJWT(token)
    if err != nil {
        return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
    }

    if role != "admin" {
        return nil, status.Errorf(codes.PermissionDenied, "only admins can add routes")
    }

    log.Println("Authorization successful, role:", role)

    route := &models.Route{
        RouteName:   req.RouteName,
        StartStopID: int(req.StartStopId),
        EndStopID:   int(req.EndStopId),
        CategoryID:  int(req.CategoryId),
    }

    err = s.usecase.AddRoute(route)
    if err != nil {
        log.Printf("Failed to add route: %v", err)
        return nil, status.Errorf(codes.Internal, "failed to add route: %v", err)
    }

    return &pb.AddRouteResponse{Message: "Route added successfully"}, nil
}

func (s *RouteServer) UpdateRoute(ctx context.Context, req *pb.UpdateRouteRequest) (*pb.UpdateRouteResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	authHeader, exists := md["authorization"]
	if !exists || len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token missing")
	}

	token := authHeader[0]
	_, role, err := utils.ParseJWT(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	if role != "admin" {
		return nil, status.Errorf(codes.PermissionDenied, "only admins can update routes")
	}
	existingRoute, err := s.usecase.GetRouteByID(int(req.RouteId))
	if err != nil {
		log.Printf("Route not found: %v", err)
		return nil, status.Errorf(codes.NotFound, "route not found")
	}
	existingRoute.RouteName = req.RouteName
	existingRoute.StartStopID = int(req.StartStopId)
	existingRoute.EndStopID = int(req.EndStopId)
	existingRoute.CategoryID = int(req.CategoryId)

	err = s.usecase.UpdateRoute(existingRoute)
	if err != nil {
		log.Printf("Failed to update route: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update route")
	}

	return &pb.UpdateRouteResponse{Message: "Route updated successfully"}, nil
}

func (s *RouteServer) DeleteRoute(ctx context.Context, req *pb.DeleteRouteRequest) (*pb.DeleteRouteResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	authHeader, exists := md["authorization"]
	if !exists || len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token missing")
	}

	token := authHeader[0]
	_, role, err := utils.ParseJWT(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	if role != "admin" {
		return nil, status.Errorf(codes.PermissionDenied, "only admins can delete routes")
	}

	existingRoute, err := s.usecase.GetRouteByID(int(req.RouteId))
	if err != nil {
		log.Printf("Route not found: %v", err)
		return nil, status.Errorf(codes.NotFound, "route not found")
	}
	err = s.usecase.DeleteRoute(existingRoute.RouteID)
	if err != nil {
		log.Printf("Failed to delete route: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete route")
	}

	return &pb.DeleteRouteResponse{Message: "Route deleted successfully"}, nil
}


func (s *RouteServer) GetAllRoutes(ctx context.Context, req *pb.GetAllRoutesRequest) (*pb.GetAllRoutesResponse, error) {
	routes, err := s.usecase.GetAllRoutes()
	if err != nil {
		log.Printf("Failed to fetch routes: %v", err)
		return nil, err
	}

	var routeList []*pb.Route
	for _, route := range routes {
		routeList = append(routeList, &pb.Route{
			RouteId:     int32(route.RouteID),
			RouteName:   route.RouteName,
			StartStopId: int32(route.StartStopID),
			EndStopId:   int32(route.EndStopID),
			CategoryId:  int32(route.CategoryID),
			CreatedAt:   route.CreatedAt.String(),
			UpdatedAt:   route.UpdatedAt.String(),
		})
	}

	return &pb.GetAllRoutesResponse{Routes: routeList}, nil
}
