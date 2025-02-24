package main

import (
	"context"
	"fmt"
	"kafka-log/common"
	"kafka-log/config"
	"kafka-log/elk"
	"kafka-log/etcd"
	"kafka-log/kafka"

	"github.com/sirupsen/logrus"
)

var wrapError error

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05", // 设置json里的日期输出格式
	})

	// 初始化 config
	configObj, err := config.NewConfig()
	if err != nil {
		logrus.Error("newConfig failed, err:", err)
		return
	}
	logrus.Infof("configObj : %+v ", configObj)

	// 初始化 etcd
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

	//初始化 Kafka
	// kafkaConsumerManager, err := kafka.InitKafkaConsumerManager(configObj.Kafakaddress.Address)
	// if err != nil {
	// 	logrus.Error("InitKafka failed, err:", err)
	// 	return
	// }

	// 從 etcd  獲取 value
	getResp, err := etcdManager.Get(ctx, "elasticsearch_config")
	if err != nil {
		wrapError = fmt.Errorf("etcd Get() failed , err :%w", err)
		logrus.Error(wrapError)
		return
	}
	// logrus.Infof("etcd Get() success, getResp:%+v", getResp)

	// 初始化 ElasticConfig
	elasticConfig, err := common.NewElasticConfigFromEtcd(getResp)
	if err != nil {
		logrus.Error(wrapError)
		return
	}
	logrus.Infof("elasticConfig : %+v", elasticConfig)

	// 初始化 elk
	elkSearchManager, err := elk.NewElkSearchManager(elasticConfig.UrlAddress, elasticConfig.ApiKey)
	if err != nil {
		wrapError = fmt.Errorf("NewElkSearchManager failed , err :%w", err)
		logrus.Error(wrapError)
		return
	}
	logrus.Infof("collectEntryList : %+v", collectEntryList)
	for _, collectEntry := range collectEntryList {
		//初始化 Kafka group
		kafkaConsumerGroupManager, err := kafka.InitKafkaConsumerGroupManager([]string{configObj.Kafakaddress.Address}, collectEntry.Topic)
		if err != nil {
			logrus.Error("InitKafkaConsumerGroupManager failed, err:", err)
			continue
		}
		kafkaConsumerGroupManager.GetTopicData(ctx, collectEntry.Topic, elkSearchManager)
		kafkaConsumerGroupManager.WgWait()
	}

}
