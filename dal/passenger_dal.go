package dal

import (
	"Project/MyProject/dal/models"
	"Project/MyProject/db"
	"gorm.io/gorm"
)

type PassengerDalImplement interface {
	AddPassenger(passenger *models.Passenger) error
	GetPassengerByID(passengerID int) (*models.Passenger, error)
	GetPassenger(name string, idNumber string) (*models.Passenger, error)
	AddUserPassenger(userPassenger *models.UserPassenger) error
	DeleteUserPassenger(userPassengerID int) error
	GetUserPassenger(userPassengerID int) (*models.UserPassenger, error)
	ListUserPassengers(userID int) ([]models.UserPassenger, error)
}

type PassengerDal struct {
	DB *gorm.DB
}

func NewPassengerDal() PassengerDalImplement {
	return &PassengerDal{
		DB: db.DB,
	}
}

func (p *PassengerDal) AddPassenger(passenger *models.Passenger) error {
	err := p.DB.Create(passenger).Error
	return err
}

func (p *PassengerDal) GetPassenger(name string, idNumber string) (*models.Passenger, error) {
	passenger := &models.Passenger{}
	err := p.DB.First(passenger, "name = ? AND id_number = ?", name, idNumber).Error
	return passenger, err
}

func (p *PassengerDal) GetPassengerByID(passengerID int) (*models.Passenger, error) {
	passenger := &models.Passenger{}
	err := p.DB.First(passenger, "passenger_id = ?", passengerID).Error
	return passenger, err
}

func (p *PassengerDal) AddUserPassenger(userPassenger *models.UserPassenger) error {
	err := p.DB.Create(userPassenger).Error
	return err
}

func (p *PassengerDal) DeleteUserPassenger(userPassengerID int) error {
	err := p.DB.Delete(&models.UserPassenger{}, userPassengerID).Error
	return err
}

func (p *PassengerDal) GetUserPassenger(userPassengerID int) (*models.UserPassenger, error) {
	userPassenger := &models.UserPassenger{}
	err := p.DB.Preload("Passenger").First(userPassenger, "id = ?", userPassengerID).Error
	return userPassenger, err
}

func (p *PassengerDal) ListUserPassengers(userID int) ([]models.UserPassenger, error) {
	var passengers []models.UserPassenger
	err := p.DB.Preload("Passenger").Find(&passengers, "user_id = ?", userID).Error
	return passengers, err
}
