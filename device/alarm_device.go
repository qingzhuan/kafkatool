package device

import "sync"

type OnlineDeviceManager struct {
	rwMutex sync.RWMutex
	OnlineDevice map[string]*Device
}

func NewOnlineDeviceManager() *OnlineDeviceManager{
	return &OnlineDeviceManager{
		OnlineDevice: map[string]*Device{},
	}
}

func (odm *OnlineDeviceManager) GetOneOnlineDevice(sceneType string) (d *Device) {
	info := GetDeviceInfo(sceneType)
	listDevice := info.Data.List

	odm.rwMutex.Lock()
	defer odm.rwMutex.Unlock()

	for v,device := range listDevice {
		if _, ok := odm.OnlineDevice[device.DeviceId]; !ok {
			d = &listDevice[v]
		}
	}
	if d != nil {
		odm.OnlineDevice[d.DeviceId] = d
	}
	return
}