package conf

type AppConf struct {
	KafkaConf   `ini:"kafka"`
	TaillogConf `ini:"taillog"`
}

type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}
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
