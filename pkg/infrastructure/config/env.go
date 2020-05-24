package config

import "github.com/kelseyhightower/envconfig"

type (
	Config struct {
		DB *DBConfig
		LoggerConfig
	}

	LoggerConfig struct {
		LogLevel string `envconfig:"LOG_LEVEL" default:"debug"`
	}

	DBConfig struct {
		ConnectionURL     string        `envconfig:"DATABASE_DSN" required:"true"`
		Driver            string        `envconfig:"DATABASE_DRIVER" default:"postgres"`
		MigrationsPath    string        `envconfig:"MIGRATIONS_PATH" default:"./migrations"`
	}
)

func LoadConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
