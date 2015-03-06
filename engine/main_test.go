package engine

import (
	"github.com/sorah/etcvault/keys"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var testRsaPrivateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
Name: the-key

MIICXAIBAAKBgQDE0H3AjeUvlOA5ueZ1q6hukF+aRFbW2h8qW2OIw88+EN4qLani
lTvTUO3V91hGhHe2CnnUOey1iAHnSPGx66XW3oN/Wuk+wK1tg1ivcCLHIOlRu22g
8DuS8TC92jhjkFVCgGasXNFGECiyF6J9WsYrF6F/OKvUVpEjWgyRMPMMuQIDAQAB
AoGAMOlbhyH8ZhHKk64GfxHU/v00NSNsrWJxwlYJ63A2LceFXtgQUzYhMwf2w2j/
8C51jbEWy85FbGvLhU4UetIEWW0OK5Y+J2juGD0ez1FX+EzmiO+khpGtYQ6OY56a
3g4FPsUuCj1gw2oBDDQ2e38RyqY9Nj3PWo4H5Y7ZbSWwSQ0CQQDSNABnC7AiM2K3
5uXqZiXx68RoLrYtGkXhgyZBIUZ+g6nbhBqpPEI9pql55yCjmx/zeY6VVipOffO2
EEUpdnG/AkEA77G9SK8lqxMeH+GRL70jYNXBqdxYhKrWlFzom+VrHIyo//limocH
dPJiEEIyPJQXeru2r2mWxVg98q+j3CUvhwJAIzebKaiHpfM+Atmog5EBonqBuYK5
+ux/8LxsWFUe3mtoteJ4JQp3fqTBmC7lBQQkYkJnZRW+mM/5WPN44u15OQJBAJPO
Wbehcav9vPzR3vK+QjurdKHnI5qjsnCInlPL8/IF9wzp3tkFXR7LfJckCtB6TcQ8
Ttn6VaPZ11F456WQNK8CQETVQARcp/v4bWtVHfJKyBcx92FkclVNXae5aHpmvIjI
LUu9LpYOrkcaL1d7SFPhWZUsI+crYKuLAb9tXG/AnJY=
-----END RSA PRIVATE KEY-----`)

var testRsaPublicKey = []byte(`-----BEGIN RSA PUBLIC KEY-----
Name: the-key

MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDE0H3AjeUvlOA5ueZ1q6hukF+a
RFbW2h8qW2OIw88+EN4qLanilTvTUO3V91hGhHe2CnnUOey1iAHnSPGx66XW3oN/
Wuk+wK1tg1ivcCLHIOlRu22g8DuS8TC92jhjkFVCgGasXNFGECiyF6J9WsYrF6F/
OKvUVpEjWgyRMPMMuQIDAQAB
-----END RSA PUBLIC KEY-----`)

var testKeychain *keys.Keychain

func TestMain(m *testing.M) {
	tmpDir, err := ioutil.TempDir("", "engine_test")
	if err != nil {
		panic(err)
	}

	testKeychain = keys.NewKeychain(tmpDir)

	if err := ioutil.WriteFile(path.Join(testKeychain.Path, "the-key.pem"), testRsaPrivateKey, 0600); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(path.Join(testKeychain.Path, "pubkey.pub"), testRsaPublicKey, 0644); err != nil {
		panic(err)
	}

	defer func() {
		if err := os.RemoveAll(testKeychain.Path); err != nil {
			panic(err)
		}
	}()

	os.Exit(m.Run())
}
