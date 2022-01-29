package kafkahandler

import (
	"github.com/segmentio/kafka-go"
	"kafkatool/constant"
	"sync"
)

var KafkaHandler = NewKafkaHandler()

type KafkaHandlerManager struct {
	rwMutex sync.RWMutex
	TopicWriterHandler map[string]*kafka.Writer
}

func NewKafkaHandler() (kafkaHandler *KafkaHandlerManager){
	kafkaHandler = &KafkaHandlerManager{
		TopicWriterHandler: make(map[string]*kafka.Writer),
	}
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	kafkaHandler.TopicWriterHandler[constant.InfoDecodeTopic] = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{constant.KafKaHost},
		Topic:    constant.InfoDecodeTopic,
		Balancer: &kafka.LeastBytes{},
	})

	kafkaHandler.TopicWriterHandler[constant.ElevatorDecodeTopic] = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{constant.KafKaHost},
		Topic:    constant.ElevatorDecodeTopic,
		Balancer: &kafka.LeastBytes{},
	})
	return
}