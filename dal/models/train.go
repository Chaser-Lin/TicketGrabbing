package models

type Train struct {
	TrainID string `json:"train_id" gorm:"PrimaryKey;type:varchar(20)"`
	Speed   uint32 `json:"speed" gorm:"type:int"`
	Seats   uint32 `json:"seats" gorm:"type:int"`
}
