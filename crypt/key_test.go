package crypt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"reflect"
	"testing"
)

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

func TestLoadKeyPublic(t *testing.T) {
	key, err := LoadKey(testRsaPublicKey)

	if err != nil {
		t.Errorf("error %#v", err)
	}

	if key.Name != "the-key" {
		t.Errorf("unexpected key.Name %#v", key.Name)
	}

	if key.Private != nil {
		t.Errorf("unexpected key.Private %#v", key.Private)
	}

	if key.Public.E != rsaKey.Public().(*rsa.PublicKey).E {
		t.Errorf("unexpected key.Public", key.Public)
	}
}

func TestLoadKeyPublicNoHeader(t *testing.T) {
	key, err := LoadKey(testRsaPublicKeyNoHeader)

	if err != nil {
		t.Errorf("error %#v", err)
	}

	if key.Name != "" {
		t.Errorf("unexpected key.Name %#v", key.Name)
	}

	if key.Private != nil {
		t.Errorf("unexpected key.Private %#v", key.Private)
	}

	if key.Public.E != rsaKey.Public().(*rsa.PublicKey).E {
		t.Errorf("unexpected key.Public", key.Public)
	}
}

func TestLoadKeyPublicButNotRsa(t *testing.T) {
	key, err := LoadKey(testEcdsaPublicKey)

	if err != ErrNotRsaKey {
		t.Errorf("unexpected error %#v", err)
	}

	if key != nil {
		t.Errorf("unexpected key %#v", err)
	}
}

func TestLoadKeyPublicButInvalid(t *testing.T) {
	// broken
	key, err := LoadKey([]byte(`-----BEGIN RSA PUBLIC KEY-----
Wbehcav9vPzR3vK+QjurdKHnI5qjsnCInlPL8/IF9wzp3tkFXR7LfJckCtB6TcQ8
Ttn6VaPZ11F456WQNK8CQETVQARcp/v4bWtVHfJKyBcx92FkclVNXae5aHpmvIjI
LUu9LpYOrkcaL1d7SFPhWZUsI+crYKuLAb9tXG/AnJY=
-----END RSA PRIVATE KEY-----`))

	if err == nil {
		t.Errorf("unexpected error %#v", err)
	}

	if key != nil {
		t.Errorf("unexpected key %#v", err)
	}
}

func TestLoadKeyPrivate(t *testing.T) {
	key, err := LoadKey(testRsaPrivateKey)

	if err != nil {
		t.Errorf("error %#v", err)
	}

	if key.Name != "the-key" {
		t.Errorf("unexpected key.Name %#v", key.Name)
	}

	if key.Private == nil {
		t.Errorf("unexpected key.Private %#v", key.Private)
	}

	if key.Private.E != rsaKey.E {
		t.Errorf("unexpected key.Private %#v", key.Private)
	}
}

func TestLoadKeyPrivateNoHeader(t *testing.T) {
	key, err := LoadKey(testRsaPrivateKeyNoHeader)

	if err != nil {
		t.Errorf("error %#v", err)
	}

	if key.Name != "" {
		t.Errorf("unexpected key.Name %#v", key.Name)
	}

	if key.Private == nil {
		t.Errorf("unexpected key.Private %#v", key.Private)
	}

	if key.Private.E != rsaKey.E {
		t.Errorf("unexpected key.Private %#v", key.Private)
	}
}

func TestLoadKeyPrivateButInvalid(t *testing.T) {
	// broken
	key, err := LoadKey([]byte(`-----BEGIN RSA PRIVATE KEY-----
GSIb3DQEBAQUAA4GNADCBiQKBgQDE0H3AjeUvlOA5ueZ1q6hukF+aRFbW2h8qW2O
Iw88+EN4qLanilTvTUO3V91hGhHe2CnnUOey1iAHnSPGx66XW3oNWuk+wK1tg1iv
cCLHIOlRu22g8DuS8TC92jhjkFVCgGasXNFGECiyF6J9WsYrF6FOKvUVpEjWgyRM
PMMuQIDAQAB
-----END RSA PUBLIC KEY-----`))

	if err == nil {
		t.Errorf("unexpected error %#v", err)
	}

	if key != nil {
		t.Errorf("unexpected key %#v", err)
	}
}

func TestLoadKeyMissing(t *testing.T) {
	key, err := LoadKey([]byte{})

	if err != ErrMissingPem {
		t.Errorf("unexpected error %#v", err)
	}

	if key != nil {
		t.Errorf("unexpected key %#v", err)
	}
}

func TestLoadKeyInvalid(t *testing.T) {
	key, err := LoadKey([]byte(`-----BEGIN SOMETHING KEY-----
PMMuQIDAQAB
-----END SOMETHING KEY-----`))
	if err != ErrMissingPem {
		t.Errorf("unexpected error %#v", err)
	}

	if key != nil {
		t.Errorf("unexpected key %#v", err)
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

	if pem == nil {
		t.Errorf("couldn't decode pem: %#v", pemBytes)
	}

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

	if pem == nil {
		t.Errorf("couldn't decode pem: %#v", pemBytes)
	}

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
