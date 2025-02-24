package common

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewElasticConfigFromEtcd(etcdGetResp *clientv3.GetResponse) (elasticConfig *ElasticConfig, err error) {
	elasticConfig = &ElasticConfig{}
	// logrus.Infof("etcd Get() success, getResp.Kvs[0].Value :%+v", getResp.Kvs[0].Value)
	err = json.Unmarshal(etcdGetResp.Kvs[0].Value, &elasticConfig)
	if err != nil {
		wrapError := fmt.Errorf("NewElkSearchManager failed , err :%w", err)
		logrus.Error(wrapError)
		return nil, err
	}
	return elasticConfig, nil
}
