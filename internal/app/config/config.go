package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	ConfigSpiders string     `yaml:"config_spiders"`
	HTTPServer    HTTPServer `yaml:"http_server"`
	Nats          Nats       `yaml:"nats"`
	Redis         Redis      `yaml:"redis"`
}

type HTTPServer struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Redis struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Nats struct {
	Server      string `yaml:"server"`
	StreamName  string `yaml:"stream_name"`
	Subject     string `yaml:"subject"`
	DurableName string `yaml:"durable_name"`

	WorkerCount int `yaml:"worker_count"`
}

func LoadConfig(filename string) (*Config, error) {
	cfg := Config{}

	err := cleanenv.ReadConfig(filename, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
