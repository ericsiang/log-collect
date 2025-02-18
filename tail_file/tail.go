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

var (
	once        sync.Once
	tailManager *TailManager
)

type TailManager struct {
	tailChan     chan struct{}
	tailTaskList []*tailTask
}

func NewTailManager() *TailManager {
	once.Do(func() {
		tailManager = &TailManager{
			tailChan:     make(chan struct{}, 0),
			tailTaskList: make([]*tailTask, 0, 10),
		}
	})
	return tailManager
}

func (t *TailManager) SendToTailChan() {
	logrus.Info("TailManager SendToTailChan")
	t.tailChan <- struct{}{}
}

func (t *TailManager) ReloadInitTailTask(ctx context.Context, wg *sync.WaitGroup, kafkaProducerManager *kafka.KafkaProducerManager, logData *common.LogData) (err error) {
	for {
		select {
		case <-ctx.Done():
			logrus.Info("TailManager ReloadInitTailTask ctx.Done()")
			wg.Done()
			return
		case <-t.tailChan:
			err = t.InitTailTask(ctx, wg, kafkaProducerManager, logData.CollectEntryList)
			if err != nil {
				logrus.Error("InitTail failed, err:", err)
				return
			}
		}
	}
}

type tailTask struct {
	path       string
	topic      string
	TailObj    *tail.Tail
	taskCancel context.CancelFunc
}

func NewTailTask(taskCancel context.CancelFunc, path, topic string) (tt *tailTask) {
	tt = &tailTask{
		path:       path,
		topic:      topic,
		taskCancel: taskCancel,
	}
	return tt
}

func (t *tailTask) InitTail() (err error) {
	/*
		offset為相對偏移量，而whence決定相對位置：
		0為相對檔案開頭，1為相對目前位置，2為相對檔案結尾。它會傳回新的偏移量（相對開頭）和可能的錯誤。

		whence參數
		io.SeekStart // 0 ， io.SeekCurrent // 1 ， io.SeekEnd // 2
	*/
	seek := &tail.SeekInfo{Offset: 0, Whence: 2}
	config := tail.Config{
		Follow:    true, //true:一直監聽(同tail -f) false:一次後即結束
		ReOpen:    true, //true:文件被删,阻塞并等待新建此文件(同tail -F) false:文件被删,程序结束
		MustExist: true, //true:文件不存在即退出 false:文件不存在即等待
		Poll:      true,
		Location:  seek,
	}
	fileName := t.path + "/" + t.topic
	logrus.Info("tail fileName:", fileName)
	t.TailObj, err = tail.TailFile(fileName, config)
	if err != nil {
		logrus.Error("tail create tailObj for path:", t.path, " , err:", err)
		return err
	}
	logrus.Info("tail Init success")
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

func (t *TailManager) InitTailTask(ctx context.Context, wg *sync.WaitGroup, kafka *kafka.KafkaProducerManager, collectEntryList []common.CollectEntry) (err error) {
	for _, tailTask := range t.tailTaskList {
		tailTask.taskCancel() //將之前的 run() goroutine 終止
	}
	t.tailTaskList = make([]*tailTask, 0, 10) //清舊資料
	for _, entry := range collectEntryList {
		taskCtx, taskCancel := context.WithCancel(context.Background())
		tt := NewTailTask(taskCancel, entry.Path, entry.Topic)
		err = tt.InitTail()
		if err != nil {
			logrus.Errorf("tail create tailObj for path:%s,err:%v", entry.Path, err)
			return
		}
		t.tailTaskList = append(t.tailTaskList, tt)
		wg.Add(1)
		go tt.run(taskCtx, wg, kafka)
	}
	return nil
}
