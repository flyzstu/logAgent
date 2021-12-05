package conf

import "time"

type AppConf struct {
	KafkaConf   `ini:"kafka"`
	TaillogConf `ini:"taillog"`
	EtcdConf    `ini:"etcd"`
}
type EtcdConf struct {
	Address string        `ini:"address"`
	Timeout time.Duration `ini:"timeout"`
	Key     string        `ini:"collect_log_key"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	// Topic   string `ini:"topic"`
	ChanMaxSize int `ini:"chan_max_size"`
}

//-----unused ↓--------
type TaillogConf struct {
	FileName string `ini:"path"`
}

// 结构体和分区双向映射

// func main() {
// 	cfg := new(AppConf)
// 	err := ini.MapTo(cfg, "./config.ini")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("cfg: %v", cfg)
// }
