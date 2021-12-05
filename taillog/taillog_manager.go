package taillog

import (
	"logagent/etcd"
)

type taillogMgr struct {
	logEntry []*etcd.LogEntry
	// tskMap   map[string]TailTask  //配置热更新
}

var tskMgr *taillogMgr // 管理者

func Init(logEntryConf []*etcd.LogEntry) {
	tskMgr = &taillogMgr{
		logEntry: logEntryConf, // 把当前日志收集项的配置存起来
	}
	for _, logEntry := range logEntryConf {
		// conf: *etcd.LogEntry
		// logEntry.path: 要收集日志的路径
		NewTailTask(logEntry.Path, logEntry.Topic)
	}
}
