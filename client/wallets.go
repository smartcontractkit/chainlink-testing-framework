package client

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	WalletFileLocation string = "../wallets.yml"
)

// GetEthWallets looks for wallets associated with a certain ETH based blockchain, looking for them in the order
// ENV var > config file > secrets manager
func GetEthWallets(name string) BlockchainWallets {
	// Check in env variables
	walletsFromEnv := os.Getenv(name)
	if walletsFromEnv != "" {
		return processEthWallets(walletsFromEnv)
	}

	// Check in the wallets file
	keyFile, err := ioutil.ReadFile(WalletFileLocation)
	if err != nil {
		log.Fatal(err)
	}

	// This might end up a bit dangerous? Maybe struct it?
	var config map[string]string
	err = yaml.Unmarshal(keyFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	if config[name] != "" {
		return processEthWallets(config[name])
	}

	// TODO Implement AWS or whatever secrets management we settle on

	return nil
}

// Processes ethereum private key strings into actual wallets
func processEthWallets(walletKeys string) BlockchainWallets {
	var processedWallets []BlockchainWallet
	splitKeys := strings.Split(walletKeys, ",")
	for _, key := range splitKeys {
		wallet, err := NewEthereumWallet(strings.TrimSpace(key))
		if err != nil {
			log.Fatal(err)
		}
		processedWallets = append(processedWallets, wallet)
	}
	return &Wallets{
		defaultWallet: 0,
		wallets:       processedWallets,
	}
}
