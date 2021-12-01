package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

// 专门往kafka写日志的模块
var (
	client sarama.SyncProducer // 声明一个全局的链接kafka的生产者client
)

// Init 初始化client
func Init(addrs []string) (err error) {
	config := sarama.NewConfig()
	// tailf包使用
	config.Producer.RequiredAcks = sarama.WaitForAll          //发送完数据需要等待follow回复ack
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 先选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel中返回
	// 链接kafka
	client, err = sarama.NewSyncProducer(addrs, config)

	if err != nil {
		fmt.Println("producer closed, err :", err)
		return
	}
	return
}
func SendTokafka(topic, data string) {
	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)
	// 发送到kafka
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return
	}
	fmt.Printf("pid: %v, offset:%v\n", pid, offset)
}
