package container

import (
	"testing"
)

func TestParseBasic(t *testing.T) {
	container, err := ParseBasic("ETCVAULT::42:foo::ETCVAULT")

	if err != nil {
		t.Errorf("unexpected err: %#v", err)
	}

	if container.Version != "42" {
		t.Errorf("unexpected container.Version: %#v", container.Version)
	}

	if container.Content != "foo" {
		t.Errorf("unexpected container.Content: %#v", container.Content)
	}
}

func TestParseBasicInvalid(t *testing.T) {
	_, err := ParseBasic("foo")

	if err != ErrInvalid {
		t.Errorf("unexpected err: %#v", err)
	}
}

func TestParseBasicNoTrail(t *testing.T) {
	_, err := ParseBasic("ETCVAULT::foo")

	if err != ErrInvalid {
		t.Errorf("unexpected err: %#v", err)
	}
}

func TestParseBasicNoHead(t *testing.T) {
	_, err := ParseBasic("foo::ETCVAULT")

	if err != ErrInvalid {
		t.Errorf("unexpected err: %#v", err)
	}
}

func TestParseBasicVersionOrContentMissing(t *testing.T) {
	_, err := ParseBasic("ETCVAULT::foo::ETCVAULT")

	if err != ErrParse {
		t.Errorf("unexpected err: %#v", err)
	}
}
