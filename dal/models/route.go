package models

type Route struct {
	RouteID int    `json:"route_id" gorm:"type:int;primaryKey;autoIncrement:true"`
	Start   string `json:"start" gorm:"type:varchar(20);index:idx_start_end,unique;not null"`
	End     string `json:"end" gorm:"type:varchar(20);index:idx_start_end,unique;not null"`
	Length  uint32 `json:"length" gorm:"type:int;not null"`
}
