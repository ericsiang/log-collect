package main

import (
	"context"
	"fmt"
	"kafka-log/common"
	"kafka-log/etcd"
	"kafka-log/kafka"

	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

func main() {
	// config ini
	var configObj = new(common.Config) //生成指针便于参数传递
	err := ini.MapTo(configObj, "config.ini")
	if err != nil {
		fmt.Println("log config failed,err:", err)
	}
	logrus.Infof("configObj : %+v ", configObj)

	// 初始化etcd
	etcdManager, err := etcd.NewEtcdManager(configObj.Etcdaddress.Address)
	if err != nil {
		logrus.Error("NewEtcdManager failed, err:", err)
		return
	}
	ctx, _ := context.WithCancel(context.Background())
	collectEntryList, err := etcdManager.GetConfWithCollectEntry(ctx, configObj.Etcdaddress.Key)
	if err != nil {
		logrus.Error("etcd GetConfWithCollectEntry failed, err:", err)
		return
	}
	logrus.Infof("collectEntryList :%+v", collectEntryList)

	//初始化Kafka
	kafkaConsumerManager, err := kafka.InitKafkaConsumerManager(configObj.Kafakaddress.Address, configObj.Kafakaddress.MessageSize)
	if err != nil {
		logrus.Error("InitKafka failed, err:", err)
		return
	}

}
