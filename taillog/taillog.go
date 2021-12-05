package taillog

import (
	"fmt"
	"logagent/kafka"

	"github.com/hpcloud/tail"
)

// 具体的一个日志收集的任务
type TailTask struct {
	path     string     // 收集路径
	topic    string     // kafka主题
	instance *tail.Tail // 日志收集实例
}

// TailTask构造函数
func NewTailTask(path, topic string) (tailObj *TailTask) {
	tailObj = &TailTask{
		path:  path,
		topic: topic,
	}
	tailObj.init() // 根据路径去打开对应的日志
	return
}

func (t *TailTask) init() {
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开重新创建的文件(失败了会尝试重新读取)
		Follow:    true,                                 // 继续寻找新行 (tail -f)
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从哪个地方开始读取
		MustExist: false,                                // 日志文件是否必须存在
		Poll:      true,                                 // 轮询文件更改而不是使用inotify
	}
	var err error
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Printf("tail file failed, err: %v", err)
	}

	go t.run() // 直接去采集日志信息并发送到kafka
}

func (t *TailTask) run() {
	for {
		select {
		case line := <-t.instance.Lines: // 从tailObj的通道里一行一行读取日志数据
			// 3.2 发往kafka
			// kafka.SendTokafka(t.topic, line.Text) // 函数调用函数：通道
			// 优化：先把日志信息发送到一个通道中
			kafka.SendToChan(t.topic, line.Text)
			// kafaka包中有一个单独的goroutinue去取日志数据信息发送到kafka
		}
	}
}
