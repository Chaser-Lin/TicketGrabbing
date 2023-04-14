package dao

import (
	"Project/MyProject/db"
	"Project/MyProject/models"
	"gorm.io/gorm"
)

type RouteDaoImplement interface {
	GetRouteByID(routeID int) (*models.Route, error)
	GetRoute(start, end string) (*models.Route, error)
	ListRoutes() ([]models.Route, error)
	AddRoute(route *models.Route) error
	UpdateRouteVisibility(routeID int) error
	//DeleteRoute(routeID int) error
}

type RouteDao struct {
	DB *gorm.DB
}

func NewRouteDao() RouteDaoImplement {
	return &RouteDao{
		DB: db.DB,
	}
}

func (u *RouteDao) GetRoute(start, end string) (*models.Route, error) {
	route := &models.Route{}
	err := u.DB.Where("start = ? AND end = ?", start, end).First(route).Error
	return route, err
}

func (u *RouteDao) GetRouteByID(routeID int) (*models.Route, error) {
	route := &models.Route{}
	err := u.DB.Where("route_id = ?", routeID).First(route).Error
	return route, err
}

func (u *RouteDao) ListRoutes() ([]models.Route, error) {
	var routes []models.Route
	err := u.DB.Where("visibility = ?", true).Order("CONVERT(start USING gbk), CONVERT(end USING gbk)").Find(&routes).Error
	return routes, err
}

func (u *RouteDao) AddRoute(route *models.Route) error {
	return u.DB.Create(route).Error
}

func (u *RouteDao) UpdateRouteVisibility(routeID int) error {
	return u.DB.Model(models.Route{}).Where("route_id = ?", routeID).UpdateColumn("visibility", false).Error
}

//func (u *RouteDao) DeleteRoute(routeID int) error {
//	return u.DB.Delete(models.Route{}, routeID).Error
//}
