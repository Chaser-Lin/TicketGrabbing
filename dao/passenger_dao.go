package dao

import (
	"Project/MyProject/db"
	"Project/MyProject/models"
	"gorm.io/gorm"
)

type PassengerDaoImplement interface {
	AddPassenger(passenger *models.Passenger) error
	GetPassengerByID(passengerID int) (*models.Passenger, error)
	GetPassenger(name string, idNumber string) (*models.Passenger, error)
	AddUserPassenger(userPassenger *models.UserPassenger) error
	DeleteUserPassenger(userPassengerID int) error
	GetUserPassenger(userPassengerID int) (*models.UserPassenger, error)
	ListUserPassengers(userID int) ([]models.UserPassenger, error)
	CheckPassengerBelongToUser(passengerID, userID int) error
}

type PassengerDao struct {
	DB *gorm.DB
}

func NewPassengerDao() PassengerDaoImplement {
	return &PassengerDao{
		DB: db.DB,
	}
}

func (p *PassengerDao) AddPassenger(passenger *models.Passenger) error {
	return p.DB.Create(passenger).Error
}

func (p *PassengerDao) GetPassenger(name string, idNumber string) (*models.Passenger, error) {
	passenger := &models.Passenger{}
	err := p.DB.First(passenger, "name = ? AND id_number = ?", name, idNumber).Error
	return passenger, err
}

func (p *PassengerDao) GetPassengerByID(passengerID int) (*models.Passenger, error) {
	passenger := &models.Passenger{}
	err := p.DB.First(passenger, "passenger_id = ?", passengerID).Error
	return passenger, err
}

func (p *PassengerDao) AddUserPassenger(userPassenger *models.UserPassenger) error {
	err := p.DB.Create(userPassenger).Error
	return err
}

func (p *PassengerDao) DeleteUserPassenger(userPassengerID int) error {
	err := p.DB.Delete(&models.UserPassenger{}, userPassengerID).Error
	return err
}

func (p *PassengerDao) GetUserPassenger(userPassengerID int) (*models.UserPassenger, error) {
	userPassenger := &models.UserPassenger{}
	err := p.DB.First(userPassenger, "id = ?", userPassengerID).Error
	return userPassenger, err
}

func (p *PassengerDao) ListUserPassengers(userID int) ([]models.UserPassenger, error) {
	var passengers []models.UserPassenger
	err := p.DB.Find(&passengers, "user_id = ?", userID).Error
	return passengers, err
}

func (p *PassengerDao) CheckPassengerBelongToUser(passengerID, userID int) error {
	return p.DB.First(&models.UserPassenger{}, "passenger_id = ? AND user_id = ?", passengerID, userID).Error
}
