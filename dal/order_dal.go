package dal

import (
	"Project/MyProject/dal/models"
	"Project/MyProject/db"
	"gorm.io/gorm"
)

type OrderDalImplement interface {
	ListOrders(userID int) ([]models.Order, error)
	GetOrder(orderID string) (*models.Order, error)
	AddOrder(order *models.Order) error
	UpdateOrderVisibility(orderID string) error
	UpdateOrderStatus(orderID string, status int) error
	IsValidOrderExist(userID int, ticketID int) bool
}

type OrderDal struct {
	DB *gorm.DB
}

func NewOrderDal() OrderDalImplement {
	return &OrderDal{
		DB: db.DB,
	}
}

func (o *OrderDal) ListOrders(userID int) ([]models.Order, error) {
	var orders []models.Order
	err := o.DB.Preload("Ticket").Preload("Passenger").Find(&orders, "user_id = ? AND visibility = ?", userID, true).Error
	return orders, err
}

func (o *OrderDal) GetOrder(orderID string) (*models.Order, error) {
	order := &models.Order{}
	err := o.DB.Preload("Ticket").Preload("Passenger").Where("order_id = ? AND visibility = ?", orderID, true).First(order).Error
	return order, err
}

func (o *OrderDal) AddOrder(order *models.Order) error {
	err := o.DB.Create(order).Error
	return err
}

func (o *OrderDal) UpdateOrderVisibility(orderID string) error {
	err := o.DB.Model(models.Order{}).Where("order_id = ?", orderID).Update("visibility", false).Error
	return err
}

func (o *OrderDal) UpdateOrderStatus(orderID string, status int) error {
	err := o.DB.Model(models.Order{}).Where("order_id = ?", orderID).Update("status", status).Error
	return err
}

func (o *OrderDal) IsValidOrderExist(userID int, ticketID int) bool {
	n := o.DB.First(&models.Order{}, "user_id = ? AND ticket_id = ? AND status <= ?", userID, ticketID, 1).RowsAffected
	return n > 0
}
