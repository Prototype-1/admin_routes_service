package usecase

import (
	"context"
	"fmt"

	authpb "github.com/Prototype-1/admin-auth-service/proto" 
	"github.com/Prototype-1/admin_routes_service/internal/models"
	"github.com/Prototype-1/admin_routes_service/internal/repository"
)

type RouteUsecase struct {
	RouteRepo  repository.RouteRepository
	AuthClient authpb.AdminServiceClient 
}

func NewRouteUsecase(routeRepo repository.RouteRepository, authClient authpb.AdminServiceClient) *RouteUsecase {
	return &RouteUsecase{
		RouteRepo:  routeRepo,
		AuthClient: authClient, 
	}
}

func (u *RouteUsecase) VerifyAdmin(token string) error {
	ctx := context.Background()
	_, err := u.AuthClient.GetAllUsers(ctx, &authpb.Empty{})
	if err != nil {
		return fmt.Errorf("unauthorized admin: %v", err)
	}
	return nil
}

func (u *RouteUsecase) AddRoute(token string, route models.Route) error {
	if err := u.VerifyAdmin(token); err != nil {
		return err
	}
	return u.RouteRepo.AddRoute(route)
}

func (u *RouteUsecase) UpdateRoute(token string, id int, route models.Route) error {
	if err := u.VerifyAdmin(token); err != nil {
		return err
	}
	route.RouteID = id
	return u.RouteRepo.UpdateRoute(route)
}

func (u *RouteUsecase) DeleteRoute(token string, id int) error {
	if err := u.VerifyAdmin(token); err != nil {
		return err
	}
	return u.RouteRepo.DeleteRoute(id)
}

func (u *RouteUsecase) GetAllRoutes(token string) ([]models.Route, error) {
	if err := u.VerifyAdmin(token); err != nil {
		return nil, err
	}
	return u.RouteRepo.GetAllRoutes()
}
