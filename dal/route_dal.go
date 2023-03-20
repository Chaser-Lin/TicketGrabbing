package dal

import (
	"Project/MyProject/dal/models"
	"Project/MyProject/db"
	"gorm.io/gorm"
)

type RouteDalImplement interface {
	GetRouteByID(routeID int) (*models.Route, error)
	GetRoute(start, end string) (*models.Route, error)
	ListRoutes() ([]models.Route, error)
	AddRoute(route *models.Route) error
	DeleteRoute(routeID int) error
}

type RouteDal struct {
	RouteDalImplement
	DB *gorm.DB
}

func NewRouteDal() RouteDalImplement {
	return &RouteDal{
		DB: db.DB,
	}
}

func (u *RouteDal) GetRoute(start, end string) (*models.Route, error) {
	route := &models.Route{}
	err := u.DB.Where("start = ? AND end = ?", start, end).First(route).Error
	return route, err
}

func (u *RouteDal) GetRouteByID(routeID int) (*models.Route, error) {
	route := &models.Route{}
	err := u.DB.Where("route_id = ?", routeID).First(route).Error
	return route, err
}

func (u *RouteDal) ListRoutes() ([]models.Route, error) {
	var routes []models.Route
	err := u.DB.Find(&routes).Error
	return routes, err
}

func (u *RouteDal) AddRoute(route *models.Route) error {
	err := u.DB.Create(route).Error
	return err
}

func (u *RouteDal) DeleteRoute(routeID int) error {
	err := u.DB.Delete(models.Route{}, routeID).Error
	return err
}
