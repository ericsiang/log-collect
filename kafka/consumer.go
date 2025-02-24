package kafka

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type KafkaConsumerCustomProcessInterface interface {
	Process(topic string, msg *sarama.ConsumerMessage) error
}

type KafkaConsumerManager struct {
	consumer sarama.Consumer
	// msgChan chan *sarama.ConsumerMessage
	wg *sync.WaitGroup
}

func InitKafkaConsumerManager(addr string) (kafka *KafkaConsumerManager, err error) {
	consumer, err := sarama.NewConsumer([]string{addr}, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumerManager{
		consumer: consumer,
		wg:       &sync.WaitGroup{},
	}, nil
}

func (k *KafkaConsumerManager) GetTopicData(ctx context.Context, topic string, consumerCustomProcess KafkaConsumerCustomProcessInterface) (err error) {
	partitionlist, err := k.consumer.Partitions(topic) //获取topic的所有分区
	if err != nil {
		logrus.Error("failed get partition list,err:", err)
		return err
	}

	for partition := range partitionlist { // 遍历所有分区
		//根据消费者对象创建一个分区对象
		pc, err := k.consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			logrus.Error("failed get partition consumer,err:", err)
			return err
		}
		// 异步从每个分区消费信息
		k.wg.Add(1)
		go k.GetPartitionConsumerData(ctx, pc, consumerCustomProcess, topic)
	}
	return nil
}

func (k *KafkaConsumerManager) GetPartitionConsumerData(ctx context.Context, partition_consumer sarama.PartitionConsumer, consumerCustomProcess KafkaConsumerCustomProcessInterface, topic string) (err error) {
	logrus.Info("GetPartitionConsumerData() topic:", topic)
	for {
		select {
		case <-ctx.Done():
			k.wg.Done()
			k.Close()
			return
		case msg := <-partition_consumer.Messages(): // 從kafka通道中取出数据
			logrus.Infof("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
			consumerCustomProcess.Process(topic, msg)
		}
	}
}

func (k *KafkaConsumerManager) Close() {
	k.consumer.Close()
}

func (k *KafkaConsumerManager) WgAdd() {
	k.wg.Add(1)
}
func (k *KafkaConsumerManager) WgDone() {
	k.wg.Done()
}

func (k *KafkaConsumerManager) WgWait() {
	k.wg.Wait()
}
