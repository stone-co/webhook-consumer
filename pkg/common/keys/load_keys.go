package keys

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"gopkg.in/square/go-jose.v2"
)

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

// LoadPublicKey loads a public key from PEM/DER/JWK-encoded data.
func LoadPublicKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	var aggErr error

	// Try to load SubjectPublicKeyInfo
	pub, err := x509.ParsePKIXPublicKey(input)
	if err == nil {
		return pub, nil
	}
	aggErr = err

	cert, err := x509.ParseCertificate(input)
	if err == nil {
		return cert.PublicKey, nil
	}
	aggErr = fmt.Errorf("%s: %w", aggErr, err)

	jwk, err := LoadJSONWebKey(data, true)
	if err == nil {
		return jwk, nil
	}

	return nil, fmt.Errorf("%s: %w", aggErr, err)
}

// LoadPublicKeyFromJWK loads a public key from JWK-encoded data.
func LoadPublicKeyFromJWK(data []byte) (*jose.JSONWebKey, error) {
	jwk, err := LoadJSONWebKey(data, true)
	if err == nil {
		return jwk, nil
	}

	return nil, fmt.Errorf("square/go-jose: jwk parse error, got '%s'", err)
}
