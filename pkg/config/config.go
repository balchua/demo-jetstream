package config

type AppConfiguration struct {
	I Infra `mapstructure:"infra"`
}

type Infra struct {
	SeedPath string `mapstructure:"seedPath"`
	NatsUri  string `mapstructure:"natsUri"`
}
