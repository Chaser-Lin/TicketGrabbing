package services

import (
	"Project/MyProject/cache"
	"Project/MyProject/dao"
	"Project/MyProject/event"
	"Project/MyProject/response"
	"encoding/json"
	"log"
)

type SpikeServiceReq struct {
	TicketID    int `json:"ticket_id" form:"ticket_id" binding:"required,number"`
	PassengerID int `json:"passenger_id" form:"passenger_id" binding:"required,number"`
	//UserID   int `json:"user_id" form:"user_id" binding:"required,number"`
}

type SpikeServiceImplement interface {
	BuyTicket(userID int, ticketID int, passengerID int) error
}

type SpikeService struct {
	KafkaProducer *event.Producer
	orderDao      dao.OrderDaoImplement
}

func NewSpikeService(producer *event.Producer) SpikeServiceImplement {
	return &SpikeService{
		KafkaProducer: producer,
		orderDao:      dao.NewOrderDao(),
	}
}

func (s *SpikeService) BuyTicket(userID int, ticketID int, passengerID int) error {
	// 查询是否有余票
	if !cache.Limit(cache.GetStockKey(ticketID)) {
		return response.ErrNoRemainingTicket
	}

	// 若同一时间有多个相同的请求进来，在orderlimit判断通过的情况下，实际上依然有可能会有两个相同的请求（两个请求的间隔快到连redis也来不及更新）
	// 就会导致某个用户购买多张相同车票的情况，对于这种情况在kafka消费者中进行了处理，对于重复请求会回退库存，导致出现少卖的情况
	// 如果用户已经购买过该车次的车票，直接返回
	exist, err := cache.OrderLimit(passengerID, ticketID)
	if err != nil {
		log.Printf("BuyTicket cache.OrderLimit err: (%v)\n", err)
		return response.ErrRedisOperation
	}
	if exist {
		// 把预扣的库存加回来
		log.Println("before send message")
		if err = cache.StockAddOne(cache.GetStockKey(ticketID)); err != nil {
			log.Printf("BuyTicket cache.StockAddOne err: (%v)\n", err)
		}
		return response.ErrSameOrderExist
	}

	message := MessageService{
		event.Message{
			TicketID:    ticketID,
			UserID:      userID,
			PassengerID: passengerID,
		},
	}
	byteMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("BuyTicket json.Marshal err: (%v)\n", err)
		err = cache.StockAddOne(cache.GetStockKey(ticketID))
		log.Printf("BuyTicket cache.StockAddOne err: (%v)\n", err)
		return response.ErrFailedChangeToJson
	}

	s.KafkaProducer.SendMessage(byteMessage)

	return nil
}
