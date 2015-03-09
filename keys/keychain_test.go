package keys

import (
	"crypto/rsa"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

// helpers

func GetKeychain() *Keychain {
	tmpDir, err := ioutil.TempDir("", "keychain_test")
	if err != nil {
		panic(err)
	}
	return NewKeychain(tmpDir)
}

func DestroyKeychain(kc *Keychain) {
	if err := os.RemoveAll(kc.Path); err != nil {
		panic(err)
	}
}

// test

func TestKeychainFindBothPrivateAndPublicKey(t *testing.T) {
	keychain := GetKeychain()
	defer DestroyKeychain(keychain)

	if err := ioutil.WriteFile(path.Join(keychain.Path, "the-key.pem"), testRsaPrivateKey, 0600); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(path.Join(keychain.Path, "the-key.pub"), testRsaPublicKey, 0644); err != nil {
		panic(err)
	}

	key, err := keychain.Find("the-key")
	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}

	if key.Name != "the-key" {
		t.Errorf("unexpected key.Name %#v", key.Name)
	}

	if key.Private.E != rsaKey.E {
		t.Errorf("unexpected key.Private %#v", key.Private)
	}

	if key.Public.E != rsaKey.Public().(*rsa.PublicKey).E {
		t.Errorf("unexpected key.Public %#v", key.Public)
	}
}

func TestKeychainFindPrivateKey(t *testing.T) {
	keychain := GetKeychain()
	defer DestroyKeychain(keychain)

	if err := ioutil.WriteFile(path.Join(keychain.Path, "the-key.pem"), testRsaPrivateKey, 0600); err != nil {
		panic(err)
	}

	key, err := keychain.Find("the-key")
	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}

	if key.Name != "the-key" {
		t.Errorf("unexpected key.Name %#v", key.Name)
	}

	if key.Private.E != rsaKey.E {
		t.Errorf("unexpected key.Private %#v", key.Private)
	}

	if key.Public.E != rsaKey.Public().(*rsa.PublicKey).E {
		t.Errorf("unexpected key.Public %#v", key.Public)
	}
}

func TestKeychainFindPublicKey(t *testing.T) {
	keychain := GetKeychain()
	defer DestroyKeychain(keychain)

	if err := ioutil.WriteFile(path.Join(keychain.Path, "the-key.pub"), testRsaPublicKey, 0644); err != nil {
		panic(err)
	}

	key, err := keychain.Find("the-key")
	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}

	if key.Name != "the-key" {
		t.Errorf("unexpected key.Name %#v", key.Name)
	}

	if key.Private != nil {
		t.Errorf("unexpected key.Private %#v", key.Private)
	}

	if key.Public.E != rsaKey.Public().(*rsa.PublicKey).E {
		t.Errorf("unexpected key.Public %#v", key.Public)
	}
}

func TestKeychainFindUnexist(t *testing.T) {
	keychain := GetKeychain()
	defer DestroyKeychain(keychain)

	_, err := keychain.Find("the-key")
	if err != ErrKeyNotFound {
		t.Errorf("unexpected error %#v", err)
	}
}

func TestKeychainSavePrivateKey(t *testing.T) {
	keychain := GetKeychain()
	defer DestroyKeychain(keychain)

	key, err := LoadKey(testRsaPrivateKey)
	if err != nil {
		panic(err)
	}

	key.Name = "new-key"

	err = keychain.Save(key)

	if err != nil {
		t.Errorf("unexpected error %#v", err.Error())
	}

	filepath := path.Join(keychain.Path, "new-key.pem")

	if fi, err := os.Stat(filepath); err == nil {
		if fi.Mode() != 0600 {
			t.Errorf("unexpected file mode %i", fi.Mode())
		}
	} else {
		t.Errorf("expected file stat fail:", err.Error())
	}

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Errorf("failed to read file %s", err.Error())
	}

	if !reflect.DeepEqual(key.PrivatePem(), bytes) {
		t.Errorf("key file content unexpected %#v", bytes)
	}
}

func TestKeychainSavePublicKey(t *testing.T) {
	keychain := GetKeychain()
	defer DestroyKeychain(keychain)

	key, err := LoadKey(testRsaPublicKey)
	if err != nil {
		panic(err)
	}

	key.Name = "new-key"

	err = keychain.Save(key)

	if err != nil {
		t.Errorf("unexpected error %#v", err.Error())
	}

	filepath := path.Join(keychain.Path, "new-key.pub")

	if _, err := os.Stat(filepath); err != nil {
		t.Errorf("expected file stat fail:", err.Error())
	}

	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Errorf("failed to read file %s", err.Error())
	}

	if !reflect.DeepEqual(key.PublicPem(), bytes) {
		t.Errorf("key file content unexpected %#v", bytes)
	}
}

func TestKeychainSaveAlreadyExist(t *testing.T) {
	keychain := GetKeychain()
	defer DestroyKeychain(keychain)

	if err := ioutil.WriteFile(path.Join(keychain.Path, "the-key.pem"), testRsaPrivateKey, 0600); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(path.Join(keychain.Path, "the-key.pub"), testRsaPublicKey, 0644); err != nil {
		panic(err)
	}

	key, err := LoadKey(testRsaPrivateKey)
	if err != nil {
		panic(err)
	}
	key.Name = "the-key"

	err = keychain.Save(key)

	if err != ErrKeyAlreadyExists {
		t.Errorf("unexpected error %#v", err.Error())
	}
}
