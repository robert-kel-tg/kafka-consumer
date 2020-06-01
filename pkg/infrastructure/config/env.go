package config

import (
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type (
	Config struct {
		DB *DBConfig
		LoggerConfig
		ConsumerConfig
	}

	LoggerConfig struct {
		LogLevel string `envconfig:"LOG_LEVEL" default:"debug"`
	}

	DBConfig struct {
		ConnectionURL     string        `envconfig:"DATABASE_DSN" required:"true"`
		Driver            string        `envconfig:"DATABASE_DRIVER" default:"postgres"`
		MigrationsPath    string        `envconfig:"MIGRATIONS_PATH" default:"./migrations"`
	}

 	ConsumerConfig struct {
		Topics []string `envconfig:"TOPICS"`
	}
)

func LoadConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (c *Config) AddKafkaConf() (kafka.ConfigMap, error) {
	config := kafka.ConfigMap{
		//TODO change to vars
		"bootstrap.servers": "broker:29092",
		"group.id":          "demoGroupID",
		"auto.offset.reset": "earliest",
		"go.events.channel.enable": true,
		"enable.partition.eof": false,
		"session.timeout.ms": 6000,
	}

	return config, nil
}
