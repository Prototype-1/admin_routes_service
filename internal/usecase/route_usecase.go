package usecase

import (
	"github.com/Prototype-1/admin_routes_service/internal/models"
	"github.com/Prototype-1/admin_routes_service/internal/repository"
)

type RouteUsecase interface {
	AddRoute(route *models.Route) error
	UpdateRoute(route *models.Route) error
	DeleteRoute(routeID int) error
	GetAllRoutes() ([]models.Route, error)
}

type routeUsecaseImpl struct {
	repo repository.RouteRepository
}

func NewRouteUsecase(repo repository.RouteRepository) RouteUsecase {
	return &routeUsecaseImpl{repo: repo}
}

func (u *routeUsecaseImpl) AddRoute(route *models.Route) error {
	return u.repo.AddRoute(route)
}

func (u *routeUsecaseImpl) UpdateRoute(route *models.Route) error {
	return u.repo.UpdateRoute(route)
}

func (u *routeUsecaseImpl) DeleteRoute(routeID int) error {
	return u.repo.DeleteRoute(routeID)
}

func (u *routeUsecaseImpl) GetAllRoutes() ([]models.Route, error) {
	return u.repo.GetAllRoutes()
}
