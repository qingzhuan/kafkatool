package main

import (
	"kafkatool/config"
	kafkahandler "kafkatool/kafka"
)

func init() {
	config.InitConfig()
}

func main() {
	go kafkahandler.ForeverWriterCarInfoMsg()
	go kafkahandler.ForeverWriterElevatorMsg()
	select {}
}
