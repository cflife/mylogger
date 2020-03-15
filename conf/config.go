package conf

type AppConf struct {
	KafkaConf
	EtcdConf
}

type KafkaConf struct {
	Address     string `ini:"address"`
	ChanMaxSize int    `ini:"chan_max_size"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	Key     string `ini:"collect_log_key"`
	TimeOut int    `ini:"timeOut"`
}
