package dao

import (
	"Project/MyProject/db"
	"Project/MyProject/models"
	"gorm.io/gorm"
	"time"
)

type TicketDaoImplement interface {
	GetTicket(ticketID int) (*models.Ticket, error)
	GetTickets(start, end string, date string) ([]models.Ticket, error)
	GetTicketsOnSale(start, end string, now time.Time) ([]models.Ticket, error)
	GetAllTickets() ([]models.Ticket, error)
	GetAllTicketsOnSale(now time.Time) ([]models.Ticket, error)
	AddTicket(ticket *models.Ticket) error
	UpdateTicketEndTime(ticketID int, now time.Time) error
	UpdateStockMinusOne(ticketID int) error
	UpdateStockAddOne(ticketID int) error
	//DeleteTicket(ticketID int) error
}

type TicketDao struct {
	DB *gorm.DB
}

func NewTicketDao() TicketDaoImplement {
	return &TicketDao{
		DB: db.DB,
	}
}

func (t *TicketDao) GetTicket(ticketID int) (*models.Ticket, error) {
	ticket := &models.Ticket{}
	err := t.DB.Where("ticket_id = ?", ticketID).First(ticket).Error
	return ticket, err
}

func (t *TicketDao) GetAllTickets() ([]models.Ticket, error) {
	tickets := make([]models.Ticket, 0)
	err := t.DB.Order("departure_time").Find(&tickets).Error
	return tickets, err
}

func (t *TicketDao) GetAllTicketsOnSale(now time.Time) ([]models.Ticket, error) {
	tickets := make([]models.Ticket, 0)
	err := t.DB.Order("departure_time").Find(&tickets, "start_time <= ? AND end_time >= ?", now, now).Error
	return tickets, err
}

func (t *TicketDao) GetTicketsOnSale(start, end string, now time.Time) ([]models.Ticket, error) {
	tickets := make([]models.Ticket, 0)
	err := t.DB.Order("departure_time").Find(&tickets, "start = ? AND end = ? AND start_time <= ? AND end_time >= ?", start, end, now, now).Error
	return tickets, err
}

func (t *TicketDao) GetTickets(start, end string, date string) ([]models.Ticket, error) {
	tickets := make([]models.Ticket, 0)
	now := time.Now()
	err := t.DB.Order("departure_time").Find(&tickets, "start = ? AND end = ? AND Date(departure_time) = ? AND departure_time > ? and start_time <= ? && end_time >= ?", start, end, date, now, now, now).Error
	return tickets, err
}

func (t *TicketDao) AddTicket(ticket *models.Ticket) error {
	return t.DB.Create(ticket).Error
}

func (t *TicketDao) UpdateTicketEndTime(ticketID int, now time.Time) error {
	return t.DB.Model(models.Ticket{}).Where("ticket_id = ?", ticketID).UpdateColumn("end_time", now).Error
}

func (t *TicketDao) UpdateStockMinusOne(ticketID int) error {
	return t.DB.Model(models.Ticket{}).
		Where("start_time < ?", time.Now()). // 筛出符合条件的记录，加快查询速度
		Where("ticket_id = ?", ticketID).
		UpdateColumn("stock", gorm.Expr("stock - ?", 1)).Error
}

func (t *TicketDao) UpdateStockAddOne(ticketID int) error {
	return t.DB.Model(models.Ticket{}).
		Where("start_time < ?", time.Now()). // 筛出符合条件的记录，加快查询速度
		Where("ticket_id = ?", ticketID).
		UpdateColumn("stock", gorm.Expr("stock + ?", 1)).Error
}

//func (t *TicketDao) DeleteTicket(ticketID int) error {
//	return t.DB.Delete(models.Ticket{}, ticketID).Error
//}
