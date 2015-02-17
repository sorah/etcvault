package crypt

import (
	"crypto/rand"
	"crypto/rsa"
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
