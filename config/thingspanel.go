package config

type Thingspanel struct {
	Address    string `mapstructure:"address" yaml:"address,omitempty"`
	OffineTime int64  `mapstructure:"offine_time" yaml:"offine_time,omitempty"`
}
