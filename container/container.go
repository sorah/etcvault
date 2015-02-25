package container

type Container interface {
	Version() string
	String() string
}

func Parse(str string) (Container, error) {
	basic, err := ParseBasic(str)

	if err != nil {
		return nil, err
	}

	switch basic.Version {
	case "1":
		return ParseV1(str)
	case "plain1", "plain":
		return ParsePlain1(str)
	default:
		return nil, ErrUnknownVersion
	}
}
