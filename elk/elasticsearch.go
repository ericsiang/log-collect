package elk

import (
	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
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
	// info, err := es8Client.Info()
	// if err != nil {
	// 	logrus.Error("ElasticSearch connect failed, err:", err)
	// }

	return &ElkSearchManager{
		Client: es8Client,
	}, nil
	// logrus.Info("ElasticSearch connect success, info:", info)
}

