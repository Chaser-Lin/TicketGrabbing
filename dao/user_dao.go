package dao

import (
	"Project/MyProject/db"
	"Project/MyProject/models"
	"gorm.io/gorm"
)

type UserDaoImplement interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(userID int) (*models.User, error)
	AddUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(userID int) error
}

type UserDao struct {
	DB *gorm.DB
}

func NewUserDao() UserDaoImplement {
	return &UserDao{
		DB: db.DB,
	}
}

func (u *UserDao) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := u.DB.Where("email = ?", email).First(user).Error
	return user, err
}

func (u *UserDao) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := u.DB.Where("username = ?", username).First(user).Error
	return user, err
}

func (u *UserDao) GetUserByID(userID int) (*models.User, error) {
	user := &models.User{}
	err := u.DB.Where("user_id = ?", userID).First(user).Error
	return user, err
}

func (u *UserDao) AddUser(user *models.User) error {
	return u.DB.Create(user).Error
}

func (u *UserDao) UpdateUser(user *models.User) error {
	return u.DB.Model(user).Updates(user).Error
}

func (u *UserDao) DeleteUser(userID int) error {
	return u.DB.Delete(models.User{}, userID).Error
}
