package global

import (
	"http-procotol-plugin/config"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	Conf       config.Conf
	Mqtt       mqtt.Client
	DevicesMap sync.Map
)
