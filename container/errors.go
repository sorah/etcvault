package container

import (
	"errors"
)

var ErrParse = errors.New("couldn't parse")
var ErrInvalid = errors.New("it's not in container form (invalid)")
var ErrDifferentVersion = errors.New("it's in different version")
var ErrUnknownVersion = errors.New("Unknown version")
