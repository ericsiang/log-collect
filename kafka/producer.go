package kafka

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/Shopify/sarama"
)

type KafkaProducerManager struct {
	producer sarama.SyncProducer
	msgChan  chan *sarama.ProducerMessage
}

func InitKafkaProducer(addr string, chanSize int64) (kafka *KafkaProducerManager, err error) {
	//初始化MsgChan
	msgChan := make(chan *sarama.ProducerMessage, chanSize)
	//初始化config
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          //发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner //创建随机分区
	config.Producer.Return.Successes = true                   //成功交付的消息将在success channel返回
	// logrus.Infof("kafka Producer config : %+v", config)

	//连接Kafka
	client, err := sarama.NewSyncProducer([]string{addr}, config)
	if err != nil {
		logrus.Error("kafka connect error , err : ", err)
		return nil, err
	}
	logrus.Info("kafka connect success")
	return &KafkaProducerManager{
		producer: client,
		msgChan:  msgChan,
	}, nil
}

func (k *KafkaProducerManager) SendToMsgChan(msg *sarama.ProducerMessage) {
	k.msgChan <- msg
}

func (k *KafkaProducerManager) Send(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			k.producer.Close()
			return
		case msg := <-k.msgChan:
			pid, offset, err := k.producer.SendMessage(msg)
			if err != nil {
				logrus.Error("send msg to kafka failed,err:", err)
			} else {
				logrus.Info("send msg to kafka success,pid:", pid, " , offset:", offset)
			}
		}
	}

}

func (k *KafkaProducerManager) Close() {
	k.producer.Close()
}
