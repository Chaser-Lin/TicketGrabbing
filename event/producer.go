package event

import (
	"github.com/Shopify/sarama"
	"log"
)

type Producer struct {
	kafkaProducer sarama.AsyncProducer
}

func NewProducer() (*Producer, error) {
	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//随机向partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//设置使用的kafka版本
	config.Version = sarama.V0_10_0_1

	//使用配置,新建一个异步生产者
	asyncProducer, err := sarama.NewAsyncProducer([]string{kafkaConf.Address}, config)
	if err != nil {
		log.Fatal("NewProducer NewAsyncProducer err: ", err)
		return nil, err
	}

	go func(p sarama.AsyncProducer) {
		for {
			select {
			case <-p.Successes():
			case fail := <-p.Errors():
				log.Println("AsyncProducer send message err: ", fail.Err)
			}
		}
	}(asyncProducer)

	return &Producer{
		kafkaProducer: asyncProducer,
	}, nil

}

func (p *Producer) SendMessage(message []byte) {
	msg := &sarama.ProducerMessage{
		Topic: kafkaConf.Topic,
	}
	msg.Value = sarama.ByteEncoder(message)
	//使用通道发送
	p.kafkaProducer.Input() <- msg
}
