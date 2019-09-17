package utils

import (
	"io/ioutil"
	"crypto/rsa"
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