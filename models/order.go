package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID          int       `json:"id" gorm:"type:int;primaryKey;autoIncrement:true"`
	OrderID     uuid.UUID `json:"order_id" gorm:"type:varchar(36);unique;not null"`
	TicketID    int       `json:"ticket_id" gorm:"not null"`
	UserID      int       `json:"user_id" gorm:"not null"`
	PassengerID int       `json:"passenger_id" gorm:"not null"`
	Status      int       `json:"status" gorm:"type:int;not null"`  // 订单状态：0/1/2/3：未支付/已支付/已过期/已取消
	Visibility  bool      `gorm:"visibility;not null;default:true"` // 订单可见性：用户删除订单后不可见
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	ExpiredAt   time.Time `json:"expired_at" gorm:"type:timestamp;not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP  on update current_timestamp;not null"`
}
