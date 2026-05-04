package vaultjwtissuer

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestMintTokenSetsVaultSecretManagementClaim(t *testing.T) {
	client, err := NewClient("http://127.0.0.1:18123", "http://127.0.0.1:18123")
	require.NoError(t, err)

	tokenString, err := client.MintToken(TokenClaims{
		OrgID:         "org-1",
		WorkflowOwner: "0xabc123",
		RequestDigest: "digest-1",
	})
	require.NoError(t, err)

	parsed, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	require.NoError(t, err)

	claims, ok := parsed.Claims.(jwt.MapClaims)
	require.True(t, ok)
	require.Equal(t, "true", claims[ClaimVaultSecretManagementEnabled])
}
