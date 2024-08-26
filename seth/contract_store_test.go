package seth_test

import (
	"testing"

	"github.com/pkg/errors"
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
			expectedABICount: 5,
		},
		{
			name:              "can load the ABI from ABI files and from gethwrappers",
			abiPath:           "./contracts/abi",
			gethWrappersPaths: []string{"./contracts/bind"},
			expectedABICount:  10,
		},
		{
			name:              "can load the ABI only from gethwrappers",
			gethWrappersPaths: []string{"./contracts/bind"},
			expectedABICount:  5,
		},
		{
			name:              "can load the ABI from 2 gethwrappers folders",
			gethWrappersPaths: []string{"./contracts/bind", "./contracts/bind2"},
			expectedABICount:  6,
		},
		{
			name:    "can't open the ABI path",
			abiPath: "dasdsadd",
			err:     "open dasdsadd: no such file or directory",
		},
		{
			name:              "can't open the gethwrappers path",
			gethWrappersPaths: []string{"dasdsadd"},
			err:               "lstat dasdsadd: no such file or directory",
		},
		{
			name:              "correct and broken gethwrappers path",
			gethWrappersPaths: []string{"./contracts/emptyContractDir", "dasdsadd"},
			err:               "lstat dasdsadd: no such file or directory",
		},
		{
			name:              "broken and correct gethwrappers path",
			gethWrappersPaths: []string{"dasdsadd", "./contracts/emptyContractDir"},
			err:               "lstat dasdsadd: no such file or directory",
		},
		{
			name:    "empty ABI dir",
			abiPath: "./contracts/emptyContractDir",
		},
		{
			name:              "empty gethwrappers dir",
			gethWrappersPaths: []string{"./contracts/emptyContractDir"},
		},
		{
			name:              "empty ABI and gethwrappers dir",
			abiPath:           "./contracts/emptyContractDir",
			gethWrappersPaths: []string{"./contracts/emptyContractDir"},
		},
		{
			name:              "no MetaData in in gethwrapper",
			gethWrappersPaths: []string{"./contracts/noMetaDataContractDir"},
		},
		{
			name:              "empty MetaData in in gethwrapper",
			gethWrappersPaths: []string{"./contracts/emptyMetaDataContractDir"},
			err:               "failed to parse ABI content: EOF",
		},
		{
			name:              "gethwrappers dir mixes regular go files and gethwrappers",
			gethWrappersPaths: []string{"./contracts/gethWrapperAndGoFile"},
			expectedABICount:  1,
		},
		{
			name:    "invalid ABI inside ABI dir",
			abiPath: "./contracts/invalidContractDir",
			err:     "failed to parse ABI file: invalid character ':' after array element",
		},
		{
			name:              "invalid ABI in gethwrappers inside dir",
			gethWrappersPaths: []string{"./contracts/invalidContractDir"},
			err:               "failed to parse ABI content: invalid character 'i' looking for beginning of value",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			cs, err := seth.NewContractStore(tc.abiPath, tc.abiPath, tc.gethWrappersPaths)
			if err == nil {
				require.NotNil(t, cs.ABIs, "ABIs should not be nil")
				require.Equal(t, tc.expectedABICount, len(cs.ABIs), "ABIs should have the expected count")
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
			require.Equal(t, tc.err, err.Error(), "error should match")
		})
	}
}
