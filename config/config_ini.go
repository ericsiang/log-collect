package config

import (
	"fmt"
	"kafka-log/common"

	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

func NewConfig() (configObj *common.Config, err error) {
	configObj = new(common.Config) //生成指针便于参数传递
	err = ini.MapTo(configObj, "../config.ini")
	if err != nil {
		wrapError := fmt.Errorf("ini MapTo() failed , err :%w", err)
		logrus.Error("log config failed,err:", wrapError)
		return nil, wrapError
	}
	return configObj, nil
}
