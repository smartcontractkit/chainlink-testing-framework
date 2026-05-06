package chip_ingress_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	chiprouter "github.com/smartcontractkit/chainlink-testing-framework/framework/components/chiprouter"
)

type ChipRouterConfig struct {
	ChipRouter *chiprouter.Input `toml:"chip_router" validate:"required"`
}

// use config file: smoke_chip_router.toml
func TestChipRouterDynamicPortsSmoke(t *testing.T) {
	os.Setenv("CTF_CONFIGS", "smoke_chip_router.toml")
	os.Setenv("CTF_CHIP_ROUTER_IMAGE", "local-cre-chip-router:v1.0.1")
	in, err := framework.Load[ChipRouterConfig](t)
	require.NoError(t, err, "failed to load config")
	require.NotEmpty(t, os.Getenv("CTF_CHIP_ROUTER_IMAGE"), "CTF_CHIP_ROUTER_IMAGE env var is not set")

	in.ChipRouter.GRPCPort = 0
	in.ChipRouter.AdminPort = 0

	out, err := chiprouter.NewWithContext(t.Context(), in.ChipRouter)
	require.NoError(t, err, "failed to create chip router")
	require.NotEmpty(t, out.ExternalGRPCURL, "GRPCExternalURL is not set")
	require.NotEmpty(t, out.ExternalAdminURL, "AdminExternalURL is not set")

	health, err := chiprouter.Health(t.Context(), out.ExternalAdminURL)
	require.NoError(t, err, "failed to get chip router health")
	require.NotEmpty(t, health, "health is not set")

}

// use config file: smoke_chip_router.toml
func TestChipRouterStaticPortsSmoke(t *testing.T) {
	os.Setenv("CTF_CONFIGS", "smoke_chip_router.toml")
	os.Setenv("CTF_CHIP_ROUTER_IMAGE", "local-cre-chip-router:v1.0.1")
	in, err := framework.Load[ChipRouterConfig](t)
	require.NoError(t, err, "failed to load config")
	require.NotEmpty(t, os.Getenv("CTF_CHIP_ROUTER_IMAGE"), "CTF_CHIP_ROUTER_IMAGE env var is not set")
	in.ChipRouter.GRPCPort = 7197
	in.ChipRouter.AdminPort = 7198

	out, err := chiprouter.NewWithContext(t.Context(), in.ChipRouter)
	require.NoError(t, err, "failed to create chip router")
	require.NotEmpty(t, out.ExternalGRPCURL, "GRPCExternalURL is not set")
	require.NotEmpty(t, out.ExternalAdminURL, "AdminExternalURL is not set")

	_, err = chiprouter.Health(t.Context(), out.ExternalAdminURL)
	require.NoError(t, err, "failed to get chip router health")
}
