package services

import (
	"Project/MyProject/cache"
	"Project/MyProject/event"
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
	"sync"
)

type KafkaMQService struct {
	*event.Consumer
}

func NewKafkaMQService(consumer *event.Consumer) *KafkaMQService {
	return &KafkaMQService{consumer}
}

func (k *KafkaMQService) StartConsumer(orderService OrderServiceImplement, ticketService TicketServiceImplement) {

	partitionList, err := k.KafkaConsumer.Partitions(k.Topic) // 根据topic取到所有的分区
	if err != nil {
		log.Println("StartConsumer get Partitions err:", err)
		//return err
	}

	var wg sync.WaitGroup

	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := k.KafkaConsumer.ConsumePartition(k.Topic, int32(partition), sarama.OffsetNewest)
		// 这里设置了OffsetNewest，只会收到consumer运行之后producer生产的数据
		if err != nil {
			log.Println("StartConsumer ConsumePartition err:", err)
			//return err
		}
		defer pc.AsyncClose()

		wg.Add(1)
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) { // 为每个分区开一个go协程去取值
			for msg := range pc.Messages() { // 阻塞直到有值发送过来，然后再继续等待
				//fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, msg.Key, msg.Value)
				message := &MessageService{}
				err = json.Unmarshal(msg.Value, message)
				if err != nil {
					log.Println("StartConsumer Unmarshal msg.Value err:", err)
					continue
				}

				// 补偿策略：在同一时间内存在多个相同请求的情况，将库存加回来
				exist, err := cache.OrderLimit(message.PassengerID, message.TicketID)
				if err != nil {
					log.Println("StartConsumer cache.OrderLimit err:", err)
					continue
				}
				if exist {
					// 把预扣的库存加回来
					log.Printf("after send message, userID: %d, ticketID: %d\n", message.PassengerID, message.TicketID)
					if err = cache.StockAddOne(cache.GetStockKey(message.TicketID)); err != nil {
						log.Println("StartConsumer cache.StockAddOne err:", err)
					}
					continue
				}

				err = orderService.AddOrder(message)
				if err != nil {
					log.Println("StartConsumer orderService.AddOrder err:", err)
					continue
				}
				err = cache.AddOrderLimit(message.PassengerID, message.TicketID)
				if err != nil {
					log.Println("StartConsumer cache.AddOrderLimit err:", err)
				}

				// 如果减库存失败，但是redis中已经预扣成功，不会导致超卖问题
				err = ticketService.SubNumberOne(message.TicketID)
				if err != nil {
					log.Println("StartConsumer ticketService.SubNumberOne err:", err)
				}
			}
			wg.Done()
		}(pc)
	}
	wg.Wait()
	//return nil
}
