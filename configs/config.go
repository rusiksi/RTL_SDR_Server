package configs

type Config struct {
	Protocol string
	Address  string
}

func GetNetworkConfig() *Config {
	return &Config{
		Protocol: "tcp",
		Address:  "192.168.0.103:62000",
	}
}

type ConfigRabbitMQ struct {
	AMQPConnectionURL string

}

func GetRMQConfig() *ConfigRabbitMQ {
	return &ConfigRabbitMQ{
		AMQPConnectionURL: "amqp://guest:guest@localhost:5672/",
	}
}
