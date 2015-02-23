package keys

import (
	"crypto/rsa"
	"io/ioutil"
	"os"
	"path"
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
