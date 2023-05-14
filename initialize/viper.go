package initialize

import (
	"http-procotol-plugin/global"
	"log"

	"github.com/spf13/viper"
)

//viper一次性读取配置文件到struct
func Conf() {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal("no config file!")
		} else {
			log.Fatal("read config file error!")
		}
	}
	if err := v.Unmarshal(&global.Conf); err != nil {
		panic("config failed!")
	}
	log.Println("配置文件加载成功...")
}
