package blockchain

import (
	"bytes"
	"strings"
	"testing"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/stretchr/testify/require"
)

func TestParseSuiKeytoolGenerateJSON(t *testing.T) {
	t.Parallel()

	const addr = "0xabc"
	compact := `{"alias":null,"flag":0,"keyScheme":"ed25519","mnemonic":"a b c","peerId":"p","publicBase64Key":"k","suiAddress":"` + addr + `"}`

	t.Run("compact one-line JSON", func(t *testing.T) {
		t.Parallel()
		got, err := parseSuiKeytoolGenerateJSON(compact)
		require.NoError(t, err)
		require.Equal(t, addr, got.SuiAddress)
	})

	t.Run("preamble before JSON", func(t *testing.T) {
		t.Parallel()
		in := "some log line\n" + compact
		got, err := parseSuiKeytoolGenerateJSON(in)
		require.NoError(t, err)
		require.Equal(t, addr, got.SuiAddress)
	})

	t.Run("legacy newline after brace (old parser shape)", func(t *testing.T) {
		t.Parallel()
		legacy := "{\n  \"suiAddress\": \"" + addr + "\"\n}"
		got, err := parseSuiKeytoolGenerateJSON(legacy)
		require.NoError(t, err)
		require.Equal(t, addr, got.SuiAddress)
	})

	t.Run("docker multiplexed stdout", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		w := stdcopy.NewStdWriter(&buf, stdcopy.Stdout)
		_, err := w.Write([]byte(compact))
		require.NoError(t, err)
		got, err := parseSuiKeytoolGenerateJSON(buf.String())
		require.NoError(t, err)
		require.Equal(t, addr, got.SuiAddress)
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()
		_, err := parseSuiKeytoolGenerateJSON("no json here")
		require.Error(t, err)
		require.True(t, strings.Contains(err.Error(), "failed to parse"))
	})
}
