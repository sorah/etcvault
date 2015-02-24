package container

import (
	"bytes"
	"testing"
)

func TestV1ParseShort(t *testing.T) {
	result, err := ParseV1("ETCVAULT::1:key::aGVsbG8=::ETCVAULT")

	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}

	if result.Version() != "1" {
		t.Errorf("unexpected version %#v", result.Version())
	}

	if result.KeyName != "key" {
		t.Errorf("unexpected KeyName %#v", result.KeyName)
	}

	if result.ContentKey != nil {
		t.Errorf("unexpected ContentKey %#v", result.ContentKey)
	}

	if !bytes.Equal(result.Content, []byte("hello")) {
		t.Errorf("unexpected Content %#v == %#v", result.Content)
	}
}

func TestV1ParseLong(t *testing.T) {
	result, err := ParseV1("ETCVAULT::1:key:long:aG9sYQ==,aGVsbG8=::ETCVAULT")

	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}

	if result.Version() != "1" {
		t.Errorf("unexpected version %#v", result.Version())
	}

	if result.KeyName != "key" {
		t.Errorf("unexpected KeyName %#v", result.KeyName)
	}

	if !bytes.Equal(result.ContentKey, []byte(`hola`)) {
		t.Errorf("unexpected ContentKey %#v", result.ContentKey)
	}

	if !bytes.Equal(result.Content, []byte(`hello`)) {
		t.Errorf("unexpected Content %#v", result.Content)
	}
}

func TestV1ParseInvalid(t *testing.T) {
	result, err := ParseV1("hello")

	if result != nil {
		t.Errorf("unexpected result %#v", result)
	}
	if err != ErrInvalid {
		t.Errorf("unexpected error %#v", err)
	}
}

func TestV1ParseError(t *testing.T) {
	result, err := ParseV1("ETCVAULT::1::ETCVAULT")

	if result != nil {
		t.Errorf("unexpected result %#v", result)
	}
	if err != ErrParse {
		t.Errorf("unexpected error %#v", err)
	}
}

func TestV1ParseNotV1(t *testing.T) {
	result, err := ParseV1("ETCVAULT::42:foo::ETCVAULT")

	if result != nil {
		t.Errorf("unexpected result %#v", result)
	}
	if err != ErrDifferentVersion {
		t.Errorf("unexpected error %#v", err)
	}
}
