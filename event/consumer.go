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
