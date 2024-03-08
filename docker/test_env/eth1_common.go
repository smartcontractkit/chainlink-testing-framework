package test_env

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type keyStoreAndExtraData struct {
	ks             *keystore.KeyStore
	minerAccount   *accounts.Account
	accountsToFund []string
	extraData      []byte
}

var EXECUTION_CONTAINER_TYPES = []ContainerType{ContainerType_Geth, ContainerType_Nethermind, ContainerType_Erigon, ContainerType_Besu}

func generateKeystoreAndExtraData(keystoreDir string, extraAddressesToFound []string) (keyStoreAndExtraData, error) {
	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	minerAccount, err := ks.NewAccount("")
	if err != nil {
		return keyStoreAndExtraData{}, err
	}

	minerAddr := strings.Replace(minerAccount.Address.Hex(), "0x", "", 1)

	i := 1
	var accounts []string
	for addr, v := range FundingAddresses {
		if v == "" || addr == minerAddr {
			continue
		}
		f, err := os.Create(fmt.Sprintf("%s/%s", keystoreDir, fmt.Sprintf("key%d", i)))
		if err != nil {
			return keyStoreAndExtraData{}, err
		}
		_, err = f.WriteString(v)
		if err != nil {
			return keyStoreAndExtraData{}, err
		}
		i++
		accounts = append(accounts, addr)
	}

	extraAddresses := []string{}
	for _, addr := range extraAddressesToFound {
		extraAddresses = append(extraAddresses, strings.Replace(addr, "0x", "", 1))
	}

	accounts = append(accounts, minerAddr)
	accounts = append(accounts, extraAddresses...)
	accounts, _, err = deduplicateAddresses(accounts)
	if err != nil {
		return keyStoreAndExtraData{}, err
	}

	signerBytes, err := hex.DecodeString(minerAddr)
	if err != nil {
		return keyStoreAndExtraData{}, err
	}

	zeroBytes := make([]byte, 32)                      // Create 32 zero bytes
	extradata := append(zeroBytes, signerBytes...)     // Concatenate zero bytes and signer address
	extradata = append(extradata, make([]byte, 65)...) // Concatenate 65 more zero bytes

	return keyStoreAndExtraData{
		ks:             ks,
		minerAccount:   &minerAccount,
		accountsToFund: accounts,
		extraData:      extradata,
	}, nil
}
