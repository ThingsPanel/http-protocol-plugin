package config

type Mqtt struct {
	Broker           string `mapstructure:"broker" yaml:"broker,omitempty"`
	Username         string `mapstructure:"username" yaml:"username,omitempty"`
	Password         string `mapstructure:"password" yaml:"password,omitempty"`
	AttributesTopic  string `mapstructure:"attributes_topic" yaml:"attributes_topic,omitempty"`
	EventTopic       string `mapstructure:"event_topic" yaml:"event_topic,omitempty"`
	CommandTopic     string `mapstructure:"command_topic" yaml:"command_topic,omitempty"`
	TopicToSubscribe string `mapstructure:"topic_to_subscribe" yaml:"topic_to_subscribe,omitempty"`
	SubscribePool    int    `mapstructure:"subscribe_pool" yaml:"subscribe_pool,omitempty"`
}
