package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

func main() {
	fmt.Println("kafka")
	config := sarama.NewConfig()    
	config.Producer.RequiredAcks = sarama.WaitForAll          //发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner //创建随机分区
	config.Producer.Return.Successes = true                   //成功交付的消息将在success channel返回

	//创建信息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "web.log"
	msg.Value = sarama.StringEncoder("this is a test log")

	//连接KafKa
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	defer client.Close()

	//发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed,err:", err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)
}