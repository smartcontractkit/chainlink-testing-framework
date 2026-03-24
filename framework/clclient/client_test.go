package clclient

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeAptosAccounts(t *testing.T) {
	keys := &AptosKeys{
		Data: []AptosKeyData{
			{
				Attributes: AptosKeyAttributes{
					Account: "0xa",
					Address: " 0xA ",
				},
			},
			{
				Attributes: AptosKeyAttributes{
					Account: "0xB",
				},
			},
			{
				Attributes: AptosKeyAttributes{
					Account: "not-an-account",
				},
			},
		},
	}

	accounts, err := normalizeAptosAccounts(keys)
	require.NoError(t, err)
	require.Equal(t, []string{
		"0x000000000000000000000000000000000000000000000000000000000000000a",
		"0x000000000000000000000000000000000000000000000000000000000000000b",
	}, accounts)
}

func TestNormalizeAptosAccounts_ErrorsOnEmpty(t *testing.T) {
	_, err := normalizeAptosAccounts(&AptosKeys{})
	require.Error(t, err)
}
