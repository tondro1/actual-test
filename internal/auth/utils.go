package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func getSecretKey() (*rsa.PrivateKey, error) {
	SECRET_KEY_PATH := os.Getenv("PRIVATE_KEY_PATH")
	data, err := os.ReadFile(SECRET_KEY_PATH)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch priv := priv.(type) {
	case *rsa.PrivateKey:
		return priv, nil
	default:
		return nil, errors.New("private key type is not RSA")
	}

}

func GetPublicKey() (*rsa.PublicKey, error) {
	PUBLIC_KEY_PATH := os.Getenv("PUBLIC_KEY_PATH")
	data, err := os.ReadFile(PUBLIC_KEY_PATH)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("pub key type is not RSA")
	}
}