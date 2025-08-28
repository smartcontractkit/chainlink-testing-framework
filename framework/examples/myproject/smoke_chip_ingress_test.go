package examples

import (
	"context"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	chipingressset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/chip_ingress_set"
	"github.com/stretchr/testify/require"
)

type ChipConfig struct {
	ChipIngress *chipingressset.Input `toml:"chip_ingress" validate:"required"`
}

// use config file: smoke_chip.toml
func TestChipIngressSmoke(t *testing.T) {
	in, err := framework.Load[ChipConfig](t)
	require.NoError(t, err, "failed to load config")

	out, err := chipingressset.New(in.ChipIngress)
	require.NoError(t, err, "failed to create chip ingress set")
	require.NotEmpty(t, out.ChipIngress.GRPCExternalURL, "GRPCExternalURL is not set")
	require.NotEmpty(t, out.RedPanda.SchemaRegistryExternalURL, "SchemaRegistryExternalURL is not set")

	t.Run("remote chainlink-protos can be registered", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		createTopicsErr := chipingressset.CreateTopics(ctx, out.RedPanda.KafkaExternalURL, []string{"cre"})
		require.NoError(t, createTopicsErr, "failed to create topics")

		err := chipingressset.DefaultRegisterAndFetchProtos(ctx, nil, []chipingressset.ProtoSchemaSet{
			{
				URI:           "https://github.com/smartcontractkit/chainlink-protos",
				Ref:           "95decc005a91a1fd2621af9d9f00cb36d8061067",
				Folders:       []string{"workflows"},
				SubjectPrefix: "cre-",
			},
		}, out.RedPanda.SchemaRegistryExternalURL)
		require.NoError(t, err, "failed to register protos")
	})

	t.Run("local protos can be registered", func(t *testing.T) {
		t.Skip("we can only one run of these nested at a time, because they register the same protos")
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		createTopicsErr := chipingressset.CreateTopics(ctx, out.RedPanda.KafkaExternalURL, []string{"cre"})
		require.NoError(t, createTopicsErr, "failed to create topics")

		err := chipingressset.DefaultRegisterAndFetchProtos(ctx, nil, []chipingressset.ProtoSchemaSet{
			{
				URI:           "file://../../../../chainlink-protos", // works also with absolute path
				Folders:       []string{"workflows"},
				SubjectPrefix: "cre-",
			},
		}, out.RedPanda.SchemaRegistryExternalURL)
		require.NoError(t, err, "failed to register protos")
	})
}
