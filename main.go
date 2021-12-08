package main

import (
	"fmt"
	"logagent/conf"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/taillog"
	"sync"

	"gopkg.in/ini.v1"
)

var cfg = new(conf.AppConf) // 要为cfg创建一个地址

// logAgent入口程序
func main() {
	// 0. 加载配置文件
	err := ini.MapTo(cfg, "./conf/config.ini")
	if err != nil {
		fmt.Println("load conf.ini failed, err: ", err)
		return
	}
	// 1.初始化kafka链接
	err = kafka.Init([]string{cfg.KafkaConf.Address}, cfg.ChanMaxSize)
	if err != nil {
		fmt.Printf("init kafka failed, err :%v\n", err)
		return
	}
	fmt.Println("init kafka success.")
	// 2. 初始化etcd
	// fmt.Printf("cfg: %#v\n", cfg)
	err = etcd.Init(cfg.EtcdConf.Address, cfg.EtcdConf.Timeout)
	if err != nil {
		fmt.Printf("init etcd failed, err :%v\n", err)
		return
	}
	fmt.Println("init etcd client success.")
	// 2.1 从etcd中获取日志收集项的配置信息
	logEntryConf, err := etcd.GetConf(cfg.EtcdConf.Key)
	if err != nil {
		fmt.Printf("get conf failed, err:%v\n", err)
		return
	}
	fmt.Printf("get conf from etcd success, %v\n", logEntryConf)
	// 2.2 派一个哨兵去监视日志收集项的变化

	for index, value := range logEntryConf {
		fmt.Printf("index:%v, err:%v\n", index, *value)
	}
	// 3. 收集日志发往kafka
	taillog.Init(logEntryConf)
	var wg sync.WaitGroup
	wg.Add(1)
	newConfChan := taillog.NewConfChan()             // 从taillog包中获取对外暴露的通道
	go etcd.WatchConf(cfg.EtcdConf.Key, newConfChan) // 哨兵发现最新的配置信息会通知上面的通道
	wg.Wait()
}
