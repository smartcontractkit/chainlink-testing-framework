package cre_storage_service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnv(t *testing.T) {
	in := &Input{
		RegistryRPCURL:        "http://anvil-1337:8545",
		RegistryChainSelector: 3379446385462418246,
		RegistryAddress:       "0xabc",
		DonFamily:             "test-don-family",
	}
	vars := env(in, []string{"CRE_STORAGE_IMAGE=custom:tag", "CRE_LOCAL_DON_FAMILY=from-host", "MALFORMED"}, []int{1, 2, 3})

	require.Equal(t, "custom:tag", vars["CRE_STORAGE_IMAGE"])
	require.Equal(t, "test-don-family", vars["CRE_LOCAL_DON_FAMILY"])
	require.Equal(t, "ctf", vars["CRE_NETWORK"])
	require.Equal(t, "true", vars["CRE_NETWORK_EXTERNAL"])
	require.Equal(t, "http://anvil-1337:8545", vars["CRE_LOCAL_REGISTRY_RPC_URL"])
	require.Equal(t, "3379446385462418246", vars["CRE_LOCAL_REGISTRY_CHAIN_SELECTOR"])
	require.Equal(t, "0xabc", vars["CRE_LOCAL_REGISTRY_ADDRESS"])
	require.Equal(t, "1", vars["USER_SERVICE_HOST_PORT"])
	require.Equal(t, "2", vars["NODE_SERVICE_HOST_PORT"])
	require.Equal(t, "3", vars["S3_HOST_PORT"])
	require.NotContains(t, vars, "MALFORMED")

	in.Network = "custom"
	require.Equal(t, "custom", env(in, nil, []int{1, 2, 3})["CRE_NETWORK"])
}

func TestNewRejectsIncompleteInput(t *testing.T) {
	_, err := New(&Input{ComposeFile: "compose.yaml", RegistryChainSelector: 1})
	require.ErrorContains(t, err, "missing registry_rpc_url, registry_address, don_family")

	_, err = New(nil)
	require.ErrorContains(t, err, "input is nil")
}
