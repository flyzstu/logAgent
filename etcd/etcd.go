package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 需要收集的日志信息
type LogEntry struct {
	Path  string `json:"path"`  // 日志存放的路径
	Topic string `json:"topic"` // 日志要发往kafka的主题
}

var (
	client *clientv3.Client
)

// 初始化etcd的函数
func Init(addr string, timeout time.Duration) (err error) {

	// 客户端配置
	config := clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: timeout * time.Second,
	}

	// 建立连接
	client, err = clientv3.New(config)
	if err != nil {
		fmt.Printf("init etcd failed, err:%s", err)
		return err
	}
	fmt.Println("connent to etcd succeed")

	return nil
}

// 从etcd读取配置
func GetConf(key string) (logEntryConf []*LogEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := client.Get(ctx, key)
	cancel() // 立即释放资源
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, kv := range resp.Kvs {
		fmt.Printf("key:%s value:%s\n", string(kv.Key), string(kv.Value)) // 索引键值对
		err = json.Unmarshal(kv.Value, &logEntryConf)
		if err != nil {
			fmt.Printf("unmarshal json failed, err:%v\n", err)
			return
		}
	}
	return
}
