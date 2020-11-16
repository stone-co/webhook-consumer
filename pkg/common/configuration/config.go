package configuration

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config defines the service configuration
type Config struct {
	HTTPConfig HTTPConfig
}

type HTTPConfig struct {
	Port            int           `envconfig:"API_PORT" default:"3000"`
	ShutdownTimeout time.Duration `envconfig:"API_SHUTDOWN_TIMEOUT" default:"5s"`
}

func LoadConfig() (*Config, error) {
	var config Config
	prefix := ""
	err := envconfig.Process(prefix, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
