package repository

import (
	"github.com/Prototype-1/admin_routes_service/internal/models"
	"gorm.io/gorm"
)

type RouteRepository interface {
    AddRoute(route models.Route) error
    UpdateRoute(route models.Route) error
    DeleteRoute(id int) error
    GetAllRoutes() ([]models.Route, error)
}

type RouteRepositoryImpl struct {
	DB *gorm.DB
}

func NewRouteRepository(db *gorm.DB) *RouteRepositoryImpl {
    return &RouteRepositoryImpl{DB: db}
}

func (r *RouteRepositoryImpl) AddRoute(route models.Route) error {
	return r.DB.Create(&route).Error
}

func (r *RouteRepositoryImpl) UpdateRoute(route models.Route) error {
    return r.DB.Save(&route).Error
}

func (r *RouteRepositoryImpl) DeleteRoute(id int) error {
    var route models.Route
    if err := r.DB.First(&route, id).Error; err != nil {
        return err
    }
    return r.DB.Delete(&route).Error
}

func (r *RouteRepositoryImpl) GetAllRoutes() ([]models.Route, error) {
    var routes []models.Route
    err := r.DB.Find(&routes).Error
    return routes, err
}
