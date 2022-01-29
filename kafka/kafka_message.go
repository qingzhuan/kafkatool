package kafkahandler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"kafkatool/common"
	"kafkatool/constant"
	"kafkatool/device"
	images "kafkatool/image"
	"strconv"
	"time"
)

var deviceManager = device.NewOnlineDeviceManager()

type Message struct {
	CommunityId string                   `json:"community_id"`
	DeviceId    string                   `json:"device_id"`
	FrameHeight int                      `json:"frame_height"`
	FrameId     int                      `json:"frame_id"`
	FrameUri    string                   `json:"frame_uri"`
	FrameWidth  int                      `json:"frame_width"`
	JpgPath     string                   `json:"jpg_path"`
	ModuleId    string                   `json:"module_id"`
	MsgId       string                   `json:"msg_id"`
	Objects     []map[string]interface{} `json:"objects"`
	Payload     map[string]interface{}   `json:"payload"`
	Pts         int64                    `json:"pts"`
	Timestamps  int64                    `json:"timestamps"`
}

func BuildKafkaMessage(communityId, deviceId string) (msg Message){
	msg = Message{
		CommunityId: communityId,
		DeviceId:    deviceId,
		FrameId:     40663,
		ModuleId:    "decoder0",
		Objects:     []map[string]interface{}{},
		Payload:     map[string]interface{}{},
		Pts:         0,
	}
	detectFile := <-images.ImageQueue
	pgm := detectFile.PGM
	width, height := images.GetImageWidthAndHeight(detectFile.JPG)
	msg.JpgPath = detectFile.JPG
	msg.FrameUri = pgm
	msg.FrameHeight = height
	msg.FrameWidth = width
	msg.MsgId = common.MD5Value(strconv.Itoa(time.Now().Nanosecond()))
	msg.Timestamps = time.Now().UnixNano() / 1e6
	return
}

func SendKafkaMessage(w *kafka.Writer, msg Message) {
	bytes, err2 := json.Marshal(msg)
	fmt.Println(string(bytes))
	if err2 != nil {
		return
	}
	message := kafka.Message{}
	message.Value = bytes
	err := w.WriteMessages(context.Background(), message)
	fmt.Printf("err:%#v \n", err)
}

func ForeverWriterCarInfoMsg() {
	w := KafkaHandler.TopicWriterHandler[constant.InfoDecodeTopic]
	onlineDevice := deviceManager.GetOneOnlineDevice(constant.GroundSceneType)
	onlineDevice = deviceManager.GetOneOnlineDevice(constant.GroundSceneType)
	fmt.Println("onlineDevice", deviceManager.OnlineDevice)
	for {
		msg := BuildKafkaMessage(onlineDevice.CommunityId, onlineDevice.DeviceId)
		SendKafkaMessage(w, msg)
		time.Sleep(time.Second * 1)
	}
}
