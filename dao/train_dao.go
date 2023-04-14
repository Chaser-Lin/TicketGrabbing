package dao

import (
	"Project/MyProject/db"
	"Project/MyProject/models"
	"gorm.io/gorm"
)

type TrainDaoImplement interface {
	GetTrain(trainID string) (*models.Train, error)
	ListTrains() ([]models.Train, error)
	AddTrain(train *models.Train) error
	UpdateTrainVisibility(trainID string) error
	//DeleteTrain(trainID string) error
}

type TrainDao struct {
	DB *gorm.DB
}

func NewTrainDao() TrainDaoImplement {
	return &TrainDao{
		DB: db.DB,
	}
}

func (u *TrainDao) GetTrain(trainID string) (*models.Train, error) {
	train := &models.Train{}
	err := u.DB.Where("train_id = ?", trainID).First(train).Error
	return train, err
}

func (u *TrainDao) ListTrains() ([]models.Train, error) {
	var trains []models.Train
	err := u.DB.Where("visibility = ?", true).Order("train_id").Find(&trains).Error
	return trains, err
}

func (u *TrainDao) AddTrain(train *models.Train) error {
	return u.DB.Create(train).Error
}

func (u *TrainDao) UpdateTrainVisibility(trainID string) error {
	return u.DB.Model(models.Train{}).Where("train_id = ?", trainID).UpdateColumn("visibility", false).Error
}

//func (u *TrainDao) DeleteTrain(trainID string) error {
//	return u.DB.Delete(models.Train{}, trainID).Error
//}
