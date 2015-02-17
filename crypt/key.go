package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

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
