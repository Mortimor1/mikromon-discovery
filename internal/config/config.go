package config

import (
	"github.com/Mortimor1/mikromon-discovery/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	Debug bool `yaml:"debug"`
	Http  struct {
		BindIp string `yaml:"bind_ip"`
		Port   string `yaml:"port"`
	} `yaml:"http"`
	Db struct {
		Mongo struct {
			Url      string `yaml:"url"`
			Database string `yaml:"database"`
		} `yaml:"mongo"`
	} `yaml:"db"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("Read Config")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config/config.yml", instance); err != nil {
			desc, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(desc)
			logger.Fatal(err)
		}
	})
	return instance
}
