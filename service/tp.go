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
	return err
}

func (u *TpService) Event(token string, msg []byte) error {
	//数据发送至tp的mqtt
	err := MqttSend(token, msg, global.Conf.Mqtt.EventTopic)
	return err
}

func (u *TpService) CommandReply(token string, msg []byte) error {
	//数据发送至tp的mqtt
	err := MqttSend(token, msg, global.Conf.Mqtt.CommandTopic)
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
