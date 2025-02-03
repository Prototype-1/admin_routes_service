package usecase

import (
	"github.com/Prototype-1/admin_routes_service/internal/models"
	"github.com/Prototype-1/admin_routes_service/internal/repository"
)

type RouteUsecase struct {
    RouteRepo repository.RouteRepository
}

func NewRouteUsecase(repo repository.RouteRepository) *RouteUsecase {
    return &RouteUsecase{RouteRepo: repo}
}

func (u *RouteUsecase) AddRoute(route models.Route) error {
    return u.RouteRepo.AddRoute(route)
}

func (u *RouteUsecase) UpdateRoute(id int, route models.Route) error {
    route.RouteID = id
    return u.RouteRepo.UpdateRoute(route)
}

func (u *RouteUsecase) DeleteRoute(id int) error {
    return u.RouteRepo.DeleteRoute(id)
}

func (u *RouteUsecase) GetAllRoutes() ([]models.Route, error) {
    return u.RouteRepo.GetAllRoutes()
}
