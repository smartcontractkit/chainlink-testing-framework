package examples

import (
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/stretchr/testify/require"
)

func TestPrivateJd(t *testing.T) {
	err := framework.DefaultNetwork(nil)
	require.NoError(t, err)
	_, err = jd.NewJD(&jd.Input{
		Image:            os.Getenv("CTF_JD_IMAGE"),
		CSAEncryptionKey: "d1093c0060d50a3c89c189b2e485da5a3ce57f3dcb38ab7e2c0d5f0bb2314a44", // random key for tests
	})
	require.NoError(t, err)
}
