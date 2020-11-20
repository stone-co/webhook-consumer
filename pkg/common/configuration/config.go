package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/stone-co/webhook-consumer/pkg/common/keys"
	"gopkg.in/square/go-jose.v2"
)

const (
	FileLocation = "file://"
	URLLocation  = "url://"
)

// Config defines the service configuration
type Config struct {
	HTTPConfig     HTTPConfig
	PrivateKeyPath string `envconfig:"PRIVATE_KEY_PATH" default:"tests/partner/mykey.pem"`
	PrivateKey     interface{}
	// PublicKeyLocation   string `envconfig:"PUBLIC_KEY_PATH" default:"file://./tests/stone/mykey.pub"`
	PublicKeyLocation   string `envconfig:"PUBLIC_KEY_PATH" default:"url://https://sandbox-api.openbank.stone.com.br/api/v1/discovery/keys"`
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

	config.VerificationKeyList, err = loadVerificationKeyList(config.PublicKeyLocation)
	if err != nil {
		return nil, fmt.Errorf("loading verification key %s: %v", config.PublicKeyLocation, err)
	}

	return &config, nil
}

func loadVerificationKeyList(location string) ([]interface{}, error) {
	var keyList []interface{}
	var err error

	if strings.HasPrefix(location, FileLocation) {
		keyList, err = loadVerificationKeyListFromFile(strings.TrimPrefix(location, FileLocation))
		if err != nil {
			return nil, fmt.Errorf("loading verification key from file %s: %v", location, err)
		}
	} else if strings.HasPrefix(location, URLLocation) {
		keyList, err = loadVerificationKeyListFromURL(strings.TrimPrefix(location, URLLocation))
		if err != nil {
			return nil, fmt.Errorf("loading verification key from file %s: %v", location, err)
		}
	} else {
		return nil, fmt.Errorf("invalid public key location: %s", location)
	}

	return keyList, nil
}

func loadVerificationKeyListFromFile(file string) ([]interface{}, error) {
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

func loadVerificationKeyListFromURL(serviceURL string) ([]interface{}, error) {
	keysURL, err := url.Parse(serviceURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse url %s: %v", serviceURL, err)
	}

	client := http.DefaultClient
	response, err := client.Get(keysURL.String())
	if err != nil {
		return nil, fmt.Errorf("unable to get url keys %s: %v", keysURL.String(), err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read body: %v", err)
	}

	type responseBody struct {
		Keys []jose.JSONWebKey `json:"keys"`
	}

	var r responseBody
	if err = json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("unable to unmarshal body: %v", err)
	}

	keyList := []interface{}{}

	for i := range r.Keys {
		keyList = append(keyList, &r.Keys[i])
	}

	return keyList, nil
}
