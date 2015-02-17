package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var ErrMissingPem = errors.New("invalid pem (couldn't decode)")
var ErrInvalidPem = errors.New("invalid pem (type should be RSA PUBLIC KEY, PUBLIC KEY, RSA PRIVATE KEY, or PRIVATE KEY)")
var ErrNotRsaKey = errors.New("invalid pem (key is not RSA public key or private key)")

type Key struct {
	Name    string
	Public  *rsa.PublicKey
	Private *rsa.PrivateKey
}

func NewPrivateKey(name string, rsaPrivateKey *rsa.PrivateKey) *Key {
	pubKey := rsaPrivateKey.Public().(*rsa.PublicKey)
	return &Key{
		Name:    name,
		Public:  pubKey,
		Private: rsaPrivateKey,
	}
}

func NewPublicKey(name string, rsaPublicKey *rsa.PublicKey) *Key {
	return &Key{
		Name:   name,
		Public: rsaPublicKey,
	}
}

func LoadKey(name string, pemBytes []byte) (*Key, error) {
	pem, _ := pem.Decode(pemBytes)
	if pem == nil {
		return nil, ErrMissingPem
	}

	switch pem.Type {
	case "PUBLIC KEY", "RSA PUBLIC KEY":
		parsedKey, err := x509.ParsePKIXPublicKey(pem.Bytes)
		if err != nil {
			return nil, err
		}

		var pubKey *rsa.PublicKey
		var ok bool
		if pubKey, ok = parsedKey.(*rsa.PublicKey); !ok {
			return nil, ErrNotRsaKey
		}

		return NewPublicKey(name, pubKey), nil

	case "PRIVATE KEY", "RSA PRIVATE KEY":
		privateKey, err := x509.ParsePKCS1PrivateKey(pem.Bytes)
		if err != nil {
			return nil, err
		}

		return NewPrivateKey(name, privateKey), nil

	default:
		return nil, ErrInvalidPem
	}
}

func GenerateKey(name string, bits int) (*Key, error) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	return NewPrivateKey(name, rsaKey), nil
}

func (key *Key) PublicPem() []byte {
	der, err := x509.MarshalPKIXPublicKey(key.Public)

	// normally MarshalPKIXPublicKey doesn't say error for rsa.PublicKey
	if err != nil {
		panic(err)
	}

	block := &pem.Block{
		Type:    "PUBLIC KEY",
		Headers: map[string]string{"Name": key.Name},
		Bytes:   der,
	}

	return pem.EncodeToMemory(block)
}

func (key *Key) PrivatePem() []byte {
	if key.Private == nil {
		return nil
	}

	der := x509.MarshalPKCS1PrivateKey(key.Private)
	block := &pem.Block{
		Type:    "PRIVATE KEY",
		Headers: map[string]string{"Name": key.Name},
		Bytes:   der,
	}

	return pem.EncodeToMemory(block)
}
