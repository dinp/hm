package g

import (
	"encoding/json"
	"github.com/toolkits/file"
	"log"
	"sync"
)

const (
	VERSION = "1.0.0"
)

type DBConfig struct {
	Dsn     string `json:"dsn"`
	MaxIdle int    `json:"maxIdle"`
}

type GlobalConfig struct {
	Debug           bool      `json:"debug"`
	CheckInterval   int       `json:"check_interval"`
	DockerPort      int       `json:"dockerPort"`
	ResponseTimeout int       `json:"response_timeout"`
	HealthSign      string    `json:"health_sign"`
	ServerHttpApi   string    `json:"server_http_api"`
	DB              *DBConfig `json:"db"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()

	config = &c

	if config.Debug {
		log.Println("read config file:", cfg, "successfully")
	}
}
