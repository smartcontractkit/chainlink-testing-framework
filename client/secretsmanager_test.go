package client

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"
)

func TestSecretsPrintMarshal(t *testing.T) {
	s := AWSSecret("1")
	t.Run("print the whole struct as string", func(t *testing.T) {
		//nolint
		output := fmt.Sprintf("%s", s)
		require.Equal(t, redacted, output)
	})
	t.Run("print the whole struct as +v%", func(t *testing.T) {
		output := fmt.Sprintf("%+v", s)
		require.Equal(t, redacted, output)
	})
	t.Run("spew library should not work either", func(t *testing.T) {
		output := spew.Sdump(s)
		require.Equal(t, "(client.AWSSecret) (len=1) ***", output[:len(output)-1])
	})
	t.Run("marshal the whole struct as JSON", func(t *testing.T) {
		d, err := json.Marshal(s)
		require.NoError(t, err)
		require.Equal(t, redactedJSON, string(d))
	})
	t.Run("marshal the whole struct as TOML results in err", func(t *testing.T) {
		// github.com/pelletier/go-toml/v2 since version 2 does not allow any struct to implement MarshalText()
		// so it results in error if we try to marshal secrets
		_, err := toml.Marshal(s)
		require.Error(t, err)
	})
}

func TestManualSecretsCRUD(t *testing.T) {
	t.Skip("Need AWS role to be enabled")
	// fill .envrc with AWS auth values and run manually
	// export AWS_REGION="us-west-2"
	// export AWS_ACCESS_KEY_ID=
	// export AWS_SECRET_ACCESS_KEY=
	// export AWS_SESSION_TOKEN=
	sm, err := NewAWSSecretsManager(1 * time.Minute)
	require.NoError(t, err)

	t.Run("basic single value CRUD", func(t *testing.T) {
		k := uuid.NewString()
		v := uuid.NewString()
		t.Cleanup(func() {
			err = sm.RemoveSecret(k, true)
			require.NoError(t, err)
		})
		err = sm.CreateSecret(k, v, true)
		require.NoError(t, err)
		secret, err := sm.GetSecret(k)
		require.NoError(t, err)
		require.Equal(t, secret, AWSSecret(v))
		require.Equal(t, v, AWSSecret(v).Value())
	})
}
