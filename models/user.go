package models

type User struct {
	UserID         int    `json:"user_id" gorm:"type:int;primaryKey;autoIncrement:true"`
	Username       string `json:"username" gorm:"type:varchar(50);not null"`
	HashedPassword string `json:"hashed_password" gorm:"type:varchar(64);not null"`
	Email          string `json:"email" gorm:"unique;type:varchar(50);not null"`
}
