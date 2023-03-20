package models

type User struct {
	UserID         int    `json:"user_id" gorm:"primaryKey;autoIncrement:true"`
	Username       string `json:"username" gorm:"type:varchar(50)"`
	HashedPassword string `json:"hashed_password" gorm:"type:varchar(64)"`
	Email          string `json:"email" gorm:"unique;type:varchar(50)"`
	Phone          string `json:"phone" gorm:"type:varchar(11)"`
}
