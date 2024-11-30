package framework

import "fmt"

// VerifyContract wraps the forge verify-contract command.
func VerifyContract(rpcURL, address, contractFile, contractName, blockscoutURL string) error {
	verifierURL := fmt.Sprintf("%s/api/", blockscoutURL)
	args := []string{
		"verify-contract",
		"--rpc-url", rpcURL,
		"--compiler-version=0.8.17",
		address,
		fmt.Sprintf("%s:%s", contractFile, contractName),
		"--verifier", "blockscout",
		"--verifier-url", verifierURL,
	}
	return runCommand("forge", args...)
}
