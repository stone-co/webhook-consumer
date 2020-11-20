package configuration

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stone-co/webhook-consumer/pkg/common/keys"
)

// Config defines the service configuration
type Config struct {
	HTTPConfig          HTTPConfig
	PrivateKeyPath      string `envconfig:"PRIVATE_KEY_PATH" default:"tests/partner/mykey.pem"`
	PrivateKey          interface{}
	PublicKeyLocation   string `envconfig:"PUBLIC_KEY_PATH" default:"file://./tests/stone/mykey.pub"`
	VerificationKeyList []interface{}
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

	const FileLocation = "file://"
	if strings.HasPrefix(config.PublicKeyLocation, FileLocation) {
		config.VerificationKeyList, err = loadVerificationKeyFromFile(strings.TrimPrefix(config.PublicKeyLocation, FileLocation))
		if err != nil {
			return nil, fmt.Errorf("loading verification key from file %s: %v", config.PublicKeyLocation, err)
		}
	} else {
		return nil, fmt.Errorf("invalid public key location: %s", config.PublicKeyLocation)
	}

	return &config, nil
}

func loadVerificationKeyFromFile(file string) ([]interface{}, error) {

	keyBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %v", file, err)
	}

	verificationKey, err := keys.LoadPublicKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to read public key: %v", err)
	}

	return []interface{}{verificationKey}, nil
}
