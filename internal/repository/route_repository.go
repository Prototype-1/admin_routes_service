package repository

import (
	"github.com/Prototype-1/admin_routes_service/internal/models"
	"gorm.io/gorm"
	"errors"
)

type RouteRepository interface {
	AddRoute(route *models.Route) error
	UpdateRoute(route *models.Route) error
	DeleteRoute(routeID int) error
	GetAllRoutes() ([]models.Route, error)
	GetRouteByID(routeID int) (*models.Route, error) 
}

type routeRepositoryImpl struct {
	db *gorm.DB
}

func NewRouteRepository(db *gorm.DB) RouteRepository {
	return &routeRepositoryImpl{db: db}
}

func (r *routeRepositoryImpl) AddRoute(route *models.Route) error {
	return r.db.Create(route).Error
}

func (r *routeRepositoryImpl) UpdateRoute(route *models.Route) error {
	return r.db.Model(&models.Route{}).
		Where("route_id = ?", route.RouteID).
		Updates(route).Error
}

func (r *routeRepositoryImpl) DeleteRoute(routeID int) error {
	return r.db.Delete(&models.Route{}, routeID).Error
}

func (r *routeRepositoryImpl) GetAllRoutes() ([]models.Route, error) {
	var routes []models.Route
	err := r.db.Find(&routes).Error
	return routes, err
}

func (r *routeRepositoryImpl) GetRouteByID(routeID int) (*models.Route, error) {
	var route models.Route
	err := r.db.Where("route_id = ?", routeID).First(&route).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("route not found")
		}
		return nil, err
	}
	return &route, nil
}

