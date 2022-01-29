package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type FireEscape struct {
	JpgPath string
	PgmPath string
}

type ImageResourcesConfig struct {
	GroundRandomImagePath string
	FireEscape FireEscape
	ElevatorImagePath string
}

var Config = ImageResourcesConfig{}

func InitConfig(){
	viper.SetConfigFile("./config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("读取配置文件出错，err: ", err)
		os.Exit(1)
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		log.Println("解析配置文件出错，err: ", err)
		os.Exit(1)
	}
}

