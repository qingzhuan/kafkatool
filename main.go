package main

import (
	"kafkatool/config"
	images "kafkatool/image"
	kafkahandler "kafkatool/kafka"
)

func init() {
	config.InitConfig()
}

func main() {
	go images.ProduceJpgImage()
	kafkahandler.ForeverWriterCarInfoMsg()
	select {}
}