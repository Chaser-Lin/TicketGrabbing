package models

type Train struct {
	ID         int    `json:"id" gorm:"type:int;primaryKey;autoIncrement:true"`
	TrainID    string `json:"train_id" gorm:"unique;type:varchar(20)"`
	Speed      uint32 `json:"speed" gorm:"type:int;not null"`
	Seats      uint32 `json:"seats" gorm:"type:int;not null"`
	Visibility bool   `gorm:"visibility;not null;default:true"` // 列车可见性：管理员删除列车后不可见
}
