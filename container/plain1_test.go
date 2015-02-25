package container

import (
	"testing"
)

func TestParsePlain(t *testing.T) {
	container, err := ParsePlain1("ETCVAULT::plain:foo:content::ETCVAULT")

	if err != nil {
		t.Errorf("unexpected err: %#v", err)
	}

	if container.Version() != "plain1" {
		t.Errorf("unexpected container.Version: %#v", container.Version)
	}

	if container.KeyName != "foo" {
		t.Errorf("unexpected container.KeyName: %#v", container.KeyName)
	}

	if container.Content != "content" {
		t.Errorf("unexpected container.Content: %#v", container.Content)
	}
}

func TestParsePlain1(t *testing.T) {
	container, err := ParsePlain1("ETCVAULT::plain1:foo:content::ETCVAULT")

	if err != nil {
		t.Errorf("unexpected err: %#v", err)
	}

	if container.Version() != "plain1" {
		t.Errorf("unexpected container.Version: %#v", container.Version)
	}

	if container.KeyName != "foo" {
		t.Errorf("unexpected container.KeyName: %#v", container.KeyName)
	}

	if container.Content != "content" {
		t.Errorf("unexpected container.Content: %#v", container.Content)
	}
}

func TestParsePlain1Invalid(t *testing.T) {
	_, err := ParseBasic("foo")

	if err != ErrInvalid {
		t.Errorf("unexpected err: %#v", err)
	}
}

func TestParsePlain1NoTrail(t *testing.T) {
	_, err := ParseBasic("ETCVAULT::foo")

	if err != ErrInvalid {
		t.Errorf("unexpected err: %#v", err)
	}
}

func TestParsePlain1NoHead(t *testing.T) {
	_, err := ParseBasic("foo::ETCVAULT")

	if err != ErrInvalid {
		t.Errorf("unexpected err: %#v", err)
	}
}

func TestParsePlain1NoKeyOrContent(t *testing.T) {
	_, err := ParseBasic("ETCVAULT::foo::ETCVAULT")

	if err != ErrParse {
		t.Errorf("unexpected err: %#v", err)
	}
}
