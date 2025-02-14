package main

import (
	"context"
	"encoding/json"
	"fmt"
	"kafka-log/common"
	"kafka-log/elk"
	"kafka-log/etcd"
	"kafka-log/kafka"
	"sync"

	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

var wg sync.WaitGroup

var wrapError error

func main() {
	// config ini
	var configObj = new(common.Config) //生成指针便于参数传递
	err := ini.MapTo(configObj, "../config.ini")
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

	// 初始化elk
	getResp, err := etcdManager.Get(ctx, "elasticsearch_config")
	if err != nil {
		wrapError = fmt.Errorf("etcd Get() failed , err :%w", err)
		logrus.Error(wrapError)
		return
	}
	// logrus.Infof("etcd Get() success, getResp:%+v", getResp)
	elasticConfig := common.ElasticConfig{}
	// logrus.Infof("etcd Get() success, getResp.Kvs[0].Value :%+v", getResp.Kvs[0].Value)
	err = json.Unmarshal(getResp.Kvs[0].Value, &elasticConfig)
	if err != nil {
		wrapError = fmt.Errorf("NewElkSearchManager failed , err :%w", err)
		logrus.Error(wrapError)
		return
	}
	logrus.Infof("elasticConfig : %+v", elasticConfig)
	elkSearchManager, err := elk.NewElkSearchManager(elasticConfig.UrlAddress, elasticConfig.ApiKey)
	if err != nil {
		wrapError = fmt.Errorf("NewElkSearchManager failed , err :%w", err)
		logrus.Error(wrapError)
		return
	}

	for _, collectEntry := range collectEntryList {
		kafkaConsumerManager.GetTopicData(ctx, &wg, collectEntry.Topic, elkSearchManager)
	}
	wg.Wait()
}
