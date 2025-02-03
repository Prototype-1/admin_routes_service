package handler

import (
	"context"
	"log"

	"github.com/Prototype-1/admin_routes_service/internal/usecase"
	pb "github.com/Prototype-1/admin_routes_service/proto"
	"github.com/Prototype-1/admin_routes_service/internal/models"
	"github.com/Prototype-1/admin_routes_service/utils"
)

type RouteHandler struct {
	pb.UnimplementedRouteServiceServer 
	RouteUsecase *usecase.RouteUsecase
}

func NewRouteHandler(routeUsecase *usecase.RouteUsecase) *RouteHandler {
	return &RouteHandler{RouteUsecase: routeUsecase}
}

func (h *RouteHandler) AddRoute(ctx context.Context, req *pb.AddRouteRequest) (*pb.AddRouteResponse, error) {
	route := models.Route{
		RouteName:   req.GetRouteName(),
		StartStopID: int(req.GetStartStopId()),
		EndStopID:   int(req.GetEndStopId()),
		CategoryID:  int(req.GetCategoryId()),
	}

	err := h.RouteUsecase.AddRoute(route)
	if err != nil {
		log.Println("Failed to add route:", err)
		return nil, err
	}

	return &pb.AddRouteResponse{Message: "Route added successfully"}, nil
}

func (h *RouteHandler) UpdateRoute(ctx context.Context, req *pb.UpdateRouteRequest) (*pb.UpdateRouteResponse, error) {
	route := models.Route{
		RouteID:     int(req.GetRouteId()),
		RouteName:   req.GetRouteName(),
		StartStopID: int(req.GetStartStopId()),
		EndStopID:   int(req.GetEndStopId()),
		CategoryID:  int(req.GetCategoryId()),
	}

	err := h.RouteUsecase.UpdateRoute(route.RouteID, route)
	if err != nil {
		log.Println("Failed to update route:", err)
		return nil, err
	}

	return &pb.UpdateRouteResponse{Message: "Route updated successfully"}, nil
}

func (h *RouteHandler) DeleteRoute(ctx context.Context, req *pb.DeleteRouteRequest) (*pb.DeleteRouteResponse, error) {
	err := h.RouteUsecase.DeleteRoute(int(req.GetRouteId()))
	if err != nil {
		log.Println("Failed to delete route:", err)
		return nil, err
	}

	return &pb.DeleteRouteResponse{Message: "Route deleted successfully"}, nil
}

func (h *RouteHandler) GetAllRoutes(ctx context.Context, req *pb.GetAllRoutesRequest) (*pb.GetAllRoutesResponse, error) {
	routes, err := h.RouteUsecase.GetAllRoutes()
	if err != nil {
		log.Println("Failed to get all routes:", err)
		return nil, err
	}

	var grpcRoutes []*pb.Route
	for _, route := range routes {
		grpcRoutes = append(grpcRoutes, &pb.Route{
			RouteId:    int32(route.RouteID),
			RouteName:  route.RouteName,
			StartStopId: int32(route.StartStopID),
			EndStopId:   int32(route.EndStopID),
			CategoryId:  int32(route.CategoryID),
			CreatedAt:   utils.FormatTime(route.CreatedAt),
			UpdatedAt:   utils.FormatTime(route.UpdatedAt),
		})
	}

	return &pb.GetAllRoutesResponse{Routes: grpcRoutes}, nil
}
