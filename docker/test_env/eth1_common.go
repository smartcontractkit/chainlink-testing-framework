package test_env

import (
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type keyStoreAndExtraData struct {
	ks           *keystore.KeyStore
	minerAccount *accounts.Account
	extraData    []byte
}

func generateKeystoreAndExtraData(keystoreDir string) (keyStoreAndExtraData, error) {
	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	minerAccount, err := ks.NewAccount("")
	if err != nil {
		return keyStoreAndExtraData{}, err
	}

	minerAddr := strings.Replace(minerAccount.Address.Hex(), "0x", "", 1)
	signerBytes, err := hex.DecodeString(minerAddr)
	if err != nil {
		return keyStoreAndExtraData{}, err
	}

	zeroBytes := make([]byte, 32)                      // Create 32 zero bytes
	extradata := append(zeroBytes, signerBytes...)     // Concatenate zero bytes and signer address
	extradata = append(extradata, make([]byte, 65)...) // Concatenate 65 more zero bytes

	return keyStoreAndExtraData{
		ks:           ks,
		minerAccount: &minerAccount,
		extraData:    extradata,
	}, nil
}
