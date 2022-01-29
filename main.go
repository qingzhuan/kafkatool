package main

import (
	"kafkatool/config"
	kafkahandler "kafkatool/kafka"
)

func init() {
	config.InitConfig()
}

func main() {
	kafkahandler.ForeverWriterCarInfoMsg()
	select {}
}
