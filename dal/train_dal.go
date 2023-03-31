package dal

import (
	"Project/MyProject/dal/models"
	"Project/MyProject/db"
	"gorm.io/gorm"
)

type TrainDalImplement interface {
	GetTrain(trainID string) (*models.Train, error)
	ListTrains() ([]models.Train, error)
	AddTrain(train *models.Train) error
	DeleteTrain(trainID string) error
}

type TrainDal struct {
	DB *gorm.DB
}

func NewTrainDal() TrainDalImplement {
	return &TrainDal{
		DB: db.DB,
	}
}

func (u *TrainDal) GetTrain(trainID string) (*models.Train, error) {
	train := &models.Train{}
	err := u.DB.Where("train_id = ?", trainID).First(train).Error
	return train, err
}

func (u *TrainDal) ListTrains() ([]models.Train, error) {
	var trains []models.Train
	err := u.DB.Find(&trains).Error
	return trains, err
}

func (u *TrainDal) AddTrain(train *models.Train) error {
	err := u.DB.Create(train).Error
	return err
}

func (u *TrainDal) DeleteTrain(trainID string) error {
	err := u.DB.Delete(models.Train{}, trainID).Error
	return err
}
