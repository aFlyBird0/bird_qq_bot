package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

// GlobalConfig 默认全局配置
var GlobalConfig *Config

var watchConfig chan struct{}

// init 使用 ./application.yaml 初始化全局配置
func init() {
	GlobalConfig = &Config{
		viper.New(),
	}
	GlobalConfig.SetConfigName("application")
	GlobalConfig.SetConfigType("yaml")
	GlobalConfig.AddConfigPath(".")
	GlobalConfig.AddConfigPath("./config")

	err := GlobalConfig.ReadInConfig()
	if err != nil {
		logrus.WithField("config", "globalconfig").WithError(err).Panicf("unable to read global config")
	}

	watchConfig = make(chan struct{}, 2)

	GlobalConfig.WatchConfig() //实时监听配置变化
	GlobalConfig.OnConfigChange(func(e fsnotify.Event) {
		err := GlobalConfig.ReadInConfig()
		if err == nil {
			// 配置文件发生变更之后会调用的回调函数
			logrus.WithField("config", "globalconfig").Info("Config file changed:", e.Name, e.Op)
			watchConfig <- struct{}{}
		} else {
			logrus.WithField("config", "globalconfig").WithError(err).Error("Config file changed, but unable to read in new config")
		}
	})
}

func GetWatcher() chan struct{} {
	return watchConfig
}
