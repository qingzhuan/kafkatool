package device

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kafkatool/constant"
	"log"
	"net/http"
	"os"
	"strings"
)

type Device struct {
	DeviceId    string `json:"device_id"`
	Status      int    `json:"status"`
	CommunityId string `json:"community_id"`
	SceneType   string `json:"scene_type"`
}

type DeviceData struct {
	List  []Device `json:"list"`
	Count int64    `json:"count"`
}

type DeviceResp struct {
	Data DeviceData `json:"data"`
}

func GetDeviceInfo(sceneType string) (deviceResp DeviceResp) {
	c := http.Client{}
	reader := strings.NewReader(fmt.Sprintf(`{"status": 1, "scene_type": "%s"}`, sceneType))
	request, err := http.NewRequest("POST", constant.CommunityHost + "/inner/api/device/list", reader)
	if err != nil {
		log.Printf("获取设备列表设备，请检查错误 err: %#v\n", err)
		os.Exit(1)
		return
	}
	defer request.Body.Close()
	response, err := c.Do(request)
	if err != nil {
		log.Printf("获取设备列表设备，请检查错误 err: %#v\n", err)
		os.Exit(1)
		return
	}
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %#v\n", err)
		return
	}
	//fmt.Println(string(all))
	err = json.Unmarshal(all, &deviceResp)
	if err != nil {
		log.Println("json.Unmarshal", err)
	}
	return
}