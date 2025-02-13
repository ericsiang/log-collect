package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

func main() {
	filename := "msg.log"
	config := tail.Config{
		Follow:    true,  //进行跟随
		ReOpen:    true,  //重新打开
		MustExist: false, //文件打开失败不报错
		Poll:      true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
	}
	tail, err := tail.TailFile(filename, config)
	if err != nil {
		fmt.Println("tail file failed,err:", err)
		return
	}
	for {
		line, ok := <-tail.Lines
		if !ok {
			fmt.Println("tail file close reopen, filename: ", tail.Filename)
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Println("line:", line.Text)
	}
}
