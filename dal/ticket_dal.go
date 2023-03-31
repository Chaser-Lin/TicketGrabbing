package dal

import (
	"Project/MyProject/dal/models"
	"Project/MyProject/db"
	"gorm.io/gorm"
	"time"
)

type TicketDalImplement interface {
	GetTicket(ticketID int) (*models.Ticket, error)
	AddTicket(ticket *models.Ticket) error
	DeleteTicket(ticketID int) error
	GetTickets(start, end string, date string) ([]models.Ticket, error)
	GetTicketsOnSale(start, end string, now time.Time) ([]models.Ticket, error)
	GetAllTickets() ([]models.Ticket, error)
	GetAllTicketsOnSale(now time.Time) ([]models.Ticket, error)
	UpdateStockMinusOne(ticketID int) error
	UpdateStockAddOne(ticketID int) error
}

type TicketDal struct {
	DB *gorm.DB
}

func NewTicketDal() TicketDalImplement {
	return &TicketDal{
		DB: db.DB,
	}
}

func (t *TicketDal) GetTicket(ticketID int) (*models.Ticket, error) {
	ticket := &models.Ticket{}
	err := t.DB.Where("ticket_id = ?", ticketID).First(ticket).Error
	return ticket, err
}

func (t *TicketDal) GetAllTickets() ([]models.Ticket, error) {
	tickets := make([]models.Ticket, 0)
	err := t.DB.Find(&tickets).Error
	return tickets, err
}

func (t *TicketDal) GetAllTicketsOnSale(now time.Time) ([]models.Ticket, error) {
	tickets := make([]models.Ticket, 0)
	err := t.DB.Find(&tickets, "start_time <= ? AND end_time >= ?", now, now).Error
	return tickets, err
}

func (t *TicketDal) GetTicketsOnSale(start, end string, now time.Time) ([]models.Ticket, error) {
	tickets := make([]models.Ticket, 0)
	err := t.DB.Find(&tickets, "start = ? AND end = ? AND start_time <= ? AND end_time >= ?", start, end, now, now).Error
	return tickets, err
}

func (t *TicketDal) GetTickets(start, end string, date string) ([]models.Ticket, error) {
	tickets := make([]models.Ticket, 0)
	err := t.DB.Find(&tickets, "start = ? AND end = ? AND Date(departure_time) = ?", start, end, date).Error
	return tickets, err
}

func (t *TicketDal) AddTicket(ticket *models.Ticket) error {
	err := t.DB.Create(ticket).Error
	return err
}

func (t *TicketDal) DeleteTicket(ticketID int) error {
	err := t.DB.Delete(models.Ticket{}, ticketID).Error
	return err
}

func (t *TicketDal) UpdateStockMinusOne(ticketID int) error {
	err := t.DB.Model(models.Ticket{}).
		Where("start_time < ?", time.Now()). // 筛出符合条件的记录，加快查询速度
		Where("ticket_id = ?", ticketID).
		UpdateColumn("stock", gorm.Expr("stock - ?", 1)).Error
	return err
}

func (t *TicketDal) UpdateStockAddOne(ticketID int) error {
	err := t.DB.Model(models.Ticket{}).
		Where("start_time < ?", time.Now()). // 筛出符合条件的记录，加快查询速度
		Where("ticket_id = ?", ticketID).
		UpdateColumn("stock", gorm.Expr("stock + ?", 1)).Error
	return err
}
