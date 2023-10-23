package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"http-procotol-plugin/global"
	"http-procotol-plugin/utils"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/robfig/cron/v3"
)

//相关业务逻辑
type TpService struct{}

var TpSer = new(TpService)

func (u *TpService) Config() (interface{}, error) {
	data := utils.ReadFormConfig()
	return data, nil
}

func (u *TpService) UpdateDevice(d utils.Device) error {
	global.DevicesMap.Store(d.DeviceConfig.AccessToken, d)
	return nil
}

func (u *TpService) AddDevice(d utils.Device) (err error) {
	global.DevicesMap.Store(d.DeviceConfig.AccessToken, d)
	return nil
}

func (u *TpService) DeleteDevice(d utils.Device) error {
	global.DevicesMap.Range(func(key, value any) bool {
		if value.(utils.Device).DeviceId == d.DeviceId {
			global.DevicesMap.Delete(key)
		}
		return true
	})
	return nil

}

func (u *TpService) Attributes(token string, msg []byte) error {
	//状态发送至tp的mqtt
	err := MqttSendOther(token, "1", global.Conf.Mqtt.StatusTopic)
	if err != nil {
		log.Println("mqtt发送状态失败...", err.Error())
	}
	//数据发送至tp的mqtt
	var deviceId string
	d, _ := global.DevicesMap.Load(token)
	if device, ok := d.(utils.Device); ok {
		device.SetLastMsgTime(utils.GetNowTime())
		device.SetStatus("1")
		global.DevicesMap.Store(token, device)
		deviceId = device.DeviceId
	}

	if global.Conf.Mqtt.StatusTopic == "timescaledb" {
		//状态发送至tp的mqtt
		err = MqttSend(token, msg, global.Conf.Mqtt.AttributesTopic)
	} else {
		//状态发送至tp的mqtt
		err = MqttSend(token, msg, global.Conf.Mqtt.AttributesTopic+"/"+deviceId)
	}

	return err
}

func (u *TpService) Event(token string, msg []byte) error {
	//数据发送至tp的mqtt
	err := MqttSend(token, msg, global.Conf.Mqtt.EventTopic)
	if err != nil {
		log.Println("mqtt发送事件失败...", err.Error())
	}
	//状态发送至tp的mqtt
	err = MqttSendOther(token, "1", global.Conf.Mqtt.StatusTopic)
	d, _ := global.DevicesMap.Load(token)
	if device, ok := d.(utils.Device); ok {
		device.SetLastMsgTime(utils.GetNowTime())
		device.SetStatus("1")
		global.DevicesMap.Store(token, device)
	}
	return err
}

func (u *TpService) CommandReply(token string, msg []byte) error {
	//数据发送至tp的mqtt
	err := MqttSend(token, msg, global.Conf.Mqtt.CommandTopic)
	//状态发送至tp的mqtt
	err = MqttSendOther(token, "1", global.Conf.Mqtt.StatusTopic)
	//修改map里设备的最后一次消息时间
	d, _ := global.DevicesMap.Load(token)
	if device, ok := d.(utils.Device); ok {
		device.SetLastMsgTime(utils.GetNowTime())
		device.SetStatus("1")
		global.DevicesMap.Store(token, device)
	}
	return err
}

//从tp获取设备信息，将token储存在map里
func TpDeviceAccessToken(token string) error {
	requrl := "http://" + global.Conf.Thingspanel.Address + "/api/plugin/device/config"
	request, _ := http.NewRequest("POST", requrl, bytes.NewBuffer([]byte("{\"AccessToken\":\""+token+"\"}")))
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	var res utils.Resdata
	if err := json.Unmarshal(body, &res); err != nil {
		return err
	}
	if res.Code != 200 {
		return errors.New("失败")
	}
	if res.Data.AccessToken == "" {
		return errors.New("失败")
	}
	//if res.Data.ProtocolType != "http-mqtt" {
	//	return errors.New("失败")
	//}
	if res.Data.DeviceConfig.OffineTime == 0 {
		res.Data.DeviceConfig.OffineTime = global.Conf.Thingspanel.OffineTime
	}
	global.DevicesMap.Store(res.Data.AccessToken, res.Data)
	return nil
}

//扫描所有设备状态，将状态发送至tp的mqtt
func OnOfflineCron() {
	crontab := cron.New()
	spec := "0/10 * * * *" //每10秒执行一次
	task := func() {
		global.DevicesMap.Range(func(key, value any) bool {
			device := value.(utils.Device)
			if device.DeviceConfig.Status == "1" && utils.GetNowTime()-device.DeviceConfig.LastMsgTime > device.DeviceConfig.OffineTime {
				log.Println("设备离线:", device.AccessToken)
				device.SetStatus("0")
				global.DevicesMap.Store(key, device)
				//状态发送至tp的mqtt
				err := MqttSendOther(device.AccessToken, "0", global.Conf.Mqtt.StatusTopic)
				if err != nil {
					log.Println("mqtt发送状态失败...", err.Error())
				}
			}
			return true
		})
	}
	crontab.AddFunc(spec, task)
	crontab.Start()
}
