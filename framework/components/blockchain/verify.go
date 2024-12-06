package blockchain

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// VerifyContract wraps the forge verify-contract command.
func VerifyContract(out *Output, address, foundryDir, contractFile, contractName string) error {
	args := []string{
		"verify-contract",
		"--rpc-url", out.Nodes[0].HostHTTPUrl,
		"--chain-id",
		out.ChainID,
		"--compiler-version=0.8.24",
		address,
		fmt.Sprintf("%s:%s", contractFile, contractName),
		"--verifier", "blockscout",
		"--verifier-url", "http://localhost/api/",
	}
	return framework.RunCommandDir(foundryDir, "forge", args...)
}
