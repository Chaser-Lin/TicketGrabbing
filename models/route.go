package models

type Route struct {
	RouteID    int    `json:"route_id" gorm:"type:int;primaryKey;autoIncrement:true"`
	Start      string `json:"start" gorm:"type:varchar(20);not null"` // index:idx_start_end,unique;
	End        string `json:"end" gorm:"type:varchar(20);not null"`   // index:idx_start_end,unique;
	Length     uint32 `json:"length" gorm:"type:int;not null"`
	Visibility bool   `gorm:"visibility;not null;default:true"` // 路线可见性：管理员删除路线后不可见
}
