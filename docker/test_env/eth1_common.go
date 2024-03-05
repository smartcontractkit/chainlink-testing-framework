package test_env

import (
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

const (
	RootFundingAddr   = `0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266`
	RootFundingWallet = `{"address":"f39fd6e51aad88f6f4ce6ab8827279cfffb92266","crypto":{"cipher":"aes-128-ctr","ciphertext":"c36afd6e60b82d6844530bd6ab44dbc3b85a53e826c3a7f6fc6a75ce38c1e4c6","cipherparams":{"iv":"f69d2bb8cd0cb6274535656553b61806"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"80d5f5e38ba175b6b89acfc8ea62a6f163970504af301292377ff7baafedab53"},"mac":"f2ecec2c4d05aacc10eba5235354c2fcc3776824f81ec6de98022f704efbf065"},"id":"e5c124e9-e280-4b10-a27b-d7f3e516b408","version":3}`

	DEFAULT_EVM_NODE_HTTP_PORT = "8544"
	DEFAULT_EVM_NODE_WS_PORT   = "8545"
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
