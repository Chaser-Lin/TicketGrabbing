package models

type Train struct {
	TrainID string `json:"train_id" gorm:"primaryKey;type:varchar(20)"`
	Speed   uint32 `json:"speed" gorm:"type:int;not null"`
	Seats   uint32 `json:"seats" gorm:"type:int;not null"`
}
