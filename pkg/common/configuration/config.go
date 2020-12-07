package configuration

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config defines the service configuration
type Config struct {
	HTTPConfig     HTTPConfig
	PrivateKeyPath string `envconfig:"PRIVATE_KEY_PATH" default:"tests/partner/fakekey.pem"`
	// PublicKeyLocation can be used to specify a file or a URL.
	// To specify a file: "file://./tests/stone/fakekey1.pub.jwt"
	// To specify a URL: "url://https://sandbox-api.openbank.stone.com.br/api/v1/discovery/keys"
	PublicKeyLocation string `envconfig:"PUBLIC_KEY_PATH" default:"url://https://sandbox-api.openbank.stone.com.br/api/v1/discovery/keys"`
	// NotifierList has stdout and proxy availables.
	NotifierList  string `envconfig:"NOTIFIER_LIST" default:"stdout"`
	ProxyNotifier ProxyNotifierConfig
}

type HTTPConfig struct {
	Port            int           `envconfig:"API_PORT" default:"3000"`
	ShutdownTimeout time.Duration `envconfig:"API_SHUTDOWN_TIMEOUT" default:"5s"`
}

type ProxyNotifierConfig struct {
	Url     string        `envconfig:"PROXY_NOTIFIER_URL"`
	Timeout time.Duration `envconfig:"PROXY_NOTIFIER_TIMEOUT" default:"10s"`
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

func (cfg Config) String() string {
	return fmt.Sprintf("port:[%d] shutdown_timeout:[%s] private_key_path:[%s] public_key_location:[%s] notifier_list:[%s]",
		cfg.HTTPConfig.Port, cfg.HTTPConfig.ShutdownTimeout, cfg.PrivateKeyPath, cfg.PublicKeyLocation, cfg.NotifierList)
}
