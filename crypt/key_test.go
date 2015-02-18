package crypt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"reflect"
	"testing"
)

var testRsaPrivateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
Name: the-key

MIICXAIBAAKBgQDE0H3AjeUvlOA5ueZ1q6hukF+aRFbW2h8qW2OIw88+EN4qLani
lTvTUO3V91hGhHe2CnnUOey1iAHnSPGx66XW3oN/Wuk+wK1tg1ivcCLHIOlRu22g
8DuS8TC92jhjkFVCgGasXNFGECiyF6J9WsYrF6F/OKvUVpEjWgyRMPMMuQIDAQAB
AoGAMOlbhyH8ZhHKk64GfxHU/v00NSNsrWJxwlYJ63A2LceFXtgQUzYhMwf2w2j/
8C51jbEWy85FbGvLhU4UetIEWW0OK5Y+J2juGD0ez1FX+EzmiO+khpGtYQ6OY56a
3g4FPsUuCj1gw2oBDDQ2e38RyqY9Nj3PWo4H5Y7ZbSWwSQ0CQQDSNABnC7AiM2K3
5uXqZiXx68RoLrYtGkXhgyZBIUZ+g6nbhBqpPEI9pql55yCjmx/zeY6VVipOffO2
EEUpdnG/AkEA77G9SK8lqxMeH+GRL70jYNXBqdxYhKrWlFzom+VrHIyo//limocH
dPJiEEIyPJQXeru2r2mWxVg98q+j3CUvhwJAIzebKaiHpfM+Atmog5EBonqBuYK5
+ux/8LxsWFUe3mtoteJ4JQp3fqTBmC7lBQQkYkJnZRW+mM/5WPN44u15OQJBAJPO
Wbehcav9vPzR3vK+QjurdKHnI5qjsnCInlPL8/IF9wzp3tkFXR7LfJckCtB6TcQ8
Ttn6VaPZ11F456WQNK8CQETVQARcp/v4bWtVHfJKyBcx92FkclVNXae5aHpmvIjI
LUu9LpYOrkcaL1d7SFPhWZUsI+crYKuLAb9tXG/AnJY=
-----END RSA PRIVATE KEY-----`)

var testRsaPublicKey = []byte(`-----BEGIN RSA PUBLIC KEY-----
Name: the-key

MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDE0H3AjeUvlOA5ueZ1q6hukF+a
RFbW2h8qW2OIw88+EN4qLanilTvTUO3V91hGhHe2CnnUOey1iAHnSPGx66XW3oN/
Wuk+wK1tg1ivcCLHIOlRu22g8DuS8TC92jhjkFVCgGasXNFGECiyF6J9WsYrF6F/
OKvUVpEjWgyRMPMMuQIDAQAB
-----END RSA PUBLIC KEY-----`)

var testRsaPrivateKeyNoHeader = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDE0H3AjeUvlOA5ueZ1q6hukF+aRFbW2h8qW2OIw88+EN4qLani
lTvTUO3V91hGhHe2CnnUOey1iAHnSPGx66XW3oN/Wuk+wK1tg1ivcCLHIOlRu22g
8DuS8TC92jhjkFVCgGasXNFGECiyF6J9WsYrF6F/OKvUVpEjWgyRMPMMuQIDAQAB
AoGAMOlbhyH8ZhHKk64GfxHU/v00NSNsrWJxwlYJ63A2LceFXtgQUzYhMwf2w2j/
8C51jbEWy85FbGvLhU4UetIEWW0OK5Y+J2juGD0ez1FX+EzmiO+khpGtYQ6OY56a
3g4FPsUuCj1gw2oBDDQ2e38RyqY9Nj3PWo4H5Y7ZbSWwSQ0CQQDSNABnC7AiM2K3
5uXqZiXx68RoLrYtGkXhgyZBIUZ+g6nbhBqpPEI9pql55yCjmx/zeY6VVipOffO2
EEUpdnG/AkEA77G9SK8lqxMeH+GRL70jYNXBqdxYhKrWlFzom+VrHIyo//limocH
dPJiEEIyPJQXeru2r2mWxVg98q+j3CUvhwJAIzebKaiHpfM+Atmog5EBonqBuYK5
+ux/8LxsWFUe3mtoteJ4JQp3fqTBmC7lBQQkYkJnZRW+mM/5WPN44u15OQJBAJPO
Wbehcav9vPzR3vK+QjurdKHnI5qjsnCInlPL8/IF9wzp3tkFXR7LfJckCtB6TcQ8
Ttn6VaPZ11F456WQNK8CQETVQARcp/v4bWtVHfJKyBcx92FkclVNXae5aHpmvIjI
LUu9LpYOrkcaL1d7SFPhWZUsI+crYKuLAb9tXG/AnJY=
-----END RSA PRIVATE KEY-----`)

var testRsaPublicKeyNoHeader = []byte(`-----BEGIN RSA PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDE0H3AjeUvlOA5ueZ1q6hukF+a
RFbW2h8qW2OIw88+EN4qLanilTvTUO3V91hGhHe2CnnUOey1iAHnSPGx66XW3oN/
Wuk+wK1tg1ivcCLHIOlRu22g8DuS8TC92jhjkFVCgGasXNFGECiyF6J9WsYrF6F/
OKvUVpEjWgyRMPMMuQIDAQAB
-----END RSA PUBLIC KEY-----`)

var testEcdsaPrivateKey = []byte(`-----BEGIN PRIVATE KEY-----
MGgCAQEEHF+pP6QjO+LH97mzJlaiqZ1y5DynKEjUSXy7hVSgBwYFK4EEACGhPAM6
AASwJR+5yutBOBaKlxjheM+VPm4kfeXoxnjN85OHAfYeyEPS95kZZKqbpvX8d8NF
Z4+YLPZEMaBs7g==
-----END PRIVATE KEY-----`)

var testEcdsaPublicKey = []byte(`-----BEGIN PUBLIC KEY-----
ME4wEAYHKoZIzj0CAQYFK4EEACEDOgAEsCUfucrrQTgWipcY4XjPlT5uJH3l6MZ4
zfOThwH2HshD0veZGWSqm6b1/HfDRWePmCz2RDGgbO4=
-----END PUBLIC KEY-----`)

var rsaKey rsa.PrivateKey

func TestMain(m *testing.M) {
	pem, _ := pem.Decode(testRsaPrivateKey)
	if pem == nil {
		panic("invalid pem")
	}

	priv, err := x509.ParsePKCS1PrivateKey(pem.Bytes)
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
