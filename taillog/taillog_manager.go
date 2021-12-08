package taillog

import (
	"fmt"
	"logagent/etcd"
	"time"
)

// tailTask管理者
type taillogMgr struct {
	logEntry    []*etcd.LogEntry
	tskMap      map[string]*TailTask //配置热更新
	newConfChan chan []*etcd.LogEntry
}

var tskMgr *taillogMgr

func Init(logEntryConf []*etcd.LogEntry) {
	tskMgr = &taillogMgr{
		logEntry:    logEntryConf,                   // 把当前日志收集项的配置存起来
		tskMap:      make(map[string]*TailTask, 16), // 分配内存
		newConfChan: make(chan []*etcd.LogEntry),    // 无缓冲区的通道:没人收就卡住了
	}
	for _, logEntry := range logEntryConf {
		// conf: *etcd.LogEntry
		// logEntry.path: 要收集日志的路径
		// 初始化的时候启动了多少tailtask需要记录下来，为了后续判断方便
		tailObj := NewTailTask(logEntry.Path, logEntry.Topic)
		// logEntry.Path：用于标识同一台机器上不同的日志收集项
		// logEntry.Path 和 logEntry.Topic 两个字符串连接起来，用于判断日志的收集项的Path和Topic任一出现修改的情况
		mkey := fmt.Sprintf("%s_%s", logEntry.Path, logEntry.Topic)
		tskMgr.tskMap[mkey] = tailObj
	}
	go tskMgr.run()
}

// 监听自己的newConfChan，有了新的配置过来之后就做处理

func (t *taillogMgr) run() {
	for {
		select {
		case newConf := <-t.newConfChan:
			// 配置是一个列表，需要遍历
			for _, conf := range newConf {
				// 判断conf.Path在不在启动的map(tskMgr.tskMap[logEntry.Path])里面
				mkey := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
				_, ok := t.tskMap[mkey]
				if ok {
					// 原来就有，不需要操作
					continue
				} else {
					// 原来没有的，当成一个新的来处理
					tailObj := NewTailTask(conf.Path, conf.Topic)
					mkey := fmt.Sprintf("%s_%s", conf.Path, conf.Topic)
					t.tskMap[mkey] = tailObj
				}

			}
			// 2. 找出原来t.tskMap有但是newConf里面没有的，删去t.tskMap里面的任务
c
			// 2.配置删除
			// 3.配置变更
			fmt.Println("新的配置来了: ", newConf)
		default:
			time.Sleep(time.Second)
		}
	}
}

//把tskMgr的newConfChan向外暴露
func NewConfChan() chan<- []*etcd.LogEntry {
	return tskMgr.newConfChan
}
