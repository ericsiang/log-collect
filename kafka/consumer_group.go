package kafka

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type KafkaConsumerGroupManager struct {
	consumerGroup   sarama.ConsumerGroup
	wg              *sync.WaitGroup
	consumerHandler *consumerHandler
}

type consumerHandler struct {
	consumerCustomProcess KafkaConsumerCustomProcessInterface
}

func InitKafkaConsumerGroupManager(addr []string, groupID string) (*KafkaConsumerGroupManager, error) {
	//初始化 sarama config
	config := sarama.NewConfig()
	config.Version= sarama.MaxVersion 
	consumerGroup, err := sarama.NewConsumerGroup(addr, groupID, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumerGroupManager{
		consumerGroup:   consumerGroup,
		wg:              &sync.WaitGroup{},
		consumerHandler: &consumerHandler{},
	}, nil
}

func (k *KafkaConsumerGroupManager) GetTopicData(ctx context.Context, topic string, consumerCustomProcess KafkaConsumerCustomProcessInterface) {
	k.wg.Add(1)
	k.consumerHandler.consumerCustomProcess = consumerCustomProcess
	for{
		err := k.consumerGroup.Consume(ctx, []string{topic}, k.consumerHandler)
		if err != nil {
			logrus.Error("failed group GetTopicData,err:", err)
			// return
		}
	}
	// return nil
}

func (k *KafkaConsumerGroupManager) WgAdd() {
	k.wg.Add(1)
}
func (k *KafkaConsumerGroupManager) WgDone() {
	k.wg.Done()
}

func (k *KafkaConsumerGroupManager) WgWait() {
	k.wg.Wait()
}

func (c *consumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg := <-claim.Messages():
			logrus.Infof("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
			c.consumerCustomProcess.Process(msg.Topic, msg) // 寫入 ElasticSearch
			session.MarkMessage(msg, "")                    // 標記訊息為已消費
		}
	}

}
