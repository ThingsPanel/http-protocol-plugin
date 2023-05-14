package config

type Thingspanel struct {
	Address  string `mapstructure:"address" yaml:"address,omitempty"`
	Username string `mapstructure:"username" yaml:"username,omitempty"`
	Password string `mapstructure:"password" yaml:"password,omitempty"`
}
