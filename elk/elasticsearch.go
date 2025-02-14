package elk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

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
	info, err := es8Client.Info()
	if err != nil {
		logrus.Error("ElasticSearch connect failed, err:", err)
	}
	logrus.Info("ElasticSearch connect success, info:", info)
	return &ElkSearchManager{
		Client: es8Client,
	}, nil
}

func (e *ElkSearchManager) SendToES(topic string, data []byte) error {
	logrus.Infof("topic:%s, data:%s", topic, data)
	type LogData struct {
		Title   string `json:"title"`
		Content string `json:"Content"`
	}

	logData, err := json.Marshal(LogData{
		Title:   topic,
		Content: string(data),
	})
	if err != nil {
		logrus.Error("json Marshal failed, err:", err)
		return err
	}

	// 新增資料，如果索引不存在，則會自動建立索引
	req := esapi.IndexRequest{
		Index: topic,
		Body:  bytes.NewReader(logData),
	}
	res, err := req.Do(context.Background(), e.Client)
	if err != nil {
		logrus.Error("esapi IndexRequest failed, err:", err)
		return err
	}
	defer res.Body.Close()
	if res.Status() != "201" {
		logrus.Error("esapi IndexRequest failed, res:", res)
		errMsg := fmt.Sprint("esapi IndexRequest failed, res:", res.Body)
		return errors.New(errMsg)
	}
	return nil
}
