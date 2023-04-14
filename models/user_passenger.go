package models

type UserPassenger struct {
	ID          int `json:"id" gorm:"type:int;primaryKey;autoIncrement:true"`
	UserID      int `json:"user_id" gorm:"type:int;index:idx_user_passenger,unique;not null"`
	PassengerID int `json:"passenger_id" gorm:"type:int;index:idx_user_passenger,unique;not null"`
}
