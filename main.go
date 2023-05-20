package main

import (
	"http-procotol-plugin/global"
	"http-procotol-plugin/initialize"
	"http-procotol-plugin/service"
)

func main() {
	initialize.Conf()               //加载配置文件
	global.Mqtt = initialize.Mqtt() //连接mqtt
	service.MqttSubscribe()         //订阅mqtt
	service.OnOfflineCron()         //定时检查设备是否离线
	initialize.RunServer()          //启动服务

}
