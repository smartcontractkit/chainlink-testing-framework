package framework

/*
This client helps us to store and load secrets in AWS Secrets Manager
It also prevents secrets from being printed by mistake
*/

import (
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var (
	redacted     = "***"
	redactedJSON = "{\"key\":\"is_redacted\"}"
	redactedTOML = "[Key]\n\nis = redacted"
)

// AWSSecret is a wrapper preventing accidental printing or marshalling
type AWSSecret string

// Value is used to return masked secret value
func (s AWSSecret) Value() string { return string(s) }

// The String method is used to print values passed as an operand
// to any format that accepts a string or to an unformatted printer
// such as Print.
func (s AWSSecret) String() string { return redacted }

// The GoString method is used to print values passed as an operand
// to a %#v format.
func (s AWSSecret) GoString() string { return redacted }

// MarshalText encodes the receiver into UTF-8-encoded text and returns the result.
func (s AWSSecret) MarshalText() ([]byte, error) { return []byte(redactedTOML), nil }

// MarshalJSON Marshaler is the interface implemented by types that
// can marshal themselves into valid JSON.
func (s AWSSecret) MarshalJSON() ([]byte, error) { return []byte(redactedJSON), nil }

var (
	_ fmt.Stringer           = (*AWSSecret)(nil)
	_ fmt.GoStringer         = (*AWSSecret)(nil)
	_ encoding.TextMarshaler = (*AWSSecret)(nil)
	_ json.Marshaler         = (*AWSSecret)(nil)
)

// AWSSecretsManager is an AWS Secrets Manager service wrapper
type AWSSecretsManager struct {
	Client         *secretsmanager.Client
	RequestTimeout time.Duration
	l              zerolog.Logger
}

// NewAWSSecretsManager create a new connection to AWS Secrets Manager
func NewAWSSecretsManager(requestTimeout time.Duration) (*AWSSecretsManager, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return nil, fmt.Errorf("region is required for AWSSecretsManager, use env variable: export AWS_REGION=...: %w", err)
	}
	cfg.Region = region
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config, %v", err)
	}
	l := log.Logger.With().Str("Component", "AWSSecretsManager").Logger()
	l.Info().Msg("Connecting to AWS Secrets Manager")
	return &AWSSecretsManager{
		Client:         secretsmanager.NewFromConfig(cfg),
		RequestTimeout: requestTimeout,
		l:              l,
	}, nil
}

// CreateSecret creates a specific secret by key
func (sm *AWSSecretsManager) CreateSecret(key string, val string, override bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), sm.RequestTimeout)
	defer cancel()

	sm.l.Debug().Str("Key", key).Msg("Creating secret by key")
	k := &key
	v := &val
	_, err := sm.Client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:                        k,
		SecretString:                v,
		ForceOverwriteReplicaSecret: override,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create a secret by key")
	}
	return nil
}

// GetSecret gets a specific secret by key
func (sm *AWSSecretsManager) GetSecret(key string) (AWSSecret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), sm.RequestTimeout)
	defer cancel()

	sm.l.Debug().Str("Key", key).Msg("Reading secret by key")
	k := &key
	out, err := sm.Client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: k,
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to read a secret by key")
	}
	return AWSSecret(*out.SecretString), nil
}

// RemoveSecret removes a specific secret by key
func (sm *AWSSecretsManager) RemoveSecret(key string, noRecovery bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), sm.RequestTimeout)
	defer cancel()

	sm.l.Debug().Str("Key", key).Msg("Removing secret by key")
	k := &key
	b := &noRecovery
	_, err := sm.Client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
		SecretId:                   k,
		ForceDeleteWithoutRecovery: b,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to remove a secret by key")
	}
	return nil
}
