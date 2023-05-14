package config

//配置文件结构体
type Conf struct {
	Server      `mapstructure:"server" yaml:"server"`
	Mqtt        `mapstructure:"mqtt" yaml:"mqtt,omitempty"`
	Thingspanel `mapstructure:"thingspanel" yaml:"thingspanel,omitempty"`
}
