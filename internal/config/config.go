package config

import (
	"github.com/Mortimor1/mikromon-discovery/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
	"time"
)

type Config struct {
	Debug bool `yaml:"debug"`
	Http  struct {
		BindIp string `yaml:"bind_ip"`
		Port   string `yaml:"port"`
	} `yaml:"http"`
	Service struct {
		Core []string `yaml:"core"`
	} `yaml:"service"`
	Discovery struct {
		Interval time.Duration `yaml:"interval"`
	} `yaml:"discovery"`
}

var instance *Config
var once sync.Once

var configPaths = []string{"config/discovery.yml", "discovery.yml", "/etc/mikromon/discovery.yml"}

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("Read Config")
		instance = &Config{}

		var errorRead error
		for _, path := range configPaths {
			errorRead = cleanenv.ReadConfig(path, instance)
			if errorRead == nil {
				logger.Infof("Config loaded success path: %s", path)
				break
			}
		}

		if errorRead != nil {
			logger.Fatal(errorRead)
		}
	})
	return instance
}
