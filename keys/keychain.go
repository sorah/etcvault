package keys

import (
	"errors"
	"os"
	"path"
)

var ErrKeyNotFound = errors.New("couldn't find specified key")

type Keychain struct {
	Path string
}

func NewKeychain(path string) *Keychain {
	return &Keychain{
		Path: path,
	}
}

func (keychain *Keychain) Find(name string) (*Key, error) {
	privateKeyPath := path.Join(keychain.Path, name+".pem")
	publicKeyPath := path.Join(keychain.Path, name+".pub")

	if _, err := os.Stat(privateKeyPath); err == nil {
		return LoadKeyFromFile(privateKeyPath)
	} else if _, err := os.Stat(publicKeyPath); err == nil {
		return LoadKeyFromFile(publicKeyPath)
	} else {
		return nil, ErrKeyNotFound
	}
}

// func (keychain *Keychain) Save(string name) error {
// }
