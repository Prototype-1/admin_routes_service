package handler

import (
	"context"
	"log"
	"strconv"
	"github.com/Prototype-1/admin_routes_service/internal/usecase"
	pb "github.com/Prototype-1/admin_routes_service/proto"
	"github.com/Prototype-1/admin_routes_service/internal/models"
	"github.com/Prototype-1/admin_routes_service/utils"
	mw "github.com/Prototype-1/admin_routes_service/internal/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RouteHandler struct {
	pb.UnimplementedRouteServiceServer
	RouteUsecase *usecase.RouteUsecase
}

func NewRouteHandler(routeUsecase *usecase.RouteUsecase) *RouteHandler {
	return &RouteHandler{RouteUsecase: routeUsecase}
}

func (h *RouteHandler) getAdminID(ctx context.Context) (int, error) {
    adminID, err := mw.GetAdminID(ctx)
    if err != nil {
        log.Printf("Failed to get admin ID from context: %v", err)
        return 0, status.Errorf(codes.Unauthenticated, "Authorization failed")
    }
    return adminID, nil
}

func (h *RouteHandler) AddRoute(ctx context.Context, req *pb.AddRouteRequest) (*pb.AddRouteResponse, error) {
	adminID, err := h.getAdminID(ctx)
	if err != nil {
		return nil, err
	}

	err = h.RouteUsecase.AddRoute(strconv.Itoa(adminID), models.Route{
		RouteName:   req.GetRouteName(),
		StartStopID: int(req.GetStartStopId()),
		EndStopID:   int(req.GetEndStopId()),
		CategoryID:  int(req.GetCategoryId()),
	})
	if err != nil {
		log.Printf("Failed to add route: %v", err)
		return nil, err
	}

	return &pb.AddRouteResponse{Message: "Route added successfully"}, nil
}

func (h *RouteHandler) UpdateRoute(ctx context.Context, req *pb.UpdateRouteRequest) (*pb.UpdateRouteResponse, error) {
	adminID, err := h.getAdminID(ctx)
	if err != nil {
		return nil, err
	}

	err = h.RouteUsecase.UpdateRoute(strconv.Itoa(adminID), int(req.GetRouteId()), models.Route{
		RouteName:   req.GetRouteName(),
		StartStopID: int(req.GetStartStopId()),
		EndStopID:   int(req.GetEndStopId()),
		CategoryID:  int(req.GetCategoryId()),
	})
	if err != nil {
		log.Printf("Failed to update route: %v", err)
		return nil, err
	}

	return &pb.UpdateRouteResponse{Message: "Route updated successfully"}, nil
}

func (h *RouteHandler) DeleteRoute(ctx context.Context, req *pb.DeleteRouteRequest) (*pb.DeleteRouteResponse, error) {
	adminID, err := h.getAdminID(ctx)
	if err != nil {
		return nil, err
	}

	err = h.RouteUsecase.DeleteRoute(strconv.Itoa(adminID), int(req.GetRouteId()))
	if err != nil {
		log.Printf("Failed to delete route: %v", err)
		return nil, err
	}

	return &pb.DeleteRouteResponse{Message: "Route deleted successfully"}, nil
}

func (h *RouteHandler) GetAllRoutes(ctx context.Context, req *pb.GetAllRoutesRequest) (*pb.GetAllRoutesResponse, error) {
	adminID, err := h.getAdminID(ctx)
	if err != nil {
		return nil, err
	}

	routes, err := h.RouteUsecase.GetAllRoutes(strconv.Itoa(adminID))
	if err != nil {
		log.Printf("Failed to get all routes: %v", err)
		return nil, err
	}

	var grpcRoutes []*pb.Route
	for _, route := range routes {
		grpcRoutes = append(grpcRoutes, &pb.Route{
			RouteId:     int32(route.RouteID),
			RouteName:   route.RouteName,
			StartStopId: int32(route.StartStopID),
			EndStopId:   int32(route.EndStopID),
			CategoryId:  int32(route.CategoryID),
			CreatedAt:   utils.FormatTime(route.CreatedAt),
			UpdatedAt:   utils.FormatTime(route.UpdatedAt),
		})
	}

	return &pb.GetAllRoutesResponse{Routes: grpcRoutes}, nil
}