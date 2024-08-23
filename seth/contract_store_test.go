package seth_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/seth"
	"github.com/stretchr/testify/require"
)

func TestSmokeContractABIStore(t *testing.T) {

	type test struct {
		name    string
		abiPath string
		err     string
	}

	tests := []test{
		{
			name:    "can load the ABI",
			abiPath: "./contracts/abi",
		},
		{
			name:    "can't open the ABI path",
			abiPath: "dasdsadd",
			err:     "open dasdsadd: no such file or directory",
		},
		{
			name:    "empty ABI dir",
			abiPath: "./contracts/emptyContractDir",
		},
		{
			name:    "invalid ABI inside dir",
			abiPath: "./contracts/invalidContractDir",
			err:     "failed to parse ABI file: invalid character ':' after array element",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			cs, err := seth.NewContractStore(tc.abiPath, tc.abiPath)
			if err == nil {
				require.NotNil(t, cs.ABIs, "ABIs should not be nil")
				require.NotNil(t, cs.BINs, "BINs should not be nil")
				require.Equal(t, make(map[string][]uint8), cs.BINs)
				err = errors.New("")
			}
			require.Equal(t, tc.err, err.Error())
		})
	}
}

func TestSmokeContractBINStore(t *testing.T) {

	type test struct {
		name     string
		abiPath  string
		binPath  string
		binFound bool
		err      string
	}

	tests := []test{
		{
			name:     "can load the ABI and BIN",
			abiPath:  "./contracts/abi",
			binPath:  "./contracts/bin",
			binFound: true,
		},
		{
			name:    "can't open the BIN path",
			abiPath: "./contracts/abi",
			binPath: "./contract/i-don't-exist",
			err:     "open ./contract/i-don't-exist: no such file or directory",
		},
		{
			name:    "empty BIN dir",
			abiPath: "./contracts/abi",
			binPath: "./contracts/emptyContractDir",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			cs, err := seth.NewContractStore(tc.abiPath, tc.binPath)
			if err == nil {
				require.NotEmpty(t, cs.ABIs, "ABIs should not be empty")
				err = errors.New("")
				if tc.binFound {
					require.NotEmpty(t, cs.BINs, "BINs should not be empty")
				} else {
					require.Empty(t, cs.BINs, "BINs should be empty")
				}
			} else {
				require.Nil(t, cs, "ContractStore should be nil")
			}
			require.Equal(t, tc.err, err.Error(), "error should match")
		})
	}
}
