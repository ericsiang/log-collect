package kafka

import (
	"context"
	"kafka-log/elk"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
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

func (k *KafkaConsumerManager) GetTopicData(ctx context.Context, wg *sync.WaitGroup, topic string, elk *elk.ElkSearchManager) (err error) {
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
		wg.Add(1)
		go k.GetPartitionConsumerData(ctx, wg, pc, elk, topic)
	}
	return nil
}

func (k *KafkaConsumerManager) GetPartitionConsumerData(ctx context.Context, wg *sync.WaitGroup, partition_consumer sarama.PartitionConsumer, elk *elk.ElkSearchManager, topic string) (err error) {
	logrus.Info("GetPartitionConsumerData() topic:", topic)
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			k.Close()
			return
		case msg := <-partition_consumer.Messages(): // 從kafka通道中取出数据
			logrus.Infof("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
			elk.SendToES(topic, msg.Value)
		}
	}
}

func (k *KafkaConsumerManager) Close() {
	k.consumer.Close()
}
