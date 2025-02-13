package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	"kafka-log/common"
	"kafka-log/elk"
	"kafka-log/etcd"
	"kafka-log/kafka"
	"kafka-log/tail_file"

	"github.com/go-ini/ini"
)

var wg sync.WaitGroup

var wrapError error

func main() {
	// config ini
	var configObj = new(common.Config) //生成指针便于参数传递
	err := ini.MapTo(configObj, "config.ini")
	if err != nil {
		wrapError = fmt.Errorf("ini MapTo() failed , err :%w", err)
		logrus.Error("log config failed,err:", wrapError)
		return
	}
	logrus.Infof("configObj : %+v ", configObj)

	// logFile, err := os.OpenFile(configObj.LogFilePath.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	panic(err)
	// }
	// defer logFile.Close()

	// 初始化etcd
	etcdManager, err := etcd.NewEtcdManager(configObj.Etcdaddress.Address)
	if err != nil {
		logrus.Error("NewEtcdManager failed, err:", err)
		return
	}
	ctx, _ := context.WithCancel(context.Background())
	putResp, err := etcdManager.Put(ctx, "elasticsearch_config", "{\"api_key\":\"bmZvWl9wUUI2S2ZDRkJDbi1mUGM6Ym1oRGVMSXNRWC1LU2Rpb0NtcHlBUQ==\",\"url_address\":[\"http://localhost:9200\"]}")
	if err != nil {
		logrus.Error("etcd Put() failed, err:", err)
		return
	}
	logrus.Infof("etcd Put() success, putResp:%+v", putResp)
	collectEntryList, err := etcdManager.GetConfWithCollectEntry(ctx, configObj.Etcdaddress.Key)
	if err != nil {
		logrus.Error("etcd GetConfWithCollectEntry failed, err:", err)
		return
	}
	logrus.Infof("collectEntryList :%+v", collectEntryList)
	wg.Add(1)
	go etcdManager.Watch(ctx, configObj.Etcdaddress.Key, &wg)

	//初始化Kafka
	kafkaProducerManager, err := kafka.InitKafkaProducer(configObj.Kafakaddress.Address, configObj.Kafakaddress.MessageSize)
	if err != nil {
		logrus.Error("InitKafka failed, err:", err)
		return
	}
	wg.Add(1)
	go kafkaProducerManager.Send(ctx, &wg)

	// 初始化tail
	// tailFile, err := tail_file.Init("log/error/error.log")
	err = tail_file.InitTail(ctx, &wg, kafkaProducerManager, collectEntryList)
	if err != nil {
		logrus.Error("InitTail failed, err:", err)
		return
	}
	logrus.Infof("InitTail success")
	// defer cancel()

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
	// add index
	response, err := elkSearchManager.Client.Indices.Create("server1-log2")
	if err != nil {
		wrapError = fmt.Errorf("Indices.Create failed , err :%w", err)
		logrus.Error(wrapError)
		return
	}
	logrus.Infof("response : %+v", response)

	wg.Wait()
}
