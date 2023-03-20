package models

type Route struct {
	RouteID int    `json:"route_id" gorm:"PrimaryKey;autoIncrement:true"`
	Start   string `json:"start" gorm:"type:varchar(20);index:idx_start_end,unique"`
	End     string `json:"end" gorm:"type:varchar(20);index:idx_start_end,unique"`
	Length  uint32 `json:"length" gorm:"type:int"`
}
