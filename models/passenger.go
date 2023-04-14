package models

type Passenger struct {
	PassengerID int    `json:"passenger_id" gorm:"type:int;primaryKey;autoIncrement:true"`
	Name        string `json:"name" gorm:"type:varchar(50);not null"`
	IDNumber    string `json:"id_number" gorm:"unique;type:varchar(18);not null"`
	Phone       string `json:"phone" gorm:"type:varchar(11);not null"`
}
