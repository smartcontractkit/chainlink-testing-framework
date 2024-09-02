package client

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

func GenerateRandomETHKey(password string) (string, error) {
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate private key")
	}
	jsonKey, err := keystore.EncryptKey(&keystore.Key{
		PrivateKey: privateKey,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
	}, password, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return "", errors.Wrap(err, "failed to encrypt the keystore")
	}
	return string(jsonKey), nil
}

func TestStoreNodeKeys(t *testing.T) {
	t.Skip("Need AWS role to be enabled")
	//testPrefix := "TEST_SOAK_KEY_"
	//sm, err := NewAWSSecretsManager("us-west-2", 1*time.Minute)
	//require.NoError(t, err)

	//urls := []string{"..."}
	//keys, err := InjectNewClusterEVMKeys(urls, "pwd")
	//require.NoError(t, err)
	//secrets := keysToSecretsFormat(keys)

	//t.Run("basic CRUD", func(t *testing.T) {
	//	err = sm.StoreSecretsMapWithPrefix(testPrefix, secrets)
	//	require.NoError(t, err)
	//	_, err := sm.GetSecretsMapWithPrefix(testPrefix)
	//	require.NoError(t, err)
	//	require.Equal(t, secrets, NewAWSecrets(map[string]string{
	//		testPrefix + "mysecret":      "1",
	//		testPrefix + "anothersecret": "2",
	//	}))
	//	err = sm.RemoveSecretsByPrefix(testPrefix, true)
	//	require.NoError(t, err)
	//})
}

// Session is the form structure used for authenticating
type Session struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ImportRandomTestKey(url, key, password string) error {
	r := resty.New().SetBaseURL(url).SetDebug(true)
	session := &Session{Email: "notreal@fakeemail.ch", Password: "fj293fbBnlQ!f9vNs"}
	authResp, err := r.R().SetBody(session).Post("/sessions")
	if err != nil {
		return errors.Wrap(err, "failed to authenticate to CL node")
	}
	r.SetCookies(authResp.Cookies())
	_, err = r.R().SetBody(key).Post(fmt.Sprintf("/v2/keys/eth/import?evmChainID=1337&oldpassword=%s", password))
	if err != nil {
		return errors.Wrap(err, "failed to import EVM key to CL node")
	}
	return nil
}

//func keysToSecretsFormat(keys []string) AWSSecrets {
//	m := AWSSecrets{}
//	for i, k := range keys {
//		s := AWSSecret(k)
//		m[fmt.Sprintf("CL_NODE_%d", i)] = &s
//	}
//	return m
//}

func InjectNewClusterEVMKeys(urls []string, pwd string) ([]string, error) {
	keys := make([]string, 0)
	for _, u := range urls {
		key, err := GenerateRandomETHKey(pwd)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
		err = ImportRandomTestKey(u, key, pwd)
		if err != nil {
			return nil, err
		}
	}
	return keys, nil
}
