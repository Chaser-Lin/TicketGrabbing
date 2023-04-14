package services

import (
	"Project/MyProject/dao"
	"Project/MyProject/models"
	"Project/MyProject/response"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// 添加列车服务参数
type AddTrainService struct {
	TrainID string `json:"train_id" form:"train_id" binding:"required"`
	Speed   uint32 `json:"speed" form:"speed" binding:"required,number"`
	Seats   uint32 `json:"seats" form:"seats" binding:"required,number"`
}

// 查询列车信息服务参数
type GetTrainService struct {
	TrainID string `json:"train_id" form:"train_id" binding:"required"`
}

// 列车相关服务接口：添加列车、展示所有列车、获取指定列车信息
type TrainServiceImplement interface {
	AddTrain(*AddTrainService) error
	ListTrains() ([]models.Train, error)
	//GetTrain(*GetTrainService) (*models.Train, error)
	GetTrain(trainID string) (*models.Train, error)
	DeleteTrain(trainID string) error
}

// 实现列车相关服务接口的实例
type TrainService struct {
	TrainDal dao.TrainDaoImplement
}

func (t *TrainService) DeleteTrain(trainID string) error {
	_, err := t.GetTrain(trainID)
	if err != nil {
		return err
	}
	err = t.TrainDal.UpdateTrainVisibility(trainID)
	if err != nil {
		return response.ErrDbOperation
	}
	return nil
}

func NewTrainServices(trainDal dao.TrainDaoImplement) TrainServiceImplement {
	return &TrainService{trainDal}
}

func (t *TrainService) AddTrain(service *AddTrainService) error {
	train := &models.Train{
		TrainID: service.TrainID,
		Speed:   service.Speed,
		Seats:   service.Seats,
	}

	if err := t.TrainDal.AddTrain(train); err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 { // 1062:Duplicate，重复数据
				return response.ErrTrainExist
			}
		}
		return response.ErrDbOperation
	}
	return nil
}

func (t *TrainService) ListTrains() ([]models.Train, error) {
	trains, err := t.TrainDal.ListTrains() // trains 中没有数据时 err == nil，会返回空数据切片
	if err == gorm.ErrRecordNotFound {
		return nil, response.EmptyTrainList
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	return trains, nil
}

//func (u *TrainService) GetTrain(services *GetTrainService) (*models.Train, error) {
//	return u.TrainDao.GetTrain(services.TrainID)
//}

func (t *TrainService) GetTrain(trainID string) (*models.Train, error) {
	train, err := t.TrainDal.GetTrain(trainID)
	if err == gorm.ErrRecordNotFound {
		return nil, response.ErrTrainNotExist
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	return train, nil
}
