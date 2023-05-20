package utils

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"
)

type Device struct {
	DeviceId     string `json:"Id"`           //设备id
	AccessToken  string `json:"AccessToken"`  //设备唯一标识
	ProtocolType string `json:"ProtocolType"` //协议类型 1-mqtt
	DeviceType   string `json:"DeviceType"`   //设备类型 1-直连设备
	DeviceConfig DeviceConfig
}
type DeviceConfig struct {
	CommandUrl  string `json:"WebhookAddr"` //设备接收消息url
	OffineTime  int64  `json:"OffineTime"`  //设备离线时间阈值
	DeviceId    string `json:"Id"`          //设备id
	AccessToken string `json:"AccessToken"` //设备唯一标识
	LastMsgTime int64  `json:"LastMsgTime"` //最后一次消息时间
	Status      string `json:"Status"`      //设备状态 0-离线 1-在线
}

// SetLastMsgTime 设置最后一次消息时间
func (d *Device) SetLastMsgTime(ts int64) {
	d.DeviceConfig.LastMsgTime = ts
}

// SetStatus 设置设备状态
func (d *Device) SetStatus(status string) {
	d.DeviceConfig.Status = status
}

type Resdata struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Device `json:"data"`
}

func GetNowTime() int64 {
	return time.Now().Unix()
}

// ReadFormConfig 读取配置文件
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

//字符串替换非法字符
func ReplaceUserInput(s string) string {
	newStringInput := strings.NewReplacer("\n", " ", "\r", " ")
	return newStringInput.Replace(s)
}
