package event

import (
	"github.com/Shopify/sarama"
	"log"
)

type Consumer struct {
	KafkaConsumer sarama.Consumer
	Topic         string
}

func NewConsumer() (*Consumer, error) {
	consumer, err := sarama.NewConsumer([]string{kafkaConf.Address}, nil)
	if err != nil {
		log.Println("NewConsumer sarama.NewConsumer err:", err)
		return nil, err
	}

	return &Consumer{
		KafkaConsumer: consumer,
		Topic:         kafkaConf.Topic,
	}, nil

}

// TODO：还没写具体的业务逻辑
//func (c *Consumer) Start() error {
//	partitionList, err := c.KafkaConsumer.Partitions(c.Topic) // 根据topic取到所有的分区
//	if err != nil {
//		log.Println("NewConsumer consumer.Partitions err:", err)
//		return err
//	}
//
//	var wg sync.WaitGroup
//
//	for partition := range partitionList { // 遍历所有的分区
//		// 针对每个分区创建一个对应的分区消费者
//		pc, err := c.KafkaConsumer.ConsumePartition(c.Topic, int32(partition), sarama.OffsetNewest)
//		// 这里设置了OffsetNewest，只会收到consumer运行之后producer生产的数据
//		if err != nil {
//			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
//			return err
//		}
//		defer pc.AsyncClose()
//
//		wg.Add(1)
//		// 异步从每个分区消费信息
//		go func(sarama.PartitionConsumer) { // 为每个分区开一个go协程去取值
//			for msg := range pc.Messages() { // 阻塞直到有值发送过来，然后再继续等待
//				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, msg.Key, msg.Value)
//			}
//			wg.Done()
//		}(pc)
//	}
//	wg.Wait()
//	return nil
//}
