package usecase

import (
	"github.com/Prototype-1/admin_routes_service/internal/models"
	"github.com/Prototype-1/admin_routes_service/internal/repository"
	"errors"
)

type RouteUsecase interface {
	AddRoute(route *models.Route) error
	UpdateRoute(route *models.Route) error
	DeleteRoute(routeID int) error
	GetAllRoutes() ([]models.Route, error)
	GetRouteByID(routeID int) (*models.Route, error)
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
	existingRoute, err := u.repo.GetRouteByID(route.RouteID)
	if err != nil {
		return errors.New("route not found")
	}
	existingRoute.RouteName = route.RouteName
	existingRoute.StartStopID = route.StartStopID
	existingRoute.EndStopID = route.EndStopID
	existingRoute.CategoryID = route.CategoryID

	return u.repo.UpdateRoute(existingRoute)
}


func (u *routeUsecaseImpl) DeleteRoute(routeID int) error {
	_, err := u.repo.GetRouteByID(routeID)
	if err != nil {
		return errors.New("route not found")
	}
	return u.repo.DeleteRoute(routeID)
}


func (u *routeUsecaseImpl) GetAllRoutes() ([]models.Route, error) {
	return u.repo.GetAllRoutes()
}

func (u *routeUsecaseImpl) GetRouteByID(routeID int) (*models.Route, error) {
	return u.repo.GetRouteByID(routeID)
}
