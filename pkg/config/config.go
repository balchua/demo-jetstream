package config

type AppConfiguration struct {
	I Infra      `mapstructure:"infra"`
	P Publisher  `mapstructure:"publish"`
	M Monitor    `mapstructure:"monitor"`
	S Subscriber `mapstructure:"subscribe"`
	T Tracing    `mapstructure:"tracing"`
	A ApiConfig  `mapstructure:"api"`
}

type Infra struct {
	SeedPath string `mapstructure:"seedPath"`
	NatsUri  string `mapstructure:"natsUri"`
}

type Publisher struct {
	SeedPath string `mapstructure:"seedPath"`
	NatsUri  string `mapstructure:"natsUri"`
}

type Subscriber struct {
	SeedPath          string `mapstructure:"seedPath"`
	NatsUri           string `mapstructure:"natsUri"`
	SleepTimeInMillis int    `mapstructure:"sleepTimeInMillis"`
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

type Tracing struct {
	JaegerUrl   string `mapstructure:"jaeger-url"`
	ServiceName string `mapstructure:"service-name"`
}

type ApiConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
