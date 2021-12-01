package taillog

import (
	"fmt"

	"github.com/hpcloud/tail"
)

// 跟踪日志的模块
var (
	tailObj *tail.Tail // 全局变量
)

func Init(fileName string) (err error) {
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开重新创建的文件(失败了会尝试重新读取)
		Follow:    true,                                 // 继续寻找新行 (tail -f)
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从哪个地方开始读取
		MustExist: false,                                // 日志文件是否必须存在
		Poll:      true,                                 // 轮询文件更改而不是使用inotify
	}
	tailObj, err = tail.TailFile(fileName, config)
	if err != nil {
		fmt.Printf("tail file failed, err: %v", err)
		return
	}
	return
}
func ReadChan() <-chan *tail.Line {
	return tailObj.Lines
}
