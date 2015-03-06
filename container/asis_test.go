package container

import (
	"testing"
)

func TestParseAsis(t *testing.T) {
	container, err := ParseAsis("ETCVAULT::asis:content::ETCVAULT")

	if err != nil {
		t.Errorf("unexpected err: %#v", err)
	}

	if container.Version() != "asis" {
		t.Errorf("unexpected container.Version: %#v", container.Version)
	}

	if container.Content != "content" {
		t.Errorf("unexpected container.Content: %#v", container.Content)
	}
}
