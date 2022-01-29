package main

import (
	"kafkatool/config"
)

func init() {
	config.InitConfig()
}

func main() {
	//go kafkahandler.ForeverWriterCarInfoMsg()
	select {}
}
