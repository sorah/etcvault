package container

import (
	"bytes"
	"testing"
)

func TestParseForV1(t *testing.T) {
	rawResult, err := Parse("ETCVAULT::1:key::aGVsbG8=::ETCVAULT")

	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}

	result, ok := rawResult.(*V1)
	if !ok {
		t.Errorf("V1 has not returned")
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

func TestParseForPlain1(t *testing.T) {
	rawResult, err := Parse("ETCVAULT::plain1:key:helo::ETCVAULT")

	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}

	result, ok := rawResult.(*Plain1)
	if !ok {
		t.Errorf("Plain1 has not returned")
	}

	if result.Version() != "plain1" {
		t.Errorf("unexpected version %#v", result.Version())
	}

	if result.KeyName != "key" {
		t.Errorf("unexpected KeyName %#v", result.KeyName)
	}

	if result.Content != "helo" {
		t.Errorf("unexpected Content %#v", result.Content)
	}
}

func TestParseForUnknown(t *testing.T) {
	result, err := Parse("ETCVAULT::unknown:XXX::ETCVAULT")

	if result != nil {
		t.Errorf("unexpected result %#v", result)
	}

	if err != ErrUnknownVersion {
		t.Errorf("unexpected error %#v", err)
	}
}
