package engine

import (
	"github.com/sorah/etcvault/keys"
	"io/ioutil"
	"os"
	"path"
	"strings"
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

var testKeychain *keys.Keychain

func TestMain(m *testing.M) {
	tmpDir, err := ioutil.TempDir("", "engine_test")
	if err != nil {
		panic(err)
	}

	testKeychain = keys.NewKeychain(tmpDir)

	if err := ioutil.WriteFile(path.Join(testKeychain.Path, "the-key.pem"), testRsaPrivateKey, 0600); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(path.Join(testKeychain.Path, "pubkey.pub"), testRsaPublicKey, 0644); err != nil {
		panic(err)
	}

	defer func() {
		if err := os.RemoveAll(testKeychain.Path); err != nil {
			panic(err)
		}
	}()

	os.Exit(m.Run())
}

func TestTransformPlainToPlain(t *testing.T) {
	engine := NewEngine(testKeychain)

	result, err := engine.Transform("plain text")

	if err != nil {
		t.Errorf("unexpected err: %#v", err)
	}
	if result != "plain text" {
		t.Errorf("unexpected result: %#v", result)
	}
}

func TestTransformPlainRoundtrip(t *testing.T) {
	engine := NewEngine(testKeychain)

	encryptedText, err := engine.Transform("ETCVAULT::plain:the-key:this text should be encrypted::ETCVAULT")
	if err != nil {
		t.Errorf("unexpected err: %#v", err)
	}
	if strings.Index(encryptedText, "this text should be encrypted") != -1 {
		t.Errorf("encrypted text contains original text: %#v", encryptedText)
	}
	if strings.Index(encryptedText, "ETCVAULT::1:the-key::") != 0 {
		t.Errorf("encrypted text unexpected: %#v", encryptedText)
	}

	plainText, err := engine.Transform(encryptedText)
	if err != nil {
		t.Errorf("2 unexpected err: %#v", err)
	}
	if plainText != "this text should be encrypted" {
		t.Errorf("unexpected result: %#v", plainText)
	}
}

func TestTransformV1RoundtripShort(t *testing.T) {
	engine := NewEngine(testKeychain)

	encryptedText, err := engine.Transform("ETCVAULT::plain:the-key:this text should be encrypted::ETCVAULT")
	if err != nil {
		t.Errorf("1 unexpected err: %#v", err)
	}
	if strings.Index(encryptedText, "this text should be encrypted") != -1 {
		t.Errorf("encrypted text contains original text: %#v", encryptedText)
	}
	if strings.Index(encryptedText, "ETCVAULT::1:the-key::") != 0 {
		t.Errorf("encrypted text unexpected: %#v", encryptedText)
	}

	plainText, err := engine.Transform(encryptedText)
	if err != nil {
		t.Errorf("2 unexpected err: %#v", err)
	}
	if plainText != "this text should be encrypted" {
		t.Errorf("unexpected result: %#v", plainText)
	}
}

func TestTransformV1DecryptionShort(t *testing.T) {
	engine := NewEngine(testKeychain)
	decryptedText, err := engine.Transform("ETCVAULT::1:the-key::oXKv3edU7AjUXK1+7+Ng7y5tjByLzMe8MRL2lCxlsE03pHS2AXnd3mvar5dkbgeTU4dY8lcMPYAqRGXi2y9YJ7MD+8vKpkORczLYOBTiSXY8cuttvWY+ffjeJMSsLiHn0tDdtjvCtshSBTe9vLz75yyW8J91DUm9CriHWtQhaXw=::ETCVAULT")

	if err != nil {
		t.Errorf("1 unexpected err: %#v", err)
	}
	if decryptedText != "this text should be encrypted" {
		t.Errorf("unexpected text %#v", decryptedText)
	}
}

func TestTransformV1RoundtripLong(t *testing.T) {
	engine := NewEngine(testKeychain)

	encryptedText, err := engine.Transform("ETCVAULT::plain:the-key:this text is too long so this should be long format aaaaaaaaaaaaaaaaaaaaaaaaaa::ETCVAULT")
	if err != nil {
		t.Errorf("1 unexpected err: %#v", err.Error())
	}
	if strings.Index(encryptedText, "this text is too long so this should be long format aaaaaaaaaaaaaaaaaaaaaaaaaa") != -1 {
		t.Errorf("encrypted text contains original text: %#v", encryptedText)
	}
	if strings.Index(encryptedText, "ETCVAULT::1:the-key:long:") != 0 {
		t.Errorf("encrypted text unexpected: %#v", encryptedText)
	}

	plainText, err := engine.Transform(encryptedText)
	if err != nil {
		t.Errorf("2 unexpected err: %#v", err)
	}
	if plainText != "this text is too long so this should be long format aaaaaaaaaaaaaaaaaaaaaaaaaa" {
		t.Errorf("unexpected result: %#v", plainText)
	}
}

func TestTransformV1DecryptionLong(t *testing.T) {
	engine := NewEngine(testKeychain)
	decryptedText, err := engine.Transform("ETCVAULT::1:the-key:long:JRrn3XxO/HJEu/xYblTkxooOGvFkvnHz4AyinTceZMI2ybRbS2TyoOS+fTGZTTdUMnQ0gKhqH/KsCBjtvW/lw+CXEXVooCmpRCRyVYJIu/FH+oarHIGkpDTeJruEVaL1Jlvo0gb9Ea4zeZuKSiabY+puoTHVCEm1sEN8pHE48xA=,6LaTIBRfKOMBfHq/2JaF/ooeVe97GLGe5gJB8DBYMI30q8mynk9DoMgDKX4ROoiUXatFhSS20hvIIZEUwt62qN7ksivXSb9OybZwU22h6Kw=::ETCVAULT")
	if err != nil {
		t.Errorf("1 unexpected err: %#v", err)
	}
	if decryptedText != "this text is too long so this should be long format aaaaaaaaaaaaaaaaaaaaaaaaaa" {
		t.Errorf("unexpected text %#v", decryptedText)
	}
}
