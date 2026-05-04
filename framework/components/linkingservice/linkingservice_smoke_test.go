package linkingservice_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	linkingpb "github.com/smartcontractkit/chainlink-protos/linking-service/go/v1"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/linkingservice"
)

func TestLinkingServiceSmoke(t *testing.T) {
	image := os.Getenv(linkingservice.ImageEnvVar)
	if image == "" {
		t.Skipf("%s env var is not set", linkingservice.ImageEnvVar)
	}

	t.Cleanup(func() {
		_ = framework.RemoveTestContainers()
	})

	out, err := linkingservice.NewWithContext(t.Context(), &linkingservice.Input{
		Image:         image,
		ContainerName: framework.DefaultTCName("linking-service-smoke"),
	})
	require.NoError(t, err)

	admin := linkingservice.NewAdminClientFromOutput(out)
	require.NotNil(t, admin)
	require.NoError(t, admin.SetOwnerOrg(t.Context(), "0xabc123", "org-1"))

	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(out.LocalGRPCURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	resp, err := linkingpb.NewLinkingServiceClient(conn).GetOrganizationFromWorkflowOwner(ctx, &linkingpb.GetOrganizationFromWorkflowOwnerRequest{
		WorkflowOwner: "0xabc123",
	})
	require.NoError(t, err)
	require.Equal(t, "org-1", resp.GetOrganizationId())
}
