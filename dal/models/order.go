package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	//OrderID       uuid.UUID `json:"order_id" gorm:"PrimaryKey;type:varchar(255)"`
	//TicketID      int       `json:"ticket_id" gorm:"index:idx_ticket_owner,unique"`
	//UserID         int       `json:"owner" gorm:"type:int;index:idx_ticket_owner,unique"`
	//Passenger     string    `json:"passenger" gorm:"type:varchar(20)"`
	//Phone         string    `json:"phone" gorm:"type:varchar(11)"`
	//Price         uint32    `json:"price" gorm:"type:int"`
	//Start         string    `json:"start" gorm:"type:varchar(20)"`
	//End           string    `json:"end" gorm:"type:varchar(20)"`
	//TrainID       string    `json:"train_id" gorm:"type:varchar(20)"`
	//DepartureTime time.Time `json:"departure_time" gorm:"type:timestamp"`
	//ArrivalTime   time.Time `json:"arrival_time" gorm:"type:timestamp"`
	//Status        int       `json:"status" gorm:"type:int"` // 订单状态：0/1/2/3：未支付/已支付/已过期/已取消
	//Visibility    bool      `json:"visibility"`             // 订单可见性：用户删除订单后不可见
	//CreatedAt     time.Time `json:"created_at" gorm:"type:timestamp"`
	//ExpiredAt     time.Time `json:"expired_at" gorm:"type:timestamp"`

	OrderID    uuid.UUID `json:"order_id" gorm:"PrimaryKey;type:varchar(255)"`
	TicketID   int       `json:"ticket_id" gorm:"index"`
	UserID     int       `json:"user_id" gorm:"type:int"`
	User       User      `json:"user"`                   // gorm:"foreignkey:UserID"`
	Ticket     Ticket    `json:"ticket"`                 // gorm:"foreignkey:TicketID"`
	Status     int       `json:"status" gorm:"type:int"` // 订单状态：0/1/2/3：未支付/已支付/已过期/已取消
	Visibility bool      `json:"visibility"`             // 订单可见性：用户删除订单后不可见
	CreatedAt  time.Time `json:"created_at" gorm:"type:timestamp"`
	ExpiredAt  time.Time `json:"expired_at" gorm:"type:timestamp"`
}
