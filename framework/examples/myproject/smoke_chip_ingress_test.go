package examples

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/go-github/v72/github"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	chipingressset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/chip_ingress_set"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

type ChipConfig struct {
	ChipIngress *chipingressset.Input `toml:"chip_ingress" validate:"required"`
}

// use config file: smoke_chip.toml
func TestChipIngressSmoke(t *testing.T) {
	// t.Skip("skipping smoke test until we have a way to fetch Chip Ingress image")
	os.Setenv("CTF_CONFIGS", "smoke_chip.toml")
	in, err := framework.Load[ChipConfig](t)
	require.NoError(t, err, "failed to load config")

	out, err := chipingressset.New(in.ChipIngress)
	require.NoError(t, err, "failed to create chip ingress set")

	t.Run("chainlink-protos can be registered", func(t *testing.T) {
		require.NotEmpty(t, out.ChipIngress.GRPCExternalURL, "GRPCExternalURL is not set")
		require.NotEmpty(t, out.RedPanda.SchemaRegistryExternalURL, "SchemaRegistryExternalURL is not set")

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		var client *github.Client
		if token := os.Getenv("GITHUB_TOKEN"); token != "" {
			ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
			tc := oauth2.NewClient(ctx, ts)
			client = github.NewClient(tc)
		} else {
			client = github.NewClient(nil)
		}

		createTopicsErr := chipingressset.CreateTopics(ctx, out.RedPanda.KafkaExternalURL, []string{"cre"})
		require.NoError(t, createTopicsErr, "failed to create topics")

		err := chipingressset.DefaultRegisterAndFetchProtos(ctx, client, []chipingressset.ProtoSchemaSet{
			{
				Owner:         "smartcontractkit",
				Repository:    "chainlink-protos",
				Ref:           "95decc005a91a1fd2621af9d9f00cb36d8061067",
				Folders:       []string{"workflows"},
				SubjectPrefix: "cre-",
			},
		}, out.RedPanda.SchemaRegistryExternalURL)
		require.NoError(t, err, "failed to register protos")
	})
}
