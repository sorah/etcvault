package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
	"testing"
)

var rsaKey rsa.PrivateKey

func TestMain(m *testing.M) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	rsaKey = *priv

	os.Exit(m.Run())
}

func TestNewPrivateKey(t *testing.T) {
	key := NewPrivateKey("foo", &rsaKey)

	if key.Name != "foo" {
		t.Errorf("unexpected key.Name %#v", key.Name)
	}

	if key.Private.E != rsaKey.E {
		t.Errorf("unexpected key.Private %#v", key.Private)
	}

	if key.Public.E != rsaKey.Public().(*rsa.PublicKey).E {
		t.Errorf("unexpected key.Public", key.Public)
	}
}

func TestNewPublicKey(t *testing.T) {
	pubKey := rsaKey.Public().(*rsa.PublicKey)
	key := NewPublicKey("foo", pubKey)

	if key.Name != "foo" {
		t.Errorf("unexpected key.Name %#v", key.Name)
	}

	if key.Private != nil {
		t.Errorf("unexpected key.Private %#v", key.Private)
	}

	if key.Public.E != rsaKey.Public().(*rsa.PublicKey).E {
		t.Errorf("unexpected key.Public", key.Public)
	}
}

func TestGenerateKey(t *testing.T) {
	key, err := GenerateKey("foo", 1024)

	if err != nil {
		t.Errorf("error %#v", err)
	}

	if key.Name != "foo" {
		t.Errorf("unexpected key.Name %#v", key.Name)
	}
}
