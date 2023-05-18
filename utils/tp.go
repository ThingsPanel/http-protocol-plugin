package utils

import (
	"encoding/json"
	"log"
	"os"
)

type Device struct {
	DeviceId     string `json:"Id"`
	AccessToken  string `json:"AccessToken"`
	ProtocolType string `json:"ProtocolType"`
	DeviceType   string `json:"DeviceType"` //设备类型 1-直连设备
	DeviceConfig DeviceConfig
}
type DeviceConfig struct {
	CommandUrl  string `json:"WebhookAddr"` //设备接收消息url
	OffineTime  string `json:"OffineTime"`  //设备离线时间阈值
	DeviceId    string `json:"Id"`
	AccessToken string `json:"AccessToken"`
}

type Resdata struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Device `json:"data"`
}

func ReadFormConfig() interface{} {
	filePtr, err := os.Open("./form_config.json")
	if err != nil {
		log.Println("文件打开失败...", err.Error())
		return nil
	}
	defer filePtr.Close()
	var info interface{}
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&info)
	if err != nil {
		log.Println("解码失败", err.Error())
		return info
	} else {
		log.Println("读取文件[form_config.json]成功...")
		return info
	}
}
