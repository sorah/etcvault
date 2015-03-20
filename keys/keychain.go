package keys

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var ErrKeyNotFound = errors.New("couldn't find specified key")
var ErrKeyAlreadyExists = errors.New("another key already exists with same name")

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

func (keychain *Keychain) Save(key *Key) error {
	if _, err := keychain.Find(key.Name); err == nil {
		return ErrKeyAlreadyExists
	}
	if key.Private == nil {
		publicKeyPath := path.Join(keychain.Path, key.Name+".pub")
		return ioutil.WriteFile(publicKeyPath, key.PublicPem(), 0644)
	} else {
		privateKeyPath := path.Join(keychain.Path, key.Name+".pem")
		return ioutil.WriteFile(privateKeyPath, key.PrivatePem(), 0600)
	}
	return nil
}

func (keychain *Keychain) List() []string {
	namesMap := make(map[string]bool)

	addNames := func(ext string) {
		matches, err := filepath.Glob(path.Join(keychain.Path, fmt.Sprintf("*.%s", ext)))
		if err != nil {
			log.Printf("error looking for key list (%s): %s", ext, err.Error())
			return
		}
		for _, keyPath := range matches {
			name := strings.TrimSuffix(path.Base(keyPath), fmt.Sprintf(".%s", ext))
			namesMap[name] = true
		}
	}

	addNames("pub")
	addNames("pem")

	names := make([]string, 0, len(namesMap))
	for name, _ := range namesMap {
		names = append(names, name)
	}
	return names
}

func (keychain *Keychain) ListForEncryption() []string {
	matches, err := filepath.Glob(path.Join(keychain.Path, "*.pem"))
	if err != nil {
		log.Printf("error looking for key list (pem): %s", err.Error())
		return []string{}
	}
	names := make([]string, 0, len(matches))
	for _, keyPath := range matches {
		name := strings.TrimSuffix(path.Base(keyPath), ".pem")
		names = append(names, name)
	}
	return names
}

func (keychain *Keychain) ListForDecryption() []string {
	return keychain.List()
}
