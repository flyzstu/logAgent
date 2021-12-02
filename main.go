package main

import (
	"fmt"
	"logAgent/conf"
	"logAgent/kafka"
	"logAgent/taillog"
	"time"

	"gopkg.in/ini.v1"
)

// var cfg *conf.AppConf       // 默认值是nil，会报错
var cfg = new(conf.AppConf) // 要为cfg创建一个地址

func run() {

	// 1. 读取日志
	for {
		select {
		case line := <-taillog.ReadChan(): // 读取到日志
			//2. 发送到kafka
			kafka.SendTokafka(cfg.KafkaConf.Topic, line.Text)
		default: // 读取不到日志
			time.Sleep(time.Second)
		}
	}
}

// logAgent入口程序
func main() {
	// 0. 加载配置文件
	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Println("load conf.ini failed, err: ", err)
		return
	}
	// 1.初始化kafka链接
	err = kafka.Init([]string{cfg.KafkaConf.Address})
	if err != nil {
		fmt.Printf("init kafka failed, err :%s\n", err.Error())
		return
	}
	fmt.Println("init kafka success.")
	// 2.打开日志文件准备收集日志
	err = taillog.Init(cfg.TaillogConf.FileName)
	if err != nil {
		fmt.Printf("init taillog failed, err: %s\n", err.Error())
		return
	}
	fmt.Println("init taillog success.")
	run()
}
