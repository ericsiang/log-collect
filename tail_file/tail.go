package tail_file

import (
	"context"
	"kafka-log/common"
	"kafka-log/kafka"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"

	"github.com/hpcloud/tail"
)

type tailTask struct {
	path    string
	topic   string
	TailObj *tail.Tail
}

func NewTailTask(path, topic string) (tt *tailTask) {
	tt = &tailTask{
		path:  path,
		topic: topic,
	}
	return tt
}

func (t *tailTask) Init() (err error) {
	config := tail.Config{
		Follow:    true,
		ReOpen:    true,
		MustExist: true,
		Poll:      true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
	}
	fileName := t.path + "/" + t.topic
	logrus.Info("tail fileName:", fileName)
	t.TailObj, err = tail.TailFile(fileName, config)
	if err != nil {
		logrus.Error("tail create tailObj for path:", t.path, " , err:", err)
		return err
	}
	return nil
}

func (t *tailTask) run(ctx context.Context, wg *sync.WaitGroup, kafka *kafka.KafkaProducerManager) (err error) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			t.TailObj.Cleanup()
			return
		case line, ok := <-t.TailObj.Lines:
			if !ok {
				logrus.Warnf("tailFile.TailObj.Lines channel closed,path:%s\n", t.path)
				time.Sleep(2 * time.Second)
				continue
			}
			if len(strings.Trim(line.Text, "\r")) == 0 {
				continue
			}
			logrus.Info("get new line:", line.Text)
			msg := &sarama.ProducerMessage{}
			msg.Topic = t.topic
			msg.Value = sarama.StringEncoder(line.Text)
			kafka.SendToMsgChan(msg) //發送到 channel
		}
	}
}

func InitTail(ctx context.Context, wg *sync.WaitGroup, kafka *kafka.KafkaProducerManager, collectEntryList []common.CollectEntry) (err error) {
	for _, entry := range collectEntryList {
		tt := NewTailTask(entry.Path, entry.Topic)
		err = tt.Init()
		if err != nil {
			logrus.Errorf("tail create tailObj for path:%s,err:%v", entry.Path, err)
			return
		}
		wg.Add(1)
		go tt.run(ctx, wg, kafka)
	}
	return nil
}
