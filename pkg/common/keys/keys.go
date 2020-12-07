package keys

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/square/go-jose.v2"
)

const (
	FileLocation = "file://"
	URLLocation  = "url://"
)

type Config struct {
	PrivateKey          interface{}
	VerificationKeyList []*jose.JSONWebKey
}

func LoadKeys(privateKeyPath, publicKeyLocation string) (*Config, error) {
	var config Config

	keyBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %v", privateKeyPath, err)
	}

	config.PrivateKey, err = LoadPrivateKey(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	config.VerificationKeyList, err = loadVerificationKeyList(publicKeyLocation)
	if err != nil {
		return nil, fmt.Errorf("loading verification key %s: %v", publicKeyLocation, err)
	}

	return &config, nil
}

func loadVerificationKeyList(location string) ([]*jose.JSONWebKey, error) {
	var keyList []*jose.JSONWebKey
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

	if len(keyList) == 0 {
		return nil, fmt.Errorf("empty key list")
	}

	return keyList, nil
}

func loadVerificationKeyListFromFile(fileList string) ([]*jose.JSONWebKey, error) {
	result := []*jose.JSONWebKey{}
	for _, file := range strings.Split(fileList, ";") {
		file = strings.TrimSpace(file)
		if file == "" {
			break
		}
		keyBytes, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("reading file %s: %v", file, err)
		}

		verificationKey, err := LoadPublicKeyFromJWK(keyBytes)
		if err != nil {
			return nil, fmt.Errorf("unable to read public key: %v", err)
		}
		result = append(result, verificationKey)
	}

	if result == nil {
		return nil, fmt.Errorf("empty file list")
	}

	return result, nil
}

func loadVerificationKeyListFromURL(serviceURL string) ([]*jose.JSONWebKey, error) {
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
		Keys []*jose.JSONWebKey `json:"keys"`
	}

	var r responseBody
	if err = json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("unable to unmarshal body: %v", err)
	}

	return r.Keys, nil
}
