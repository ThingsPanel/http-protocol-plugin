package initialize

import (
	"fmt"
	"http-procotol-plugin/global"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-basic/uuid"
	"github.com/panjf2000/ants/v2"
)

func Mqtt() mqtt.Client {
	broker := os.Getenv("MQTT_HOST")
	if broker == "" {
		broker = global.Conf.Mqtt.Broker
	}
	clientid := uuid.New()
	username := global.Conf.Mqtt.Username
	password := global.Conf.Mqtt.Password

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		fmt.Printf("Connect lost: %v", err)
	}
	opts := mqtt.NewClientOptions()
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(clientid)
	opts.AddBroker(broker)
	opts.SetAutoReconnect(true) //自动重连
	opts.SetOrderMatters(false)
	opts.OnConnectionLost = connectLostHandler
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		log.Println("MQTT客户端连接成功...", broker)
	})
	p, _ := ants.NewPool(global.Conf.Mqtt.SubscribePool) //设置并发池
	log.Println("mqtt客户端订阅处理的并发池大小为", global.Conf.Mqtt.SubscribePool)
	opts.SetDefaultPublishHandler(func(c mqtt.Client, m mqtt.Message) {
		_ = p.Submit(func() {

		})
	})
	mqtt_client := mqtt.NewClient(opts)
	if token := mqtt_client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("mqtt客户端连接异常...", global.Conf.Mqtt.Broker, token.Error())
		os.Exit(1)
	}
	if token := mqtt_client.Subscribe(global.Conf.Mqtt.TopicToSubscribe, 0, nil); token.Wait() && token.Error() != nil {
		log.Println("mqtt订阅异常...", global.Conf.Mqtt.TopicToSubscribe, token.Error())
		os.Exit(1)
	} else {
		log.Println("mqtt订阅成功...", global.Conf.Mqtt.TopicToSubscribe)
	}
	return mqtt_client
}
