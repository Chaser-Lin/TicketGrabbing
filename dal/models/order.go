package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	//gorm.Model
	OrderID     uuid.UUID  `json:"order_id" gorm:"type:varchar(255);unique;not null"`
	TicketID    int        `json:"ticket_id" gorm:"not null"`
	UserID      int        `json:"user_id" gorm:"not null"`
	PassengerID int        `json:"passenger_id" gorm:"not null"`
	Passenger   *Passenger `json:"passenger"`
	Ticket      *Ticket    `json:"ticket"`
	Status      int        `json:"status" gorm:"type:int;not null"` // 订单状态：0/1/2/3：未支付/已支付/已过期/已取消
	Visibility  bool       `json:"visibility;not null"`             // 订单可见性：用户删除订单后不可见
	ExpiredAt   time.Time  `json:"expired_at" gorm:"type:timestamp;not null"`
}
