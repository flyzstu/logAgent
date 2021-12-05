package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

type logData struct {
	topic string
	data  string
}

// 专门往kafka写日志的模块
var (
	client      sarama.SyncProducer // 声明一个全局的连接kafka的生产者client
	logDataChan chan *logData
)

// Init 初始化client
func Init(addrs []string, maxSize int) (err error) {
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
	// 初始化全局通道
	logDataChan = make(chan *logData, maxSize)
	// 开启后台的goroutine从通道中取数据发往kafka
	go sendTokafka()
	return
}

// 真正往kafka发送日志的接口
func sendTokafka() {
	for {
		select {
		case ld := <-logDataChan:
			{
				// 构造一个消息
				msg := &sarama.ProducerMessage{}
				msg.Topic = ld.topic
				msg.Value = sarama.StringEncoder(ld.data)
				// 发送到kafka
				pid, offset, err := client.SendMessage(msg)
				if err != nil {
					fmt.Println("send msg failed, err:", err)
					return
				}
				fmt.Printf("pid: %v, offset:%v\n", pid, offset)
			}
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}

// 给外部暴露的一个函数，该函数只把日志数据只把日志数据发生到一个内部的channel中
func SendToChan(topic, data string) {
	msg := &logData{
		topic: topic,
		data:  data,
	}
	logDataChan <- msg // 实现异步发送

}
