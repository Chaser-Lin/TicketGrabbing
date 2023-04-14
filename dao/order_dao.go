package dao

import (
	"Project/MyProject/db"
	"Project/MyProject/models"
	"gorm.io/gorm"
)

type OrderDaoImplement interface {
	ListOrders(userID int) ([]models.Order, error)
	GetOrder(orderID string) (*models.Order, error)
	AddOrder(order *models.Order) error
	UpdateOrderVisibility(orderID string) error
	UpdateOrderStatus(orderID string, status int) error
	IsValidOrderExist(userID int, ticketID int) bool
}

type OrderDao struct {
	DB *gorm.DB
}

func NewOrderDao() OrderDaoImplement {
	return &OrderDao{
		DB: db.DB,
	}
}

func (o *OrderDao) ListOrders(userID int) ([]models.Order, error) {
	var orders []models.Order
	err := o.DB.Order("id desc").Find(&orders, "user_id = ? AND visibility = ?", userID, true).Error
	return orders, err
}

func (o *OrderDao) GetOrder(orderID string) (*models.Order, error) {
	order := &models.Order{}
	err := o.DB.Where("order_id = ? AND visibility = ?", orderID, true).First(order).Error
	return order, err
}

func (o *OrderDao) AddOrder(order *models.Order) error {
	return o.DB.Create(order).Error
}

func (o *OrderDao) UpdateOrderVisibility(orderID string) error {
	return o.DB.Model(models.Order{}).Where("order_id = ?", orderID).UpdateColumn("visibility", false).Error
}

func (o *OrderDao) UpdateOrderStatus(orderID string, status int) error {
	return o.DB.Model(models.Order{}).Where("order_id = ?", orderID).UpdateColumn("status", status).Error
}

func (o *OrderDao) IsValidOrderExist(userID int, ticketID int) bool {
	n := o.DB.First(&models.Order{}, "user_id = ? AND ticket_id = ? AND status <= ?", userID, ticketID, 1).RowsAffected
	return n > 0
}
