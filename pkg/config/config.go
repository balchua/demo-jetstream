package config

type AppConfiguration struct {
	I Infra     `mapstructure:"infra"`
	P Publisher `mapstructure:"publish"`
	M Monitor   `mapstructure:"monitor"`
}

type Infra struct {
	SeedPath string `mapstructure:"seedPath"`
	NatsUri  string `mapstructure:"natsUri"`
}

type Publisher struct {
	SeedPath string `mapstructure:"seedPath"`
	NatsUri  string `mapstructure:"natsUri"`
}

type Monitor struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Scheme       string `mapstructure:"scheme"`
	Account      string `mapstructure:"account"`
	ConsumerName string `mapstructure:"consumerName"`
	StreamName   string `mapstructure:"streamName"`
	PollSeconds  int    `mapstructure:"pollSeconds"`
}
