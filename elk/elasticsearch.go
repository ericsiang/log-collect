package elk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/sirupsen/logrus"
)

type ElkSearchManager struct {
	Client *elasticsearch8.Client
}

func NewElkSearchManager(address []string, api_key string) (elkSearchManager *ElkSearchManager, err error) {
	cfg := elasticsearch8.Config{
		Addresses: address,
		APIKey:    api_key,
	}
	es8Client, err := elasticsearch8.NewClient(cfg)
	if err != nil {
		logrus.Error("ElasticSearch connect failed, err:", err)
		return nil, err
	}
	_, err = es8Client.Info()
	if err != nil {
		logrus.Error("ElasticSearch connect failed, err:", err)
	}
	logrus.Info("ElasticSearch connect success")
	// logrus.Info("ElasticSearch connect success, info:", info)
	return &ElkSearchManager{
		Client: es8Client,
	}, nil
}

func (e *ElkSearchManager) Process(topic string, msg *sarama.ConsumerMessage) error {
	//logrus.Infof("topic:%s, msg:%+v", topic, msg)
	type LogData struct {
		Title      string    `json:"title"`
		Content    string    `json:"Content"`
		CreateTime time.Time `json:"@timestamp"`
	}
	jsonData := map[string]interface{}{}
	err := json.Unmarshal(msg.Value, &jsonData)
	if err != nil {
		logrus.Error("json Unmarshal failed, err:", err)
		return err
	}
	// logrus.Infof("logDajsonDatata : %+v", jsonData)
	layout := "2006-01-02 15:04:05" // Go 的時間格式基於這個參考值
	time, err := time.Parse(layout, jsonData["time"].(string))
	if err != nil {
		logrus.Error("time.Parse failed, err:", err)
		return err
	}
	logData, err := json.Marshal(LogData{
		Title:      topic,
		Content:    string(msg.Value),
		CreateTime: time,
	})
	// logrus.Info("logData :", string(logData))
	if err != nil {
		logrus.Error("json Marshal failed, err:", err)
		return err
	}

	// 新增資料，如果索引不存在，則會自動建立索引
	req := esapi.IndexRequest{
		Index: topic,
		Body:  bytes.NewReader(logData),
	}
	logrus.Infof("req :%+v", req)
	res, err := req.Do(context.Background(), e.Client)
	if err != nil {
		logrus.Error("esapi IndexRequest failed, err:", err)
		return err
	}
	defer res.Body.Close()
	logrus.Info("res.Status() :", res.Status())
	if !strings.Contains(res.Status(), "201") {
		logrus.Error("esapi IndexRequest failed, res:", res)
		errMsg := fmt.Sprint("esapi IndexRequest failed, res:", res.Body)
		return errors.New(errMsg)
	}
	return nil
}
