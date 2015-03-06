package engine

import (
	"strings"
	"testing"
)

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
