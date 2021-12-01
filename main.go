package main

import (
	"fmt"
	"logAgent/kafka"
	"logAgent/taillog"
	"time"
)

func run() {
	// 1. 读取日志
	for {
		select {
		case line := <-taillog.ReadChan(): // 读取到日志
			//2. 发送到kafka
			kafka.SendTokafka("web_log", line.Text)
		default: // 读取不到日志
			time.Sleep(time.Second)
		}
	}
}

// logAgent入口程序
func main() {
	// 1.初始化kafka链接
	err := kafka.Init([]string{"192.168.31.103:9092"})
	if err != nil {
		fmt.Printf("init kafka failed, err :%s\n", err.Error())
		return
	}
	fmt.Println("init kafka success.")
	// 2.打开日志文件准备收集日志
	err = taillog.Init("./my.log")
	if err != nil {
		fmt.Printf("init taillog failed, err: %s\n", err.Error())
		return
	}
	fmt.Println("init taillog success.")
	run()
}
