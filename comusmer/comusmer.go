package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

func main() {
	customer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		fmt.Println("failed init customer,err:", err)
		return
	}

	partitionlist, err := customer.Partitions("web.log-0") //获取topic的所有分区
	if err != nil {
		fmt.Println("failed get partition list,err:", err)
		return
	}

	fmt.Println("partitions:", partitionlist)

	for partition := range partitionlist { // 遍历所有分区
		//根据消费者对象创建一个分区对象
		pc, err := customer.ConsumePartition("web.log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Println("failed get partition consumer,err:", err)
			return
		}
		defer pc.Close() // 移动到这里

		go func(consumer sarama.PartitionConsumer) {
			
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(pc)
		time.Sleep(time.Second * 10)
	}
	// defer pc.AsyncClose() // 移除这行，因为已经在循环结束时关闭了
}