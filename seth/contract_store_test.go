package seth_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

func TestSmokeContractABIStore(t *testing.T) {
	type test struct {
		name              string
		abiPath           string
		gethWrappersPaths []string
		err               string
		expectedABICount  int
	}

	tests := []test{
		{
			name:             "can load the ABI only from ABI files",
			abiPath:          "./contracts/abi",
			expectedABICount: 8,
		},
		{
			name:              "can load the ABI from ABI files and from gethwrappers",
			abiPath:           "./contracts/abi",
			gethWrappersPaths: []string{"./contracts/bind"},
			expectedABICount:  11,
		},
		{
			name:              "can load the ABI only from gethwrappers",
			gethWrappersPaths: []string{"./contracts/bind"},
			expectedABICount:  7,
		},
		{
			name:              "can load the ABI from 2 gethwrappers folders",
			gethWrappersPaths: []string{"./contracts/bind", "./contracts/bind2"},
			expectedABICount:  8,
		},
		{
			name:    "can't open the ABI path",
			abiPath: "dasdsadd",
			err:     "open dasdsadd: no such file or directory",
		},
		{
			name:              "can't open the gethwrappers path",
			gethWrappersPaths: []string{"dasdsadd"},
			err:               "failed to load geth wrappers from [dasdsadd]: lstat dasdsadd: no such file or directory",
		},
		{
			name:              "correct and broken gethwrappers path",
			gethWrappersPaths: []string{"./contracts/emptyContractDir", "dasdsadd"},
			err:               "failed to load geth wrappers from [./contracts/emptyContractDir dasdsadd]: lstat dasdsadd: no such file or directory",
		},
		{
			name:              "broken and correct gethwrappers path",
			gethWrappersPaths: []string{"dasdsadd", "./contracts/emptyContractDir"},
			err:               "failed to load geth wrappers from [dasdsadd ./contracts/emptyContractDir]: lstat dasdsadd: no such file or directory",
		},
		{
			name:    "empty ABI dir",
			abiPath: "./contracts/emptyContractDir",
			err:     "no ABI files (*.abi) found in directory './contracts/emptyContractDir'.",
		},
		{
			name:              "empty gethwrappers dir",
			gethWrappersPaths: []string{"./contracts/emptyContractDir"},
			err:               "failed to load geth wrappers from [./contracts/emptyContractDir]: no geth wrapper files found in directories: [./contracts/emptyContractDir]",
		},
		{
			name:              "empty ABI and gethwrappers dir",
			abiPath:           "./contracts/emptyContractDir",
			gethWrappersPaths: []string{"./contracts/emptyContractDir"},
			err:               "no ABI files (*.abi) found in directory './contracts/emptyContractDir'.",
		},
		{
			name:              "no MetaData in one of gethwrappers",
			gethWrappersPaths: []string{"./contracts/noMetaDataContractDir"},
			expectedABICount:  1,
		},
		{
			name:              "empty MetaData in one of gethwrappers",
			gethWrappersPaths: []string{"./contracts/emptyMetaDataContractDir"},
			err:               "failed to load geth wrappers from [./contracts/emptyMetaDataContractDir]: failed to parse ABI content from 'contracts/emptyMetaDataContractDir/NetworkDebugContract_Broken.go': EOF",
		},
		{
			name:              "gethwrappers dir mixes regular go files and gethwrappers",
			gethWrappersPaths: []string{"./contracts/gethWrapperAndGoFile"},
			expectedABICount:  1,
		},
		{
			name:    "invalid ABI inside ABI dir",
			abiPath: "./contracts/invalidContractDir",
			err:     "failed to parse ABI file 'NetworkDebugContract.abi': invalid character ':' after array element",
		},
		{
			name:              "invalid ABI in gethwrappers inside dir",
			gethWrappersPaths: []string{"./contracts/invalidContractDir"},
			err:               "failed to load geth wrappers from [./contracts/invalidContractDir]: failed to parse ABI content from 'contracts/invalidContractDir/NetworkDebugContract_Broken.go': invalid character 'i' looking for beginning of value",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			cs, err := seth.NewContractStore(tc.abiPath, "", tc.gethWrappersPaths)
			if err == nil {
				require.NotNil(t, cs.ABIs, "ABIs should not be nil")
				require.Equal(t, tc.expectedABICount, len(cs.ABIs), "ABIs should have the expected count")
				require.NotNil(t, cs.BINs, "BINs should not be nil")
				require.Equal(t, make(map[string][]uint8), cs.BINs)
				err = errors.New("")
			}
			require.Contains(t, err.Error(), tc.err, "error should match")
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
			name:     "can load the ABI and BIN and gethWrappers",
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
			err:     "no BIN files (*.bin) found in directory './contracts/emptyContractDir'.",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			cs, err := seth.NewContractStore(tc.abiPath, tc.binPath, nil)
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
			require.Contains(t, err.Error(), tc.err, "error should match")
		})
	}
}
