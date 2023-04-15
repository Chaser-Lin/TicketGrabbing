package services

import (
	"Project/MyProject/cache"
	"Project/MyProject/dao"
	"Project/MyProject/event"
	"Project/MyProject/models"
	"Project/MyProject/response"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"time"
)

// 添加订单服务参数
type AddOrderService struct {
	// 传入用户相关信息
	UserID      int    `json:"user_id" form:"user_id"`     // binding:"required,number"`
	Passenger   string `json:"passenger" form:"passenger"` // binding:"required"`
	Phone       string `json:"phone" form:"phone"`         // binding:"required,number"`
	PassengerID int    `json:"passenger_id" form:"passenger_id" binding:"required,number"`
	// 传入车票相关信息
	TicketID      int       `json:"ticket_id" form:"ticket_id" binding:"required,number"`
	Price         uint32    `json:"price" form:"price"`                   // binding:"required,number"`
	Start         string    `json:"start" form:"start"`                   // binding:"required"`
	End           string    `json:"end" form:"end"`                       // binding:"required"`
	TrainID       string    `json:"train_id" form:"train_id"`             // binding:"required"`
	DepartureTime time.Time `json:"departure_time" form:"departure_time"` // binding:"required,datetime"`
	ArrivalTime   time.Time `json:"arrival_time" form:"arrival_time"`     // binding:"required.datetime"`
	// 设置订单相关信息
	Status     int       `json:"status" form:"status"`         // 订单状态：0/1/2/3：未支付/已支付/已过期/已取消
	Visibility bool      `json:"visibility" form:"visibility"` // 订单可见性：用户删除订单后不可见
	CreatedAt  time.Time `json:"created_at" form:"created_at"`
	ExpiredAt  time.Time `json:"expired_at" form:"expired_at"`
}

// 向消息队列发送消息的服务
type MessageService struct {
	event.Message
}

// 用于展示给用户的订单信息
type OrderInfo struct {
	Passenger        string    `json:"passenger"`
	OrderID          string    `json:"order_id"`
	UserID           int       `json:"user_id"`
	TrainID          string    `json:"train_id"`
	Phone            string    `json:"phone"`
	Price            uint32    `json:"price"`
	Start            string    `json:"start"`
	End              string    `json:"end"`
	DepartureTimeDes string    `json:"departure_time_des"`
	ArrivalTimeDes   string    `json:"arrival_time_des"`
	DepartureTime    time.Time `json:"departure_time"`
	ArrivalTime      time.Time `json:"arrival_time"`
	Status           string    `json:"status"` // 订单状态：未支付/已支付/已过期/已取消
	CreatedAt        time.Time `json:"created_at"`
	ExpiredAt        time.Time `json:"expired_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

var orderStatus = []string{"未支付", "购票成功", "已过期", "已取消", "已完成", "行程中", "未出行"}

// 获取订单服务参数
type GetOrderService struct {
	OrderID string `json:"order_id" form:"order_id" binding:"required"`
}

// 获取订单列表服务参数
//type ListOrderService struct {
//	userID int `json:"user_id" form:"user_id" binding:"required,number"`
//}

// 更新订单状态服务参数
type UpdateOrderStatusService struct {
	OrderID string `json:"order_id" form:"order_id" binding:"required"`
	Status  int    `json:"status" form:"status"`
}

// 删除订单服务参数
type DeleteOrderService struct {
	OrderID string `json:"order_id" form:"order_id" binding:"required"`
}

// 订单相关服务接口：添加订单、展示所有订单、获取指定订单
type OrderServiceImplement interface {
	//AddOrder(*AddOrderService) error
	AddOrder(*MessageService) error
	GetOrder(orderID string) (*models.Order, error)  // 返回数据库中的订单消息
	GetOrderInfo(orderID string) (*OrderInfo, error) // 返回处理后的订单消息
	ListOrders(userID int) ([]*OrderInfo, error)
	UpdateOrderStatus(*UpdateOrderStatusService) error
	DeleteOrder(orderID string) error
	GetOrderPassengerAndTicketID(orderID string) (int, int, error)
	//DeleteOrder(services *DeleteOrderService) error
}

// 实现订单服务接口的实例
type OrderService struct {
	OrderDao     dao.OrderDaoImplement
	TicketDao    dao.TicketDaoImplement
	PassengerDao dao.PassengerDaoImplement
}

func NewOrderServices(orderDal dao.OrderDaoImplement,
	ticketDal dao.TicketDaoImplement,
	passengerDal dao.PassengerDaoImplement) OrderServiceImplement {
	return &OrderService{
		OrderDao:     orderDal,
		TicketDao:    ticketDal,
		PassengerDao: passengerDal,
	}
}

func (o *OrderService) AddOrder(service *MessageService) error {
	orderID := uuid.New()
	now := time.Now()

	order := &models.Order{
		OrderID:     orderID,
		UserID:      service.UserID,
		PassengerID: service.PassengerID,
		TicketID:    service.TicketID,
		Status:      0,
		ExpiredAt:   now.Add(30 * time.Minute),
	}
	log.Println(order)

	if err := o.OrderDao.AddOrder(order); err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 { // 1062:Duplicate，重复数据
				return response.ErrSameOrderExist
			}
		}
		return response.ErrDbOperation
	}

	// 将订单过期时间放入队列中
	orderExpireTime := time.Now().Add(30 * time.Minute).Unix()
	if err := cache.AddOrderExpireTime(float64(orderExpireTime), orderID.String()); err != nil {
		return response.ErrRedisOperation
	}

	return nil
}

func (o *OrderService) GetOrder(orderID string) (*models.Order, error) {
	order, err := o.OrderDao.GetOrder(orderID)
	if err == gorm.ErrRecordNotFound {
		return nil, response.ErrOrderNotExist
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	return order, nil
}

func (o *OrderService) GetOrderInfo(orderID string) (*OrderInfo, error) {
	order, err := o.GetOrder(orderID)
	if err != nil {
		return nil, err
	}
	ticket, err := o.TicketDao.GetTicket(order.TicketID)
	if err != nil {
		return nil, response.ErrDbOperation
	}
	passenger, err := o.PassengerDao.GetPassengerByID(order.PassengerID)
	if err != nil {
		return nil, response.ErrDbOperation
	}
	return parseOrderToInfo(order, ticket, passenger), nil
}

func (o *OrderService) GetOrderPassengerAndTicketID(orderID string) (int, int, error) {
	order, err := o.GetOrder(orderID)
	if err != nil {
		return 0, 0, err
	}
	return order.PassengerID, order.TicketID, nil
}

func (o *OrderService) ListOrders(userID int) ([]*OrderInfo, error) {
	orders, err := o.OrderDao.ListOrders(userID)
	if err == gorm.ErrRecordNotFound {
		return nil, response.EmptyOrderList
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	orderInfos := make([]*OrderInfo, 0)
	for _, order := range orders {
		ticket, err := o.TicketDao.GetTicket(order.TicketID)
		if err != nil {
			return nil, response.ErrDbOperation
		}
		passenger, err := o.PassengerDao.GetPassengerByID(order.PassengerID)
		if err != nil {
			return nil, response.ErrDbOperation
		}
		orderInfos = append(orderInfos, parseOrderToInfo(&order, ticket, passenger))
	}
	return orderInfos, nil
}

func (o *OrderService) UpdateOrderStatus(service *UpdateOrderStatusService) error {
	// 获取订单，判断订单是否存在
	_, err := o.GetOrder(service.OrderID)
	if err != nil {
		return err
	}
	// 更新订单状态
	err = o.OrderDao.UpdateOrderStatus(service.OrderID, service.Status)
	if err != nil {
		return response.ErrUpdateOrderStatus
	}
	// 更新redis缓存中的订单状态
	err = cache.RemoveFinishOrder(service.OrderID)
	if err != nil {
		return response.ErrRedisOperation
	}
	return nil
}

func (o *OrderService) DeleteOrder(orderID string) error {
	_, err := o.GetOrder(orderID)
	if err != nil {
		return err
	}
	err = o.OrderDao.UpdateOrderVisibility(orderID)
	if err != nil {
		return response.ErrDbOperation
	}
	return nil
}

func parseOrderToInfo(order *models.Order, ticket *models.Ticket, passenger *models.Passenger) *OrderInfo {
	status := orderStatus[order.Status]

	if order.Status == 1 {
		if time.Now().After(ticket.ArrivalTime) {
			status = orderStatus[4]
		} else if time.Now().After(ticket.DepartureTime) {
			status = orderStatus[5]
		} else {
			status = orderStatus[6]
		}
	}
	return &OrderInfo{
		OrderID:          order.OrderID.String(),
		TrainID:          ticket.TrainID,
		Passenger:        passenger.Name,
		Phone:            passenger.Phone,
		Price:            ticket.Price,
		Start:            ticket.Start,
		End:              ticket.End,
		DepartureTime:    ticket.DepartureTime,
		ArrivalTime:      ticket.ArrivalTime,
		DepartureTimeDes: fmt.Sprintf("%d月%02d日 %d:%02d出发", ticket.DepartureTime.Month(), ticket.DepartureTime.Day(), ticket.DepartureTime.Hour(), ticket.DepartureTime.Minute()), //ticket.DepartureTime,
		ArrivalTimeDes:   fmt.Sprintf("预计%d月%02d日 %d:%02d到达", ticket.ArrivalTime.Month(), ticket.ArrivalTime.Day(), ticket.ArrivalTime.Hour(), ticket.ArrivalTime.Minute()),       //ticket.ArrivalTime,
		Status:           status,
		CreatedAt:        order.CreatedAt,
		ExpiredAt:        order.ExpiredAt,
		UpdatedAt:        order.UpdatedAt,
	}
}
