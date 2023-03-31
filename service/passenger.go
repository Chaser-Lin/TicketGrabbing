package service

import (
	"Project/MyProject/dal"
	"Project/MyProject/dal/models"
	"Project/MyProject/response"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// 添加乘车人服务参数
type AddPassengerService struct {
	// 传入乘客相关信息
	UserID   int    `json:"user_id" form:"user_id"`
	Name     string `json:"name" form:"name" binding:"required"`
	IDNumber string `json:"id_number" form:"id_number" binding:"required"`
	Phone    string `json:"phone" form:"phone" binding:"required,number"`
}

// 删除乘车人服务参数
type DeletePassengerService struct {
	UserPassengerID int `json:"user_passenger_id" form:"user_passenger_id" binding:"required"`
}

// 获取乘车人信息服务参数
type GetPassengerService struct {
	UserPassengerID int `json:"user_passenger_id" form:"user_passenger_id" binding:"required"`
}

// 获取乘车人列表服务参数
type ListUserPassengersService struct {
	UserID int `json:"user_id" form:"user_id"`
}

// 用于展示给用户的乘车人信息
type PassengerInfo struct {
	UserPassengerID int    `json:"user_passenger_id"`
	Name            string `json:"name"`
	IdNumber        string `json:"id_number"`
	Phone           string `json:"phone"`
}

// 乘车人相关服务接口：添加乘车人、展示所有乘车人、获取指定乘车人
type PassengerServiceImplement interface {
	AddPassenger(*AddPassengerService) error
	DeletePassenger(*DeletePassengerService) error
	GetPassenger(*GetPassengerService) (*PassengerInfo, error)
	ListPassengers(*ListUserPassengersService) ([]*PassengerInfo, error)
}

// 实现乘车人服务接口的实例
type PassengerService struct {
	PassengerDal dal.PassengerDalImplement
}

func NewPassengerServices(passengerDal dal.PassengerDalImplement) PassengerServiceImplement {
	return &PassengerService{passengerDal}
}

func (p *PassengerService) AddPassenger(service *AddPassengerService) error {
	// 判断乘客信息是否在系统中
	passenger, err := p.PassengerDal.GetPassenger(service.Name, service.IDNumber)
	if err != nil {
		// 不存在则添加乘客信息
		if err == gorm.ErrRecordNotFound {
			passenger = &models.Passenger{
				Name:     service.Name,
				IDNumber: service.IDNumber,
				Phone:    service.Phone,
			}
			err = p.PassengerDal.AddPassenger(passenger)
			if err != nil {
				return response.ErrDbOperation
			}
		}
		return response.ErrDbOperation
	}
	// 添加用户拥有的乘客信息
	userPassenger := &models.UserPassenger{
		UserID:      service.UserID,
		PassengerID: passenger.PassengerID,
	}

	if err := p.PassengerDal.AddUserPassenger(userPassenger); err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 { // 1062:Duplicate，重复数据
				return response.ErrSamePassengerExist
			}
		}
		return response.ErrDbOperation
	}

	return nil
}

func (p *PassengerService) DeletePassenger(service *DeletePassengerService) error {
	err := p.PassengerDal.DeleteUserPassenger(service.UserPassengerID)
	if err != nil {
		return response.ErrDbOperation
	}
	return nil
}

func (p *PassengerService) GetPassenger(service *GetPassengerService) (*PassengerInfo, error) {
	userPassenger, err := p.PassengerDal.GetUserPassenger(service.UserPassengerID)
	if err == gorm.ErrRecordNotFound {
		return nil, response.ErrPassengerNotExist
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	passenger, err := p.PassengerDal.GetPassengerByID(userPassenger.PassengerID)
	return parsePassengerToInfo(passenger, userPassenger.ID), nil
}

func (p *PassengerService) ListPassengers(service *ListUserPassengersService) ([]*PassengerInfo, error) {
	userPassengers, err := p.PassengerDal.ListUserPassengers(service.UserID)
	if err == gorm.ErrRecordNotFound {
		return nil, response.ErrPassengerNotExist
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	passengerInfos := make([]*PassengerInfo, len(userPassengers))
	for i, userPassenger := range userPassengers {
		passenger, err := p.PassengerDal.GetPassengerByID(userPassenger.PassengerID)
		if err != nil {
			return nil, response.ErrDbOperation
		}
		passengerInfos[i] = parsePassengerToInfo(passenger, userPassenger.ID)
	}
	return passengerInfos, nil
}

func parsePassengerToInfo(passenger *models.Passenger, userPassengerID int) *PassengerInfo {
	return &PassengerInfo{
		UserPassengerID: userPassengerID,
		Name:            passenger.Name,
		IdNumber:        passenger.IDNumber,
		Phone:           passenger.Phone,
	}
}
