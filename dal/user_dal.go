package dal

import (
	"Project/MyProject/dal/models"
	"Project/MyProject/db"
	"gorm.io/gorm"
)

type UserDalImplement interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(userID int) (*models.User, error)
	AddUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(userID int) error
}

type UserDal struct {
	DB *gorm.DB
}

func NewUserDal() UserDalImplement {
	return &UserDal{
		DB: db.DB,
	}
}

func (u *UserDal) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := u.DB.Where("email = ?", email).First(user).Error
	return user, err
}

func (u *UserDal) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := u.DB.Where("username = ?", username).First(user).Error
	return user, err
}

func (u *UserDal) GetUserByID(userID int) (*models.User, error) {
	user := &models.User{}
	err := u.DB.Where("user_id = ?", userID).First(user).Error
	return user, err
}

func (u *UserDal) AddUser(user *models.User) error {
	err := u.DB.Create(user).Error
	return err
}

func (u *UserDal) UpdateUser(user *models.User) error {
	err := u.DB.Model(&models.User{}).Updates(user).Error
	return err
}

func (u *UserDal) DeleteUser(userID int) error {
	err := u.DB.Delete(models.User{}, userID).Error
	return err
}
