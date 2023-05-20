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

	"github.com/robfig/cron"
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
	//数据发送至tp的mqtt
	err := MqttSend(token, msg, global.Conf.Mqtt.AttributesTopic)
	status := []byte("{\"status\":\"1\"}")
	//状态发送至tp的mqtt
	err = MqttSend(token, status, global.Conf.Mqtt.StatusTopic)
	d, _ := global.DevicesMap.Load(token)
	if device, ok := d.(utils.Device); ok {
		device.SetLastMsgTime(utils.GetNowTime())
		global.DevicesMap.Store(token, device)
	}
	return err
}

func (u *TpService) Event(token string, msg []byte) error {
	//数据发送至tp的mqtt
	err := MqttSend(token, msg, global.Conf.Mqtt.EventTopic)
	status := []byte("{\"status\":\"1\"}")
	//状态发送至tp的mqtt
	err = MqttSend(token, status, global.Conf.Mqtt.StatusTopic)
	d, _ := global.DevicesMap.Load(token)
	if device, ok := d.(utils.Device); ok {
		device.SetLastMsgTime(utils.GetNowTime())
		global.DevicesMap.Store(token, device)
	}
	return err
}

func (u *TpService) CommandReply(token string, msg []byte) error {
	//数据发送至tp的mqtt
	err := MqttSend(token, msg, global.Conf.Mqtt.CommandTopic)
	status := []byte("{\"status\":\"1\"}")
	//状态发送至tp的mqtt
	err = MqttSend(token, status, global.Conf.Mqtt.StatusTopic)
	//修改map里设备的最后一次消息时间
	d, _ := global.DevicesMap.Load(token)
	if device, ok := d.(utils.Device); ok {
		device.SetLastMsgTime(utils.GetNowTime())
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
		return errors.New(res.Msg)
	}
	global.DevicesMap.Store(res.Data.AccessToken, res.Data)
	return nil
}

//扫描所有设备状态，将状态发送至tp的mqtt
func OnOfflineCron() {
	crontab := cron.New()
	spec := "0/1 * * * *" //每秒执行一次
	task := func() {
		log.Println("扫描设备状态...")
		global.DevicesMap.Range(func(key, value any) bool {
			device := value.(utils.Device)
			if utils.GetNowTime()-device.DeviceConfig.LastMsgTime > device.DeviceConfig.OffineTime {
				log.Println("设备离线:", device.DeviceConfig.AccessToken)
				status := []byte("{\"status\":\"0\"}")
				//状态发送至tp的mqtt
				err := MqttSend(device.DeviceConfig.AccessToken, status, global.Conf.Mqtt.StatusTopic)
				if err != nil {
					log.Println("mqtt发送状态失败...", err.Error())
				}
			}
			return true
		})
		log.Println("扫描设备状态完成...")
	}
	crontab.AddFunc(spec, task)
	crontab.Start()
}
