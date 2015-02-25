package container

import (
	"fmt"
	"strings"
)

type Plain1 struct {
	KeyName string
	Content string
}

func ParsePlain1(str string) (*Plain1, error) {
	basic, err := ParseBasic(str)
	if err != nil {
		return nil, err
	}

	if !(basic.Version == "plain" || basic.Version == "plain1") {
		return nil, ErrDifferentVersion
	}

	keyAndContent := strings.SplitN(basic.Content, ":", 2)

	if len(keyAndContent) < 2 {
		return nil, ErrParse
	}

	return &Plain1{
		KeyName: keyAndContent[0],
		Content: keyAndContent[1],
	}, nil
}

func (container *Plain1) Version() string {
	return "plain1"
}

func (container *Plain1) String() string {
	return fmt.Sprintf("ETCVAULT::plain:%s:%s::ETCVAULT", container.KeyName, container.Content)
}
