package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"reflect"
	"testing"
)

var rsaKey rsa.PrivateKey

func TestMain(m *testing.M) {
	priv, err := rsa.GenerateKey(rand.Reader, 1024)
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

func TestPublicPem(t *testing.T) {
	key := NewPrivateKey("foo", &rsaKey)

	der, err := x509.MarshalPKIXPublicKey(key.Public)
	if err != nil {
		t.Errorf("err %#v", err)
	}

	pemBytes := key.PublicPem()

	pem, _ := pem.Decode(pemBytes)

	if pem.Type != "PUBLIC KEY" {
		t.Errorf("pem unexpected type %#v", pem.Type)
	}
	if pem.Headers["Name"] != "foo" {
		t.Errorf("pem unexpected Header['Name'] %#v", pem.Headers["Name"])
	}
	if !reflect.DeepEqual(pem.Bytes, der) {
		t.Errorf("pem unexpected bytes %#v\nbut: %#v", der, pem.Bytes)
	}
}

func TestPrivatePem(t *testing.T) {
	key := NewPrivateKey("foo", &rsaKey)
	der := x509.MarshalPKCS1PrivateKey(key.Private)

	pemBytes := key.PrivatePem()

	pem, _ := pem.Decode(pemBytes)

	if pem.Type != "PRIVATE KEY" {
		t.Errorf("pem unexpected type %#v", pem.Type)
	}
	if pem.Headers["Name"] != "foo" {
		t.Errorf("pem unexpected Header['Name'] %#v", pem.Headers["Name"])
	}
	if !reflect.DeepEqual(pem.Bytes, der) {
		t.Errorf("pem unexpected bytes %#v\nbut: %#v", der, pem.Bytes)
	}
}

func TestPrivatePemOnPublicKey(t *testing.T) {
	pubKey := rsaKey.Public().(*rsa.PublicKey)
	key := NewPublicKey("foo", pubKey)

	pemBytes := key.PrivatePem()

	if pemBytes != nil {
		t.Errorf("unexpected value %#v", pemBytes)
	}
}
