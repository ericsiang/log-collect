package etcd

import (
	"context"
	"encoding/json"
	"kafka-log/common"
	"kafka-log/tail_file"

	"sync"
	"time"

	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdManager struct {
	client *clientv3.Client
	wg     *sync.WaitGroup
}

func NewEtcdManager(address []string) (etcdManager *EtcdManager, err error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   address,         // etcd集群地址
		DialTimeout: 5 * time.Second, // 连接超时时间
	})
	if err != nil {
		logrus.Error("NewEtcdManager failed, err:", err)
		return nil, err
	}
	etcdManager = &EtcdManager{
		client: client,
		wg:     &sync.WaitGroup{},
	}
	logrus.Info("NewEtcdManager success, etcdManager:", etcdManager)
	return etcdManager, nil
}

func (e *EtcdManager) GetClient() (etcdClient *clientv3.Client) {
	return e.client
}

func (e *EtcdManager) GetKv() clientv3.KV {
	kv := clientv3.NewKV(e.client)
	return kv
}

func (e *EtcdManager) Put(ctx context.Context, key, value string) (putResp *clientv3.PutResponse, err error) {
	kv := e.GetKv()
	putResp, err = kv.Put(ctx, key, value)
	if err != nil {
		logrus.Error("EtcdManager Put failed, err:", err)
		return nil, err
	}
	return putResp, nil
}

func (e *EtcdManager) Get(ctx context.Context, key string) (getResp *clientv3.GetResponse, err error) {
	kv := e.GetKv()
	getResp, err = kv.Get(ctx, key)
	if err != nil {
		logrus.Error("EtcdManager Get failed, err:", err)
		return nil, err
	}

	return getResp, nil
}

func (e *EtcdManager) Watch(ctx context.Context, key string, logData *common.LogData) {
	e.wg.Add(1)
	watchChan := e.client.Watch(ctx, key)
	logrus.Info("EtcdManager watch key :", key)
	for {
		select {
		case <-ctx.Done():
			e.wg.Add(1)
			e.Close()
			return
		case watchResp := <-watchChan:
			for _, event := range watchResp.Events {
				logrus.Infof("Type: %s Key:%s Value:%s", event.Type, event.Kv.Key, event.Kv.Value)
				if string(event.Kv.Key) == key {
					collectEntryList, err := e.GetConfWithCollectEntry(ctx, key)
					if err != nil {
						logrus.Error("etcd GetConfWithCollectEntry failed, err:", err)
					} else {
						logData.Lock.Lock()
						logData.CollectEntryList = collectEntryList
						logData.Lock.Unlock()
						tail_file.NewTailManager().SendToTailChan()
					}
				}
			}
		}
	}
}

func (e *EtcdManager) GetConfWithCollectEntry(ctx context.Context, key string) (collectEntryList []common.CollectEntry, err error) {
	response, err := e.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get conf from etcd failed,err:%v", err)
		return
	}
	if len(response.Kvs) == 0 {
		logrus.Warningf("get len:0 conf from etcd failed,err:%v", err)
		return
	}
	collectEntryList = make([]common.CollectEntry, 0, 10)
	err = json.Unmarshal(response.Kvs[0].Value, &collectEntryList) //把值反序列化到collectEntryList
	if err != nil {
		logrus.Errorf("json unmarshal failed,err:%v", err)
		return
	}
	return
}

func (e *EtcdManager) Close() {
	e.client.Close()
}

func (e *EtcdManager) WgAdd() {
	e.wg.Add(1)
}

func (e *EtcdManager) WgDone() {
	e.wg.Done()
}

func (e *EtcdManager) WgWait() {
	e.wg.Wait()
}
