package service

import (
	"Project/MyProject/cache"
	"Project/MyProject/dal"
	"Project/MyProject/dal/models"
	"Project/MyProject/response"
	"Project/MyProject/utils"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

// 车票发售服务参数
type AddTicketService struct {
	// 通过参数传递得到的参数
	RouteID       int    `json:"route_id" form:"route_id" binding:"required,number"`
	TrainID       string `json:"train_id" form:"train_id" binding:"required"`
	Price         uint32 `json:"price" form:"price" binding:"required,number"`
	StartTime     string `json:"start_time" form:"start_time" binding:"required"`
	DepartureTime string `json:"departure_time" form:"departure_time" binding:"required"`

	// 通过routeID和TrainID获取对应信息后得到的参数
	Stock       uint32    `json:"stock" form:"stock"`
	Start       string    `json:"start" form:"start"`
	End         string    `json:"end" form:"end"`
	EndTime     time.Time `json:"end_time" form:"end_time"` // 结束购票时间，默认为列车发车时间
	ArrivalTime time.Time `json:"arrival_time" form:"arrival_time"`
}

// 获取指定路线车票服务参数
type ListTicketsService struct {
	Start string `json:"start" form:"start" binding:"required"`
	End   string `json:"end" form:"end" binding:"required"`
	Date  string `json:"date" form:"date" binding:"required"`
}

// 获取指定路线在售车票服务参数
type ListTicketsOnSaleService struct {
	Start string `json:"start" form:"start" binding:"required"`
	End   string `json:"end" form:"end" binding:"required"`
}

// 获取指定id车票服务参数
type GetTicketService struct {
	TicketID int `json:"ticket_id" form:"ticket_id" binding:"required,number"`
}

// 车票相关服务接口：车票发售、查询车票
type TicketServiceImplement interface {
	AddTicket(*AddTicketService) error
	GetTicket(ticketID int) (*models.Ticket, error)
	ListTickets(*ListTicketsService) ([]models.Ticket, error)
	ListTicketsOnSale(*ListTicketsOnSaleService) ([]models.Ticket, error)
	SubNumberOne(ticketID int) (err error)
	AddNumberOne(ticketID int) (err error)
	GetAllTickets() ([]models.Ticket, error)
	GetAllTicketsOnSale() ([]models.Ticket, error)
	//GetTicket(TicketID int) (*models.Ticket, error)
}

// 实现列车相关服务接口的实例
type TicketService struct {
	TicketDal dal.TicketDalImplement
}

func NewTicketServices(ticketDal dal.TicketDalImplement) TicketServiceImplement {
	return &TicketService{ticketDal}
}

func (t *TicketService) AddTicket(service *AddTicketService) error {
	//_, err := t.TicketDal.GetTicket(service.TicketID)
	//if err == nil {
	//	return response.ErrTicketExist
	//}

	startTime, err := utils.ParseStringToTime(service.StartTime)
	if err != nil {
		return response.ErrInvalidParam
	}
	departureTime, err := utils.ParseStringToTime(service.DepartureTime)

	//routeDal := dal.NewRouteDal()
	//route, err := routeDal.GetRouteByID(service.RouteID)
	//if err == gorm.ErrRecordNotFound {
	//	return response.ErrRouteNotExist
	//} else if err != nil {
	//	return response.ErrDbOperation
	//}

	ticket := &models.Ticket{
		RouteID:       service.RouteID,
		TrainID:       service.TrainID,
		Start:         service.Start,
		End:           service.End,
		Stock:         service.Stock,
		Price:         service.Price,
		StartTime:     startTime,
		EndTime:       departureTime,
		DepartureTime: departureTime,
		ArrivalTime:   service.ArrivalTime,
	}

	if err := t.TicketDal.AddTicket(ticket); err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 { // 1062:Duplicate，重复数据
				return response.ErrTicketExist
			}
		}
		return response.ErrDbOperation
	}
	// 加入redis缓存中
	err = cache.AddStock(cache.GetStockKey(ticket.TicketID), ticket.Stock)
	if err != nil {
		return err
	}

	return nil
}

func (t *TicketService) GetAllTickets() ([]models.Ticket, error) {
	tickets, err := t.TicketDal.GetAllTickets()
	if err == gorm.ErrRecordNotFound {
		return nil, response.EmptyTicketList
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	return tickets, nil
}

func (t *TicketService) GetAllTicketsOnSale() ([]models.Ticket, error) {
	tickets, err := t.TicketDal.GetAllTicketsOnSale(time.Now())
	if err == gorm.ErrRecordNotFound {
		return nil, response.EmptyTicketList
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	return tickets, nil
}

func (t *TicketService) ListTickets(service *ListTicketsService) ([]models.Ticket, error) {
	tickets, err := t.TicketDal.GetTickets(service.Start, service.End, service.Date)
	if err == gorm.ErrRecordNotFound {
		return nil, response.EmptyTicketList
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	return tickets, nil
}

func (t *TicketService) ListTicketsOnSale(service *ListTicketsOnSaleService) ([]models.Ticket, error) {
	tickets, err := t.TicketDal.GetTicketsOnSale(service.Start, service.End, time.Now())
	if err == gorm.ErrRecordNotFound {
		return nil, response.EmptyOnSaleTicketList
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	return tickets, nil
}

func (t *TicketService) GetTicket(ticketID int) (*models.Ticket, error) {
	ticket, err := t.TicketDal.GetTicket(ticketID)
	if err == gorm.ErrRecordNotFound {
		return nil, response.ErrTicketNotExist
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	return ticket, nil
}

func (t *TicketService) SubNumberOne(ticketID int) (err error) {
	err = t.TicketDal.UpdateStockMinusOne(ticketID)
	if err != nil {
		return response.ErrFailedSubStock
	}
	return
}

func (t *TicketService) AddNumberOne(ticketID int) (err error) {
	err = t.TicketDal.UpdateStockAddOne(ticketID)
	if err != nil {
		return response.ErrFailedSubStock
	}
	return
}

//func (u *TicketService) GetTicketByID(ticketID int) (*models.Ticket, error) {
//	return u.TicketDal.GetTicketByID(ticketID)
//}
