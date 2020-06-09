package netWorker

type Config struct {
	Protocol string
	Address  string
}

func NewConfig() *Config {
	return &Config{
		Protocol: "tcp",
		Address:  "192.168.0.103:62000",
	}
}
