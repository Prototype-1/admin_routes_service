package handlers

import (
	"context"
	"log"

	"github.com/Prototype-1/admin_routes_service/internal/models"
	"github.com/Prototype-1/admin_routes_service/internal/usecase"
	pb "github.com/Prototype-1/admin_routes_service/proto"
)

type RouteServer struct {
	pb.UnimplementedRouteServiceServer
	usecase usecase.RouteUsecase
}

func NewRouteServer(usecase usecase.RouteUsecase) *RouteServer {
	return &RouteServer{usecase: usecase}
}

func (s *RouteServer) AddRoute(ctx context.Context, req *pb.AddRouteRequest) (*pb.AddRouteResponse, error) {
	route := &models.Route{
		RouteName:   req.RouteName,
		StartStopID: int(req.StartStopId),
		EndStopID:   int(req.EndStopId),
		CategoryID:  int(req.CategoryId),
	}

	err := s.usecase.AddRoute(route)
	if err != nil {
		log.Printf("Failed to add route: %v", err)
		return nil, err
	}

	return &pb.AddRouteResponse{Message: "Route added successfully"}, nil
}

func (s *RouteServer) UpdateRoute(ctx context.Context, req *pb.UpdateRouteRequest) (*pb.UpdateRouteResponse, error) {
	route := &models.Route{
		RouteID:     int(req.RouteId),
		RouteName:   req.RouteName,
		StartStopID: int(req.StartStopId),
		EndStopID:   int(req.EndStopId),
		CategoryID:  int(req.CategoryId),
	}

	err := s.usecase.UpdateRoute(route)
	if err != nil {
		log.Printf("Failed to update route: %v", err)
		return nil, err
	}

	return &pb.UpdateRouteResponse{Message: "Route updated successfully"}, nil
}

func (s *RouteServer) DeleteRoute(ctx context.Context, req *pb.DeleteRouteRequest) (*pb.DeleteRouteResponse, error) {
	err := s.usecase.DeleteRoute(int(req.RouteId))
	if err != nil {
		log.Printf("Failed to delete route: %v", err)
		return nil, err
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
