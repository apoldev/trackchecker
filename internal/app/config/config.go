package config

import "github.com/ilyakaznacheev/cleanenv"

const (
	BrokerNats  = "nats"
	BrokerKafka = "kafka"
)

type Config struct {
	ConfigSpiders string     `yaml:"config_spiders"`
	Broker        string     `yaml:"broker" env-default:"nats"`
	GRPCServer    GRPCServer `yaml:"grpc_server"`
	HTTPServer    HTTPServer `yaml:"http_server"`
	Nats          Nats       `yaml:"nats"`
	Kafka         Kafka      `yaml:"kafka"`
	Redis         Redis      `yaml:"redis"`
}

type GRPCServer struct {
	Port int `yaml:"port" env-default:"8070"`
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
	Server       string `yaml:"server"`
	StreamName   string `yaml:"stream_name"`
	Subject      string `yaml:"subject"`
	DurableName  string `yaml:"durable_name"`
	WorkerCount  int    `yaml:"worker_count"`
	JSMaxPending int    `yaml:"js_max_pending"`
}

type Kafka struct {
	Server        string `yaml:"server"`
	ConsumerGroup string `yaml:"consumer_group"`
	Topic         string `yaml:"topic"`
	WorkerCount   int    `yaml:"worker_count"`
}

func LoadConfig(filename string) (*Config, error) {
	cfg := Config{}

	err := cleanenv.ReadConfig(filename, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
