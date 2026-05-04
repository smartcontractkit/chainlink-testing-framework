package vaultjwtissuer_test

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/vaultjwtissuer"
)

func TestVaultJWTIssuerSmoke(t *testing.T) {
	image := os.Getenv(vaultjwtissuer.ImageEnvVar)
	if image == "" {
		t.Skipf("%s env var is not set", vaultjwtissuer.ImageEnvVar)
	}

	t.Cleanup(func() {
		_ = framework.RemoveTestContainers()
	})

	out, err := vaultjwtissuer.NewWithContext(t.Context(), &vaultjwtissuer.Input{
		Image:         image,
		ContainerName: framework.DefaultTCName("vault-jwt-issuer-smoke"),
	})
	require.NoError(t, err)

	resp, err := http.Get(strings.TrimRight(out.LocalHTTPURL, "/") + "/.well-known/openid-configuration") //nolint:noctx // smoke test reads a local admin endpoint
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var openID struct {
		Issuer  string `json:"issuer"`
		JWKSURI string `json:"jwks_uri"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&openID))
	require.Equal(t, vaultjwtissuer.NormalizeIssuerURL(out.LocalHTTPURL), openID.Issuer)
	require.Equal(t, strings.TrimRight(out.LocalHTTPURL, "/")+"/.well-known/jwks.json", openID.JWKSURI)

	client, err := vaultjwtissuer.NewClientFromOutput(out)
	require.NoError(t, err)

	token, err := client.MintToken(vaultjwtissuer.TokenClaims{
		OrgID:         "org-1",
		WorkflowOwner: "0xabc123",
		RequestDigest: "digest-1",
	})
	require.NoError(t, err)
	require.Len(t, strings.Split(token, "."), 3)
}
