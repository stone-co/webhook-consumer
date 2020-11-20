package configuration

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/square/go-jose.v2"
)

// Config defines the service configuration
type Config struct {
	HTTPConfig     HTTPConfig
	PrivateKeyPath string `envconfig:"PRIVATE_KEY_PATH" default:"tests/partner/mykey.pem"`
	PrivateKey     interface{}
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

	config.PrivateKey, err = LoadPrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to read public key: %v", err)
	}

	return &config, nil
}

// LoadPrivateKey loads a private key from PEM/DER/JWK-encoded data.
func LoadPrivateKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	var priv interface{}
	priv, err0 := x509.ParsePKCS1PrivateKey(input)
	if err0 == nil {
		return priv, nil
	}

	priv, err1 := x509.ParsePKCS8PrivateKey(input)
	if err1 == nil {
		return priv, nil
	}

	priv, err2 := x509.ParseECPrivateKey(input)
	if err2 == nil {
		return priv, nil
	}

	jwk, err3 := LoadJSONWebKey(input, false)
	if err3 == nil {
		return jwk, nil
	}

	return nil, fmt.Errorf("square/go-jose: parse error, got '%s', '%s', '%s' and '%s'", err0, err1, err2, err3)
}

func LoadJSONWebKey(json []byte, pub bool) (*jose.JSONWebKey, error) {
	var jwk jose.JSONWebKey
	err := jwk.UnmarshalJSON(json)
	if err != nil {
		return nil, err
	}

	if !jwk.Valid() {
		return nil, errors.New("invalid JWK key")
	}

	if jwk.IsPublic() != pub {
		return nil, errors.New("priv/pub JWK key mismatch")
	}

	return &jwk, nil
}
