package configuration

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stone-co/webhook-consumer/pkg/common/keys"
)

// Config defines the service configuration
type Config struct {
	HTTPConfig      HTTPConfig
	PrivateKeyPath  string `envconfig:"PRIVATE_KEY_PATH" default:"tests/partner/mykey.pem"`
	PrivateKey      interface{}
	PublicKeyPath   string `envconfig:"PUBLIC_KEY_PATH" default:"tests/stone/mykey.pub"`
	VerificationKey interface{}
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

	keyBytes, err := ioutil.ReadFile(config.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %v", config.PrivateKeyPath, err)
	}

	config.PrivateKey, err = keys.LoadPrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	keyBytes, err = ioutil.ReadFile(config.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %v", config.PublicKeyPath, err)
	}

	config.VerificationKey, err = keys.LoadPublicKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to read public key: %v", err)
	}

	return &config, nil
}
