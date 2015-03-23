package engine

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/sorah/etcvault/container"
	"github.com/sorah/etcvault/keys"
)

var ErrNoPrivateKey = errors.New("no private key provided")
var ErrTooShortKey = errors.New("key too short; couldn't generate 16, 24, and 32 bytes aes key")

type Transformable interface {
	Transform(text string) (string, error)
	TransformEtcdJsonResponse(jsonData []byte) ([]byte, error)
	GetKeychain() *keys.Keychain
}

type Engine struct {
	Keychain *keys.Keychain
}

func NewEngine(keychain *keys.Keychain) *Engine {
	return &Engine{
		Keychain: keychain,
	}
}

func (engine *Engine) GetKeychain() *keys.Keychain {
	return engine.Keychain
}

func (engine *Engine) Transform(text string) (string, error) {
	s, _, e := engine.TransformAndParse(text)
	return s, e
}

func (engine *Engine) TransformAndParse(text string) (string, container.Container, error) {
	// FIXME: test for this
	rawContainer, err := container.Parse(text)

	if err != nil {
		if err == container.ErrInvalid {
			return text, nil, nil
		} else {
			return "", nil, err
		}
	}

	switch c := rawContainer.(type) {
	case *container.Plain1:
		result, err := engine.TransformPlain1(c)
		return result, c, err
	case *container.Asis:
		result, err := engine.TransformAsis(c)
		return result, c, err
	case *container.V1:
		result, err := engine.TransformV1(c)
		return result, c, err
	}
	// shouldnt reach
	panic(fmt.Errorf("BUG: unsupported container type %#v", rawContainer))
}

func (engine *Engine) TransformPlain1(c *container.Plain1) (string, error) {
	key, err := engine.Keychain.Find(c.KeyName)
	if err != nil {
		return "", err
	}

	encryptedContent, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, key.Public, []byte(c.Content), []byte{})
	if err == rsa.ErrMessageTooLong {
		return engine.transformPlain1Long(key, c)
	}
	if err != nil {
		return "", err
	}

	result := &container.V1{
		KeyName: key.Name,
		Content: encryptedContent,
	}

	return result.String(), nil
}

func (engine *Engine) TransformAsis(c *container.Asis) (string, error) {
	return c.Content, nil
}

func (engine *Engine) transformPlain1Long(key *keys.Key, c *container.Plain1) (string, error) {
	hash := sha256.New()

	rsaMaxLength := ((key.Public.N.BitLen() + 7) / 8) - (2 * hash.Size()) - 2
	contentKeyLength := 32
	if rsaMaxLength < contentKeyLength {
		contentKeyLength = 24
	}
	if rsaMaxLength < contentKeyLength {
		contentKeyLength = 16
	}
	if rsaMaxLength < contentKeyLength {
		return "", ErrTooShortKey
	}

	contentKey := make([]byte, contentKeyLength)
	if _, err := rand.Read(contentKey); err != nil {
		return "", err
	}

	encryptedContentKey, err := rsa.EncryptOAEP(hash, rand.Reader, key.Public, contentKey, []byte{})
	if err != nil {
		return "", err
	}

	cipher, err := aes.NewCipher(contentKey)
	if err != nil {
		return "", err
	}

	content := []byte(c.Content)
	encryptedContent := *(encryptAesWithPkcs7Padding(&cipher, &content))

	result := &container.V1{
		KeyName:    key.Name,
		ContentKey: encryptedContentKey,
		Content:    encryptedContent,
	}
	return result.String(), nil
}

func (engine *Engine) TransformV1(c *container.V1) (string, error) {
	if c.ContentKey == nil {
		return engine.transformV1Short(c)
	} else {
		return engine.transformV1Long(c)
	}
}

func (engine *Engine) transformV1Short(c *container.V1) (string, error) {
	key, err := engine.Keychain.Find(c.KeyName)
	if err != nil {
		return "", err
	}
	if key.Private == nil {
		return "", ErrNoPrivateKey
	}

	hash := sha256.New()
	decryptedContent, err := rsa.DecryptOAEP(hash, rand.Reader, key.Private, c.Content, []byte{})
	if err != nil {
		return "", err
	}
	return string(decryptedContent), nil
}

func (engine *Engine) transformV1Long(c *container.V1) (string, error) {
	key, err := engine.Keychain.Find(c.KeyName)
	if err != nil {
		return "", err
	}
	if key.Private == nil {
		return "", ErrNoPrivateKey
	}

	hash := sha256.New()
	decryptedContentKey, err := rsa.DecryptOAEP(hash, rand.Reader, key.Private, c.ContentKey, []byte{})
	if err != nil {
		return "", err
	}

	aes, err := aes.NewCipher(decryptedContentKey)
	if err != nil {
		return "", err
	}

	decryptedContent, err := decryptAesWithPkcs7Padding(&aes, &c.Content)
	if err != nil {
		return "", nil
	}

	return string(*decryptedContent), nil
}
