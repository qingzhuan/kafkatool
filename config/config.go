package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

type FireEscape struct {
	JpgPath                       string
	PgmPath                       string
	FireBasePath                  string
	CarImageContinueTime          int
	GroundRandomImageContinueTime int
}

type ImageResourcesConfig struct {
	GroundRandomImagePath string
	FireEscape            FireEscape
	ElevatorImagePath     string
}

var Config = &ImageResourcesConfig{}

func InitConfig() {

	viper.SetConfigFile("./config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("读取配置文件出错，err: ", err)
		os.Exit(1)
	}
	err = viper.Unmarshal(Config)
	if err != nil {
		log.Println("解析配置文件出错，err: ", err)
		os.Exit(1)
	}
	log.Printf("init config file success，content: %#v\n", *Config)
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		ReloadConfig()
		log.Printf("监听到文件变更, 变更内容：%#v\n", in)
	})
}

func ReloadConfig() {
	viper.SetConfigFile("./config.yml")
	err := viper.ReadInConfig()
	err = viper.Unmarshal(Config)
	if err != nil {
		log.Printf("监听到文件变更, 重新加载失败，err：%#v", err)
	}
}
