package models

import "time"

type Ticket struct {
	TicketID      int       `json:"ticket_id" gorm:"PrimaryKey;autoIncrement:true"`
	RouteID       int       `json:"route_id" gorm:"index:idx_route_departure-time,unique"`
	TrainID       string    `json:"train_id" gorm:"type:varchar(20)"`
	Start         string    `json:"start" gorm:"type:varchar(20)"`
	End           string    `json:"end" gorm:"type:varchar(20)"`
	Stock         uint32    `json:"stock" gorm:"type:int"`
	Price         uint32    `json:"price" gorm:"type:int"`
	StartTime     time.Time `json:"start_time" gorm:"type:timestamp"`
	EndTime       time.Time `json:"end_time" gorm:"type:timestamp"`
	DepartureTime time.Time `json:"departure_time" gorm:"type:timestamp;index:idx_route_departure-time,unique"`
	ArrivalTime   time.Time `json:"arrival_time" gorm:"type:timestamp"`
}
