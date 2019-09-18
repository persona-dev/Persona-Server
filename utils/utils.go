package utils

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
)

// ReadPublicKey look up RSA public key from ./public-key.pem.
func ReadPublicKey() (*rsa.PublicKey, error) {
	Key, err := ioutil.ReadFile("public-key.pem")
	if err != nil {
		return nil, err
	}
	ParsedKey, err := jwt.ParseRSAPublicKeyFromPEM(Key)
	return ParsedKey, err
}

func LoadPrivateKey() (*rsa.PrivateKey, error) {
	Key, err := ioutil.ReadFile("private-key.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to road private key: %s", err)
	}
	PrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(Key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Privatekey: %s", err)
	}
	return PrivateKey, nil
}
