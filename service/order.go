package service

import (
	"Project/MyProject/cache"
	"Project/MyProject/dal"
	"Project/MyProject/dal/models"
	"Project/MyProject/event"
	"Project/MyProject/response"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// 添加订单服务参数
type AddOrderService struct {
	// 传入用户相关信息
	UserID    int    `json:"user_id" form:"user_id"`     // binding:"required,number"`
	Passenger string `json:"passenger" form:"passenger"` // binding:"required"`
	Phone     string `json:"phone" form:"phone"`         // binding:"required,number"`
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
	Passenger     string    `json:"passenger"`
	OrderID       string    `json:"order_id"`
	UserID        int       `json:"user_id"`
	UserEmail     string    `json:"user_email"`
	Phone         string    `json:"phone"`
	Price         uint32    `json:"price"`
	Start         string    `json:"start"`
	End           string    `json:"end"`
	DepartureTime time.Time `json:"departure_time"`
	ArrivalTime   time.Time `json:"arrival_time"`
	Status        int       `json:"status"` // 订单状态：0/1/2/3：未支付/已支付/已过期/已取消
	CreatedAt     time.Time `json:"created_at"`
	ExpiredAt     time.Time `json:"expired_at"`
}

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
	GetOrderUserAndTicketID(orderID string) (int, int, error)
	//DeleteOrder(service *DeleteOrderService) error
}

// 实现订单服务接口的实例
type OrderService struct {
	OrderDal dal.OrderDalImplement
}

func NewOrderServices(orderDal dal.OrderDalImplement) OrderServiceImplement {
	return &OrderService{orderDal}
}

func (o *OrderService) AddOrder(service *MessageService) error {
	//_, err := o.OrderDal.GetOrder(service.Start, service.End)
	//if err == nil {
	//	return response.ErrOrderExist
	//}

	orderID := uuid.New()
	now := time.Now()

	//userDal := dal.NewUserDal()
	//user, err := userDal.GetUserByID(service.UserID)
	//if err != nil {
	//	return err
	//}
	//
	//ticketDal := dal.NewTicketDal()
	//ticket, err := ticketDal.GetTicket(service.TicketID)
	//if err != nil {
	//	return err
	//}

	order := &models.Order{
		OrderID:    orderID,
		UserID:     service.UserID,
		TicketID:   service.TicketID,
		Status:     0,
		Visibility: true,
		CreatedAt:  now,
		ExpiredAt:  now.Add(30 * time.Minute),
	}

	if err := o.OrderDal.AddOrder(order); err != nil {
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

//
//func (o *OrderService) AddOrder(service *AddOrderService) error {
//	//_, err := o.OrderDal.GetOrder(service.Start, service.End)
//	//if err == nil {
//	//	return response.ErrOrderExist
//	//}
//
//	orderID := uuid.New()
//	now := time.Now()
//
//	//userDal := dal.NewUserDal()
//	//user, err := userDal.GetUserByID(service.UserID)
//	//if err != nil {
//	//	return err
//	//}
//	//
//	//ticketDal := dal.NewTicketDal()
//	//ticket, err := ticketDal.GetTicket(service.TicketID)
//	//if err != nil {
//	//	return err
//	//}
//
//	order := &models.Order{
//		OrderID:       orderID,
//		UserID:         service.UserID,
//		Passenger:     service.Passenger,
//		Phone:         service.Phone,
//		TicketID:      service.TicketID,
//		Price:         service.Price,
//		Start:         service.Start,
//		End:           service.End,
//		TrainID:       service.TrainID,
//		DepartureTime: service.DepartureTime,
//		ArrivalTime:   service.ArrivalTime,
//		Status:        0,
//		Visibility:    true,
//		CreatedAt:     now,
//		ExpiredAt:     now.Add(30 * time.Minute),
//	}
//
//	if err := o.OrderDal.AddOrder(order); err != nil {
//		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
//			if mysqlErr.Number == 1062 { // 1062:Duplicate，重复数据
//				return response.ErrSameOrderExist
//			}
//		}
//		return response.ErrDbOperation
//	}
//	return nil
//}

func (o *OrderService) GetOrder(orderID string) (*models.Order, error) {
	order, err := o.OrderDal.GetOrder(orderID)
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
	return parseOrderToInfo(order), nil
}

func (o *OrderService) GetOrderUserAndTicketID(orderID string) (int, int, error) {
	order, err := o.GetOrder(orderID)
	if err != nil {
		return 0, 0, err
	}
	return order.UserID, order.TicketID, nil
}

func (o *OrderService) ListOrders(userID int) ([]*OrderInfo, error) {
	orders, err := o.OrderDal.ListOrders(userID)
	if err == gorm.ErrRecordNotFound {
		return nil, response.EmptyOrderList
	} else if err != nil {
		return nil, response.ErrDbOperation
	}
	orderInfos := make([]*OrderInfo, 0)
	for _, order := range orders {
		orderInfos = append(orderInfos, parseOrderToInfo(&order))
	}
	return orderInfos, nil
}

func (o *OrderService) UpdateOrderStatus(service *UpdateOrderStatusService) error {
	_, err := o.GetOrder(service.OrderID)
	if err != nil {
		return err
	}
	err = o.OrderDal.UpdateOrderStatus(service.OrderID, service.Status)
	if err != nil {
		return response.ErrUpdateOrderStatus
	}

	// 更新订单状态时，同时更新redis缓存中的订单状态
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
	err = o.OrderDal.UpdateOrderVisibility(orderID)
	if err != nil {
		return response.ErrDbOperation
	}
	return nil
}

func parseOrderToInfo(order *models.Order) *OrderInfo {
	return &OrderInfo{
		Passenger:     order.User.Username,
		OrderID:       order.OrderID.String(),
		UserID:        order.UserID,
		UserEmail:     order.User.Email,
		Phone:         order.User.Phone,
		Price:         order.Ticket.Price,
		Start:         order.Ticket.Start,
		End:           order.Ticket.End,
		DepartureTime: order.Ticket.DepartureTime,
		ArrivalTime:   order.Ticket.ArrivalTime,
		Status:        order.Status,
		CreatedAt:     order.CreatedAt,
		ExpiredAt:     order.ExpiredAt,
	}
}
