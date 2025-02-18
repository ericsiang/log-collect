package common

import "sync"

type Config struct {
	Etcdaddress  ConfigEtcd        `ini:"etcd"`
	Kafakaddress ConfigKafkaddress `ini:"kafka"`
	LogFilePath  ConfigLogFilePath `ini:"collect"`
	Etcd1        ConfigEtcd        `ini:"etcd1"`
	Etcd2        ConfigEtcd        `ini:"etcd2"`
}

type ConfigKafkaddress struct {
	Address     string `ini:"address"`
	Topic       string `ini:"topic"`
	MessageSize int64  `ini:"chan_size"`
}

type ConfigLogFilePath struct {
	Path string `ini:"logfile_path"`
}

type ConfigEtcd struct {
	Address []string `ini:"address"`
	Key     string   `ini:"collect_key"`
}

type CollectEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

type LogData struct {
	Lock             sync.Mutex
	CollectEntryList []CollectEntry
}

type ElasticConfig struct {
	ApiKey     string   `json:"api_key"`
	UrlAddress []string `json:"url_address"`
}
