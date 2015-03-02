package keys

import (
	"errors"
	"os"
	"path"
)

var ErrKeyNotFound = errors.New("couldn't find specified key")

type Keychain struct {
	Path  string
	Cache map[string]*Key
}

func NewKeychain(path string) *Keychain {
	return &Keychain{
		Path:  path,
		Cache: make(map[string]*Key),
	}
}

func (keychain *Keychain) Find(name string) (*Key, error) {
	if key, ok := keychain.Cache[name]; ok {
		return key, nil
	}

	privateKeyPath := path.Join(keychain.Path, name+".pem")
	publicKeyPath := path.Join(keychain.Path, name+".pub")

	if _, err := os.Stat(privateKeyPath); err == nil {
		key, err := LoadKeyFromFile(privateKeyPath)
		if err != nil {
			return nil, err
		}
		keychain.Cache[name] = key
		return key, nil
	} else if _, err := os.Stat(publicKeyPath); err == nil {
		key, err := LoadKeyFromFile(publicKeyPath)
		if err != nil {
			return nil, err
		}
		return key, nil
	} else {
		return nil, ErrKeyNotFound
	}
}

// func (keychain *Keychain) Save(string name) error {
// }
