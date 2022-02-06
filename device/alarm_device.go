package device

import (
	"fmt"
	"sync"
	"time"
)

type OnlineDeviceManager struct {
	rwMutex      sync.RWMutex
	OnlineDevice map[string]*Device
}

type RunningDevice struct {
}

func NewOnlineDeviceManager() *OnlineDeviceManager {
	odm := &OnlineDeviceManager{
		OnlineDevice: map[string]*Device{},
	}
	odm.TimerCheckDevice()
	return odm
}

func (odm *OnlineDeviceManager) GetOneOnlineDevice(sceneType string) (d *Device) {
	info := GetDeviceInfo(sceneType)
	listDevice := info.Data.List

	odm.rwMutex.Lock()
	defer odm.rwMutex.Unlock()

	for v, device := range listDevice {
		if _, ok := odm.OnlineDevice[device.DeviceId]; !ok {
			d = &listDevice[v]
		}
	}
	if d != nil {
		odm.OnlineDevice[d.DeviceId] = d
	}
	return
}

/* 定时检查已经加载的设备是否已被修改为别的场景 */
func (odm *OnlineDeviceManager) TimerCheckDevice() {
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		defer ticker.Stop()
		for {
			if len(odm.OnlineDevice) < 1 {
				continue
			}
			fmt.Printf("odm.OnlineDevice: %#v\n", odm.OnlineDevice)
			for d := range odm.OnlineDevice {
				isChange := odm.checkOnlineDeviceIsChange(odm.OnlineDevice[d])
				fmt.Println("isChange: ", isChange)
				// 如果别更改需要重新获取设备，并添加到OnlineMap中
				if isChange {
					odm.GetOneOnlineDevice(odm.OnlineDevice[d].SceneType)
					// 删除之前到设备
					delete(odm.OnlineDevice, d)
				}

			}
			<-ticker.C
		}

	}()
}

/* 检查添加到OnlineMap中的设备场景是否被更改 */
func (odm *OnlineDeviceManager) checkOnlineDeviceIsChange(device *Device) (isChange bool) {
	info := GetDeviceInfo(device.SceneType)
	listDevice := info.Data.List
	for _, v := range listDevice {
		if v.DeviceId == device.DeviceId {
			isChange = false
			return
		}
	}
	isChange = true
	return
}
