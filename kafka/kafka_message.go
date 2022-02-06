package kafkahandler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"kafkatool/common"
	"kafkatool/config"
	"kafkatool/constant"
	"kafkatool/device"
	images "kafkatool/image"
	"log"
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

func BuildKafkaMessage(communityId, deviceId string, imageQueue <-chan images.DetectFile) (msg Message){
	msg = Message{
		CommunityId: communityId,
		DeviceId:    deviceId,
		FrameId:     40663,
		ModuleId:    "decoder0",
		Objects:     []map[string]interface{}{},
		Payload:     map[string]interface{}{},
		Pts:         0,
	}
	//detectFile := <-images.FireEscapeImageQueue
	detectFile := <- imageQueue
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
	log.Println(string(bytes))
	if err2 != nil {
		return
	}
	message := kafka.Message{}
	message.Value = bytes
	err := w.WriteMessages(context.Background(), message)
	if err != nil {
		log.Printf("send kafka message err:%#v \n", err)
	}
}

func ForeverWriterCarInfoMsg() {
	// 生产消防通道需要的图片jpg+pgm
	go images.ProduceFireEscapeJpgImage()
	// 生产地上场景的图片（从某一个目录随机获取一张图片）
	parameter := &images.ProduceImageParameter{
		ImageQueue: images.GroundImageQueue,
		Path:       &config.Config.GroundRandomImagePath,
	}

	go images.ProduceJpgImage(parameter)
	// 定时删除生成的jpg
	go images.ClearImage(config.Config.FireEscape.FireBasePath + "/jpg")

	w := KafkaHandler.TopicWriterHandler[constant.InfoDecodeTopic]
	onlineDevice := deviceManager.GetOneOnlineDevice(constant.GroundSceneType)
	//onlineDevice = deviceManager.GetOneOnlineDevice(constant.GroundSceneType)
	fmt.Println("onlineDevice", deviceManager.OnlineDevice)
	for {
		msg := BuildKafkaMessage(onlineDevice.CommunityId, onlineDevice.DeviceId, images.FireEscapeImageQueue)
		SendKafkaMessage(w, msg)
		time.Sleep(time.Second * 1)
	}
}

func ForeverWriterElevatorMsg()  {
	// 生产电梯场景可以报警的图片
	produceElevatorAlarmImagePara := &images.ProduceImageParameter{
		ImageQueue: images.ElevatorAlarmImageQueue,
		Path:       &config.Config.ElevatorImageConfig.ElevatorAlarmImagePath,
	}
	go images.ProduceJpgImage(produceElevatorAlarmImagePara)

	// 随机生产电梯场景图片
	produceElevatorRandomImagePara := &images.ProduceImageParameter{
		ImageQueue: images.ElevatorRandomImageQueue,
		Path:       &config.Config.ElevatorImageConfig.ElevatorRandomImagePath,
	}
	go images.ProduceJpgImage(produceElevatorRandomImagePara)

	// 生产总的报警图片
	go images.ProduceElevatorSceneImage()

	// 获取设备
	onlineDevice := deviceManager.GetOneOnlineDevice(constant.ElevatorSceneType)
	if onlineDevice == nil {
		log.Println("not found elevator type device")
		return
	}

	w := KafkaHandler.TopicWriterHandler[constant.ElevatorDecodeTopic]
	for {
		msg := BuildKafkaMessage(onlineDevice.CommunityId, onlineDevice.DeviceId, images.ElevatorImageQueue)
		SendKafkaMessage(w, msg)
		time.Sleep(time.Second * 1)
	}

	
}