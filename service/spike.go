package service

import (
	"Project/MyProject/cache"
	"Project/MyProject/dal"
	"Project/MyProject/event"
	"Project/MyProject/response"
	"encoding/json"
	"log"
)

type SpikeServiceReq struct {
	TicketID int `json:"ticket_id" form:"ticket_id" binding:"required,number"`
}

type SpikeServiceImplement interface {
	BuyTicket(userID int, ticketID int) error
}

type SpikeService struct {
	KafkaProducer *event.Producer
	orderDal      dal.OrderDalImplement
}

func NewSpikeService(producer *event.Producer) SpikeServiceImplement {
	return &SpikeService{
		KafkaProducer: producer,
		orderDal:      dal.NewOrderDal(),
	}
}

func (s *SpikeService) BuyTicket(userID int, ticketID int) error {
	// 查询是否有余票
	if !cache.Limit(cache.GetStockKey(ticketID)) {
		return response.ErrNoRemainingTicket
	}

	// 如果用户已经购买过该车次的车票，直接返回
	exist := s.orderDal.IsValidOrderExist(userID, ticketID)
	if exist {
		// 把预扣的库存加回来
		err := cache.StockAddOne(cache.GetStockKey(ticketID))
		log.Printf("Selling cache.StockAddOne err: (%v)\n", err)
		return response.ErrSameOrderExist
	}

	message := MessageService{
		event.Message{
			TicketID: ticketID,
			UserID:   userID,
		},
	}
	byteMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Selling json.Marshal err: (%v)\n", err)
		err = cache.StockAddOne(cache.GetStockKey(ticketID))
		log.Printf("Selling cache.StockAddOne err: (%v)\n", err)
		return response.ErrFailedChangeToJson
	}

	s.KafkaProducer.SendMessage(byteMessage)

	return nil
}
