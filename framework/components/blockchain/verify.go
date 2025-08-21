package blockchain

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// VerifyContract wraps the forge verify-contract command.
func VerifyContract(out *Output, address, foundryDir, contractFile, contractName, compilerVersion string) error {
	args := []string{
		"verify-contract",
		"--rpc-url", out.Nodes[0].ExternalHTTPUrl,
		"--chain-id",
		out.ChainID,
		fmt.Sprintf("--compiler-version=%s", compilerVersion),
		address,
		fmt.Sprintf("%s:%s", contractFile, contractName),
		"--verifier", "blockscout",
		"--verifier-url", "http://localhost/api/",
	}
	return framework.RunCommandDir(foundryDir, "forge", args...)
}
