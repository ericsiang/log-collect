package kafka

import (
	"github.com/Shopify/sarama"
)

type KafkaConsumerManager struct {
	consumer sarama.Consumer
	// msgChan chan *sarama.ConsumerMessage
}

func InitKafkaConsumerManager(addr string, chanSize int64) (kafka *KafkaConsumerManager, err error) {
	//初始化MsgChan
	// msgChan := make(chan *sarama.ProducerMessage, chanSize)
	consumer, err := sarama.NewConsumer([]string{addr}, nil)
	if err != nil {
		panic(err)
	}
	

	return &KafkaConsumerManager{
		consumer: consumer,
	}, nil
}
