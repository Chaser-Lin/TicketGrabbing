package models

import "time"

type Ticket struct {
	TicketID      int       `json:"ticket_id" gorm:"type:int;primaryKey;autoIncrement:true"`
	RouteID       int       `json:"route_id" gorm:"index:idx_route_departure-time,unique;not null"`
	TrainID       string    `json:"train_id" gorm:"type:varchar(20);not null"`
	Start         string    `json:"start" gorm:"type:varchar(20);not null"`
	End           string    `json:"end" gorm:"type:varchar(20);not null"`
	Stock         uint32    `json:"stock" gorm:"type:int;not null"`
	Price         uint32    `json:"price" gorm:"type:int;not null"`
	StartTime     time.Time `json:"start_time" gorm:"type:timestamp;not null"`
	EndTime       time.Time `json:"end_time" gorm:"type:timestamp;not null"`
	DepartureTime time.Time `json:"departure_time" gorm:"type:timestamp;index:idx_route_departure-time,unique;not null"`
	ArrivalTime   time.Time `json:"arrival_time" gorm:"type:timestamp;not null"`
}
