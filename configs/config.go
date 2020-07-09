package configs

type Config struct {
	Protocol string
	Address  string
}

func GetNetworkConfig() *Config {
	return &Config{
		Protocol: "tcp",
		Address:  ":62000",
	}
}

type ConfigRabbitMQ struct {
	AMQPConnectionURL string
}

func GetRMQConfig() *ConfigRabbitMQ {
	return &ConfigRabbitMQ{
		//TODO: вынести в переменные окружения
		AMQPConnectionURL: "amqp://guest:guest@rabbitmq:5672/",
		//AMQPConnectionURL: "amqp://guest:guest@127.0.0.1:5672/",
	}
}
