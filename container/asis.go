package container

import (
	"fmt"
)

type Asis struct {
	Content string
}

func ParseAsis(str string) (*Asis, error) {
	basic, err := ParseBasic(str)
	if err != nil {
		return nil, err
	}

	if !(basic.Version == "asis") {
		return nil, ErrDifferentVersion
	}

	return &Asis{
		Content: basic.Content,
	}, nil
}

func (container *Asis) Version() string {
	return "asis"
}

func (container *Asis) String() string {
	return fmt.Sprintf("ETCVAULT::asis:%s::ETCVAULT", container.Content)
}
