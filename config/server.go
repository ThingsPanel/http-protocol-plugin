package config

//server相关配置
type Server struct {
	Addr string `mapstructure:"addr"  yaml:"addr,omitempty"`
}
