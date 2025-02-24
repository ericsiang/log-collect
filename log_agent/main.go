package main

import (
	"context"

	"github.com/sirupsen/logrus"

	"kafka-log/common"
	"kafka-log/config"
	"kafka-log/etcd"
	"kafka-log/kafka"
	"kafka-log/tail_file"
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
	// 在 etcd 新增 elasticsearch 用的連接參數
	putResp, err := etcdManager.Put(ctx, "elasticsearch_config", "{\"api_key\":\"djMzMUo1VUIwLWVVbldSUXJJRnY6cFRFY1dmSTBSUy1oeGgtbXBuVnZqUQ==\",\"url_address\":[\"http://localhost:9200\"]}")
	if err != nil {
		logrus.Error("etcd Put() failed, err:", err)
		return
	}
	logrus.Infof("etcd Put() success, putResp:%+v", putResp)

	// 從 etcd key 獲取的 value ，轉成 []common.CollectEntry
	collectEntryList, err := etcdManager.GetConfWithCollectEntry(ctx, configObj.Etcdaddress.Key)
	if err != nil {
		logrus.Error("etcd GetConfWithCollectEntry failed, err:", err)
		return
	}
	logrus.Infof("collectEntryList :%+v", collectEntryList)

	// 初始化 logData
	logData := common.NewLogData(collectEntryList)

	//初始化 Kafka Producer
	kafkaProducerManager, err := kafka.InitKafkaProducer(configObj.Kafakaddress.Address, configObj.Kafakaddress.MessageSize)
	if err != nil {
		logrus.Error("InitKafka failed, err:", err)
		return
	}
	go kafkaProducerManager.Send(ctx)

	// etcd 監聽 key
	go etcdManager.Watch(ctx, configObj.Etcdaddress.Key, logData)

	// 初始化 tail ，收集日志文件
	tailManager := tail_file.NewTailManager()
	go tailManager.ReloadInitTailTask(ctx, kafkaProducerManager, logData)
	tailManager.SendToTailChan() // 觸發 tail 流程執行
	// defer cancel()

	kafkaProducerManager.WgWait()
	etcdManager.WgWait()
	tailManager.WgWait()
}
