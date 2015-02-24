package container

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type V1 struct {
	KeyName    string
	ContentKey []byte
	Content    []byte
}

func ParseV1(str string) (*V1, error) {
	basic, err := ParseBasic(str)
	if err != nil {
		return nil, err
	}

	if basic.Version != "1" {
		return nil, ErrDifferentVersion
	}

	parts := strings.SplitN(basic.Content, ":", 3) // key name, format, content

	if len(parts) < 3 {
		return nil, ErrParse
	}

	keyName := parts[0]
	format := parts[1]
	contentPart := parts[2]

	var contentKey []byte
	var content []byte
	if format == "long" {
		contentKeyAndContent := strings.SplitN(contentPart, ",", 2)
		if len(parts) < 2 {
			return nil, ErrParse
		}

		contentKey, err = base64.StdEncoding.DecodeString(contentKeyAndContent[0])
		if err != nil {
			return nil, err
		}

		content, err = base64.StdEncoding.DecodeString(contentKeyAndContent[1])
		if err != nil {
			return nil, err
		}
	} else {
		content, err = base64.StdEncoding.DecodeString(contentPart)
		if err != nil {
			return nil, err
		}
	}

	return &V1{
		KeyName:    keyName,
		ContentKey: contentKey,
		Content:    content,
	}, nil
}

func (container *V1) Version() string {
	return "1"
}

func (container *V1) String() string {
	encodedContent := base64.StdEncoding.EncodeToString(container.Content)
	if container.ContentKey == nil {
		return fmt.Sprintf("ETCVAULT::1:%s::%s::ETCVAULT", container.KeyName, encodedContent)
	} else {
		encodedContentKey := base64.StdEncoding.EncodeToString(container.ContentKey)
		return fmt.Sprintf("ETCVAULT::1:%s:long:%s,%s::ETCVAULT", container.KeyName, encodedContentKey, encodedContent)
	}
}
