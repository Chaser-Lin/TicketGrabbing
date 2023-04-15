package services

import (
	"Project/MyProject/dao"
	"Project/MyProject/models"
	"Project/MyProject/response"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// 添加路线服务参数
type AddRouteService struct {
	Start  string `json:"start" form:"start" binding:"required"`
	End    string `json:"end" form:"end" binding:"required"`
	Length uint32 `json:"length" form:"length" binding:"required,number"`
}

// 获取路线服务参数
type GetRouteService struct {
	Start string `json:"start" form:"start" binding:"required"`
	End   string `json:"end" form:"end" binding:"required"`
}

// 路线相关服务接口：添加路线、展示所有路线、获取指定路线
type RouteServiceImplement interface {
	AddRoute(*AddRouteService) error
	GetRoute(*GetRouteService) (*models.Route, error)
	DeleteRoute(routeID int) error
	ListRoutes() ([]models.Route, error)
	GetRouteByID(routeID int) (*models.Route, error)
}

// 实现路线服务接口的实例
type RouteService struct {
	RouteDao dao.RouteDaoImplement
}

func NewRouteServices(routeDal dao.RouteDaoImplement) RouteServiceImplement {
	return &RouteService{routeDal}
}

func (r *RouteService) AddRoute(service *AddRouteService) error {
	route := &models.Route{
		Start:  service.Start,
		End:    service.End,
		Length: service.Length,
	}

	if err := r.RouteDao.AddRoute(route); err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 { // 1062:Duplicate，重复数据
				return response.ErrRouteExist
			}
		}
		return response.ErrDbOperation
	}
	return nil
}

func (r *RouteService) GetRoute(service *GetRouteService) (*models.Route, error) {
	route, err := r.RouteDao.GetRoute(service.Start, service.End)
	if err == gorm.ErrRecordNotFound {
		return nil, response.ErrRouteNotExist
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	return route, nil
}

func (r *RouteService) DeleteRoute(routeID int) error {
	_, err := r.GetRouteByID(routeID)
	if err != nil {
		return err
	}
	err = r.RouteDao.UpdateRouteVisibility(routeID)
	if err != nil {
		return response.ErrDbOperation
	}
	return nil
}

func (u *RouteService) GetRouteByID(routeID int) (*models.Route, error) {
	route, err := u.RouteDao.GetRouteByID(routeID)
	if err == gorm.ErrRecordNotFound {
		return nil, response.ErrRouteNotExist
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	return route, nil
}

func (r *RouteService) ListRoutes() ([]models.Route, error) {
	routes, err := r.RouteDao.ListRoutes()
	if err == gorm.ErrRecordNotFound {
		return nil, response.EmptyRouteList
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	return routes, nil
}
