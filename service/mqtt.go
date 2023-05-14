package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"http-procotol-plugin/global"
	"http-procotol-plugin/utils"
	"log"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttPayload struct {
	Token  string `json:"token"`
	Values []byte `json:"values"`
}

//发送消息
func MqttSend(accessToken string, msg []byte, topic string) (err error) {
	payload := &mqttPayload{}
	payload.Token = accessToken
	payload.Values = msg
	//转换成json
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("json转换失败", err)
		return err
	}
	fmt.Println("发送消息...", string(msg), "topic:", topic)
	t := global.Mqtt.Publish(topic, 1, false, payloadJson)
	if t.Error() != nil {
		log.Println("发送消息失败...", payloadJson, t.Error())
	} else {
		log.Println("发送...", payloadJson)
	}
	return t.Error()
}

//订阅加回调函数
func MqttSubscribe() mqtt.Token {
	return global.Mqtt.Subscribe(global.Conf.Mqtt.TopicToSubscribe, 1, DeviceMsgFunc)
}

func DeviceMsgFunc(client mqtt.Client, msg mqtt.Message) {
	log.Println("订阅的新消息是：", msg.Topic(), string(msg.Payload()))
	//将消息发送至设备
	devicemsg := &mqttPayload{}
	err := json.Unmarshal(msg.Payload(), &devicemsg)
	if err != nil {
		log.Println("json转换失败", err)
		return
	}
	//判断设备是否存在
	if _, ok := global.DevicesMap.Load(devicemsg.Token); !ok {
		log.Println("设备不存在,开始添加设备", devicemsg.Token)
		//从tp获取设备信息，将token储存在map里
		if err := TpDeviceAccessToken(devicemsg.Token); err != nil {
			log.Println("添加设备失败", err)
			return
		}
	}
	//获取设备信息
	device, _ := global.DevicesMap.Load(devicemsg.Token)
	//发送至设备
	if err := PostJSON(device.(utils.Device).DeviceConfig.CommandUrl, devicemsg.Values); err != nil {
		log.Println("发送失败", err)
	}

}

func PostJSON(url string, data interface{}) error {
	// 将数据转换成JSON格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// 发送POST请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("发送失败 %d", resp.StatusCode)
	}

	return nil
}
