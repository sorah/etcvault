package container

import (
	"fmt"
	"strings"
)

type Basic struct {
	Version string
	Content string
}

func ParseBasic(str string) (*Basic, error) {
	// ETCVAULT:::::ETCVAULT (at least 21 chars)
	if len(str) < 21 {
		return nil, ErrInvalid
	}
	if strings.Index(str, "ETCVAULT::") == -1 || strings.Index(str, "::ETCVAULT") != (len(str)-10) {
		return nil, ErrInvalid
	}
	inner := str[10 : len(str)-10]
	versionAndContent := strings.SplitN(inner, ":", 2)

	if len(versionAndContent) < 2 {
		return nil, ErrParse
	}

	return &Basic{
		Version: versionAndContent[0],
		Content: versionAndContent[1],
	}, nil
}

func (container *Basic) String() string {
	return fmt.Sprintf("ETCVAULT::%s:%s::ETCVAULT", container.Version, container.Content)
}
