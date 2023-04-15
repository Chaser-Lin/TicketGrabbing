package services

import (
	"Project/MyProject/dao"
	"Project/MyProject/models"
	"Project/MyProject/response"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"log"
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

// 获取乘车人列表服务参数
type ListUserPassengersService struct {
	UserID int `json:"user_id" form:"user_id"`
}

// 用于展示给用户的乘车人信息
type PassengerInfo struct {
	UserPassengerID int    `json:"user_passenger_id"`
	UserID          int    `json:"user_id"`
	PassengerID     int    `json:"passenger_id"`
	Name            string `json:"name"`
	IdNumber        string `json:"id_number"`
	Phone           string `json:"phone"`
}

// 乘车人相关服务接口：添加乘车人、展示所有乘车人、获取指定乘车人
type PassengerServiceImplement interface {
	AddPassenger(*AddPassengerService) error
	DeletePassenger(*DeletePassengerService) error
	GetPassenger(userPassengerID int) (*PassengerInfo, error)
	CheckPassengerBelongToUser(passengerID, userID int) (bool, error)
	ListPassengers(*ListUserPassengersService) ([]*PassengerInfo, error)
}

// 实现乘车人服务接口的实例
type PassengerService struct {
	PassengerDao dao.PassengerDaoImplement
}

func NewPassengerServices(passengerDal dao.PassengerDaoImplement) PassengerServiceImplement {
	return &PassengerService{passengerDal}
}

func (p *PassengerService) AddPassenger(service *AddPassengerService) error {
	// 判断乘客信息是否在系统中
	passenger, err := p.PassengerDao.GetPassenger(service.Name, service.IDNumber)
	if err != nil {
		// 不存在则添加乘客信息
		if err == gorm.ErrRecordNotFound {
			passenger = &models.Passenger{
				Name:     service.Name,
				IDNumber: service.IDNumber,
				Phone:    service.Phone,
			}
			err = p.PassengerDao.AddPassenger(passenger)
			if err != nil {
				return response.ErrPassengerName
			}
		} else {
			return response.ErrDbOperation
		}
	}
	// 添加用户拥有的乘客信息
	userPassenger := &models.UserPassenger{
		UserID:      service.UserID,
		PassengerID: passenger.PassengerID,
	}

	if err := p.PassengerDao.AddUserPassenger(userPassenger); err != nil {
		log.Println("AddPassenger PassengerDao.AddUserPassenger err: ", err)
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
	err := p.PassengerDao.DeleteUserPassenger(service.UserPassengerID)
	if err != nil {
		return response.ErrDbOperation
	}
	return nil
}

func (p *PassengerService) GetPassenger(userPassengerID int) (*PassengerInfo, error) {
	userPassenger, err := p.PassengerDao.GetUserPassenger(userPassengerID)
	if err == gorm.ErrRecordNotFound {
		return nil, response.ErrPassengerNotExist
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	passenger, err := p.PassengerDao.GetPassengerByID(userPassenger.PassengerID)
	return parsePassengerToInfo(passenger, userPassenger), nil
}

func (p *PassengerService) CheckPassengerBelongToUser(passengerID, userID int) (bool, error) {
	if err := p.PassengerDao.CheckPassengerBelongToUser(passengerID, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, response.ErrPassengerNotExist
		} else if err != nil {
			return false, response.ErrDbOperation
		}
	}
	return true, nil
}

func (p *PassengerService) ListPassengers(service *ListUserPassengersService) ([]*PassengerInfo, error) {
	userPassengers, err := p.PassengerDao.ListUserPassengers(service.UserID)
	if err == gorm.ErrRecordNotFound {
		return nil, response.ErrPassengerNotExist
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	passengerInfos := make([]*PassengerInfo, len(userPassengers))
	for i, userPassenger := range userPassengers {
		passenger, err := p.PassengerDao.GetPassengerByID(userPassenger.PassengerID)
		if err != nil {
			return nil, response.ErrDbOperation
		}
		passengerInfos[i] = parsePassengerToInfo(passenger, &userPassenger)
	}
	return passengerInfos, nil
}

func parsePassengerToInfo(passenger *models.Passenger, userPassenger *models.UserPassenger) *PassengerInfo {
	return &PassengerInfo{
		UserPassengerID: userPassenger.ID,
		UserID:          userPassenger.UserID,
		PassengerID:     passenger.PassengerID,
		Name:            passenger.Name,
		IdNumber:        passenger.IDNumber,
		Phone:           passenger.Phone,
	}
}
