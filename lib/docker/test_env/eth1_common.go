package test_env

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/slice"
)

type keyStoreAndExtraData struct {
	ks             *keystore.KeyStore
	minerAccount   *accounts.Account
	accountsToFund []string
	extraData      []byte
}

var FundingAddresses = map[string]string{
	"f39fd6e51aad88f6f4ce6ab8827279cfffb92266": `{"address":"f39fd6e51aad88f6f4ce6ab8827279cfffb92266","crypto":{"cipher":"aes-128-ctr","ciphertext":"c36afd6e60b82d6844530bd6ab44dbc3b85a53e826c3a7f6fc6a75ce38c1e4c6","cipherparams":{"iv":"f69d2bb8cd0cb6274535656553b61806"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"80d5f5e38ba175b6b89acfc8ea62a6f163970504af301292377ff7baafedab53"},"mac":"f2ecec2c4d05aacc10eba5235354c2fcc3776824f81ec6de98022f704efbf065"},"id":"e5c124e9-e280-4b10-a27b-d7f3e516b408","version":3}`,
	"70997970c51812dc3a010c7d01b50e0d17dc79c8": `{"address":"70997970c51812dc3a010c7d01b50e0d17dc79c8","crypto":{"cipher":"aes-128-ctr","ciphertext":"f8183fa00bc112645d3e23e29a233e214f7c708bf49d72750c08af88ad76c980","cipherparams":{"iv":"796d08e3e1f71bde89ed826abda96cda"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"03c864a22a1f7b06b1da12d8b93e024ac144f898285907c58b2abc135fc8a35c"},"mac":"5fe91b1a1821c0d9f85dfd582354ead9612e9a7e9adc38b06a2beff558c119ac"},"id":"d2cab765-5e30-42ae-bb91-f090d9574fae","version":3}`,
	"3c44cdddb6a900fa2b585dd299e03d12fa4293bc": `{"address":"3c44cdddb6a900fa2b585dd299e03d12fa4293bc","crypto":{"cipher":"aes-128-ctr","ciphertext":"2cd6ab87086c47f343f2c4d957eace7986f3b3c87fc35a2aafbefb57a06d9f1c","cipherparams":{"iv":"4e16b6cd580866c1aa642fb4d7312c9b"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"0cabde93877f6e9a59070f9992f7a01848618263124835c90d4d07a0041fc57c"},"mac":"94b7776ea95b0ecd8406c7755acf17b389b7ebe489a8942e32082dfdc1f04f57"},"id":"ade1484b-a3bb-426f-9223-a1f5e3bde2e8","version":3}`,
	"90f79bf6eb2c4f870365e785982e1f101e93b906": `{"address":"90f79bf6eb2c4f870365e785982e1f101e93b906","crypto":{"cipher":"aes-128-ctr","ciphertext":"15144214d323871e00f7b205368128061c91b77a27b7deec935f8f5b734f0d42","cipherparams":{"iv":"bb22ba8051ef9f60abded7a9f4f2c6ae"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"07331ef3035567c00830b4e50d5dd68bc877974b4ce38cd42fef755be01556c9"},"mac":"2294eacadaf2761851814451d8c7dcca20a606a0344335d98f09403aba4e82ca"},"id":"96af8cc7-97e1-4bba-8968-632b034986c2","version":3}`,
	"15d34aaf54267db7d7c367839aaf71a00a2c6a65": `{"address":"15d34aaf54267db7d7c367839aaf71a00a2c6a65","crypto":{"cipher":"aes-128-ctr","ciphertext":"057878284a6c74d3ad99910adddd6b477b383837dbf2280efea585f0f0fdb012","cipherparams":{"iv":"e6eab29d60b526f305f8d47badf48687"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"dfdca8066d2486da5cb9a909d03744e2a8c6537930271e85e7cd8a5d952c0f22"},"mac":"f8352be41c9a06d69111ca4d8fcff0eef079b68b1173cad99803538991716c5d"},"id":"a35bb452-0d57-42d5-8d25-5a00a40a4db8","version":3}`,
	"9965507d1a55bcc2695c58ba16fb37d819b0a4dc": `{"address":"9965507d1a55bcc2695c58ba16fb37d819b0a4dc","crypto":{"cipher":"aes-128-ctr","ciphertext":"5a73201500307c6aa98edd44d962b344a893768331454a61595ec848e738e9d2","cipherparams":{"iv":"5282de2b3e2b305019a2fed5c62f3383"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"6ad001831d097f175fff7d6cf61301e9620b32afd9a7a6437e6030af14576a96"},"mac":"0a55eddbd13c713aa8b8c4106b2fb62bc1d1e18e7177207a444f83a4d8426ed5"},"id":"27aed2b2-cb94-4d37-8819-b15219187bb5","version":3}`,
	"976ea74026e726554db657fa54763abd0c3a0aa9": `{"address":"976ea74026e726554db657fa54763abd0c3a0aa9","crypto":{"cipher":"aes-128-ctr","ciphertext":"a6edf11e81b38e60a549696236cb9efc026e87adc45a9521ea7b2c45a2a9fbb9","cipherparams":{"iv":"82f4c79cd4b28a8585a9c78d758f832b"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"87400e16ecc320dadff85eccbf4dbaaea2dd91e50047e4aa391799bb319c1fd8"},"mac":"80c83dad05998db6c673a97096fcfad54636458f4a3c82483686b253f8cc9b69"},"id":"fc7d7694-6206-48fc-bb25-36b523f90df6","version":3}`,
	"14dc79964da2c08b23698b3d3cc7ca32193d9955": `{"address":"14dc79964da2c08b23698b3d3cc7ca32193d9955","crypto":{"cipher":"aes-128-ctr","ciphertext":"410f258bc8b12a0250cba22cbc5e413534fcf90bf322ced6943189ad9e43b4b9","cipherparams":{"iv":"1dd6077a8bee9b3bf2ca90e6abc8a237"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"5d3358bf99bbcb82354f40e5501abf4336bc141ee05d8feed4fbe7eb8c08c917"},"mac":"9cd959fa1e8129a8deb86e0264ec81d6cde79b5a19ae259b7d00543c9037908a"},"id":"689d7ad2-fe46-4c09-9c2a-a50e607989b8","version":3}`,
	"23618e81e3f5cdf7f54c3d65f7fbc0abf5b21e8f": `{"address":"23618e81e3f5cdf7f54c3d65f7fbc0abf5b21e8f","crypto":{"cipher":"aes-128-ctr","ciphertext":"13dccac740314edea20d44e6f3592575bbcb739ec5892d635326cff3c386eb86","cipherparams":{"iv":"bf42d811cd41fa97ddcae3425f8c3211"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"d2fa67cbb5e86d5bf9a90e27b8747bac493614b45778d43e9da1c14e06b2401d"},"mac":"7d2797cf344704d8f36265238d3938e06952c78ab7dfcbac53dc7f472c93d933"},"id":"4c8e899e-80f0-4417-9b1e-c5e29049f1e7","version":3}`,
	"a0ee7a142d267c1f36714e4a8f75612f20a79720": `{"address":"a0ee7a142d267c1f36714e4a8f75612f20a79720","crypto":{"cipher":"aes-128-ctr","ciphertext":"56bc8766f47aeafae74eea333e1e890a3776d7fae6c48cbdbffb270655ce050d","cipherparams":{"iv":"a66129e6a110b3ddf93b4355aa147c58"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"15c4e8bcc80920139eb236d91194825f1fce27dd2af281e0f2752d8a5dbc48bd"},"mac":"db01e720866ce8bb7897dfc7773e064003ad53429a79732ee769cf6d02273570"},"id":"87b2d76f-1b70-4e4f-8b2a-5d1915c1177c","version":3}`,
	"bcd4042de499d14e55001ccbb24a551f3b954096": `{"address":"bcd4042de499d14e55001ccbb24a551f3b954096","crypto":{"cipher":"aes-128-ctr","ciphertext":"e455eda6e38d246c03b930f845adfc8721ca75e9f47135cd4c18dbc3e5c5440a","cipherparams":{"iv":"0b1a0a24acc1ad25b0f170f751c2cb27"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"69f324ed0787794878bf5f84d4dbbc70dec1274cad666399edc48640605f64c8"},"mac":"f07da09c460a69f943f5639545d2b3f72c1e9789f0421ad41d3078ea3db12c96"},"id":"7ec7bb3c-c486-4785-a4fc-f8f4b2fc7764","version":3}`,
	"71be63f3384f5fb98995898a86b02fb2426c5788": `{"address":"71be63f3384f5fb98995898a86b02fb2426c5788","crypto":{"cipher":"aes-128-ctr","ciphertext":"4194377a05fd3d13e0a3155dad974a003fe5f7a3b5acb35d7d97c50daa8990d4","cipherparams":{"iv":"607670778baf62b1e86394cf1980487a"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"d63b890ad7f4fcc857681faabe9319dffc53893966ef0810bf64c4f319b0ffc5"},"mac":"bfaf924959e65c8030ece259d52ed52d5d21bd74f1a67ae545d4bb289a479e16"},"id":"0c6af842-384f-49b6-b5b7-199a1e05486b","version":3}`,
}

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
	// we need to deduplicate addresses to fund, because if present they will crash the genesis and here we cannot know
	// whether the caller hasn't passed duplicates both in the slice or in the keystore.
	accounts, _, err = slice.ValidateAndDeduplicateAddresses(accounts)
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
