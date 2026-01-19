package examples

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/fullstorydev/grpcurl"
	"github.com/go-resty/resty/v2"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

type CfgCanton struct {
	BlockchainA *blockchain.Input `toml:"blockchain_a" validate:"required"`
}

func TestCantonSmoke(t *testing.T) {
	in, err := framework.Load[CfgCanton](t)
	require.NoError(t, err)

	bc, err := blockchain.NewBlockchainNetwork(in.BlockchainA)
	require.NoError(t, err)

	t.Run("Test scan endpoint", func(t *testing.T) {
		resp, err := resty.New().SetBaseURL(bc.NetworkSpecificData.CantonEndpoints.ScanAPIURL).R().
			Get("/v0/dso-party-id")
		assert.NoError(t, err)
		fmt.Println(resp)
	})
	t.Run("Test registry endpoint", func(t *testing.T) {
		resp, err := resty.New().SetBaseURL(bc.NetworkSpecificData.CantonEndpoints.RegistryAPIURL).R().
			Get("/metadata/v1/instruments")
		assert.NoError(t, err)
		fmt.Println(resp)
	})

	testParticipant := func(t *testing.T, name string, endpoints blockchain.CantonParticipantEndpoints) {
		t.Run(fmt.Sprintf("Test %s endpoints", name), func(t *testing.T) {
			require.NoError(t, err)

			// JSON Ledger API
			fmt.Println("Calling JSON Ledger API")
			resp, err := resty.New().SetBaseURL(endpoints.JSONLedgerAPIURL).SetAuthToken(endpoints.JWT).R().
				Get("/v2/packages")
			assert.NoError(t, err)
			fmt.Println(resp)

			// gRPC Ledger API - use reflection
			fmt.Println("Calling gRPC Ledger API")
			res, err := callGRPC(t.Context(), endpoints.GRPCLedgerAPIURL, "com.daml.ledger.api.v2.admin.PartyManagementService/GetParties", `{}`, []string{fmt.Sprintf("Authorization: Bearer %s", endpoints.JWT)})
			assert.NoError(t, err)
			fmt.Println(res)

			// gRPC Admin API - use reflection
			fmt.Println("Calling gRPC Admin API")
			res, err = callGRPC(t.Context(), endpoints.AdminAPIURL, "com.digitalasset.canton.admin.participant.v30.PackageService/ListDars", `{}`, []string{fmt.Sprintf("Authorization: Bearer %s", endpoints.JWT)})
			assert.NoError(t, err)
			fmt.Println(res)

			// Validator API
			fmt.Println("Calling Validator API")
			resp, err = resty.New().SetBaseURL(endpoints.ValidatorAPIURL).SetAuthToken(endpoints.JWT).R().
				Get("/v0/admin/users")
			assert.NoError(t, err)
			fmt.Println(resp)

			// HTTP Health Check
			fmt.Println("Calling HTTP Health Check")
			resp, err = resty.New().SetBaseURL(endpoints.HTTPHealthCheckURL).R().
				Get("/health")
			assert.NoError(t, err)
			fmt.Println(resp)

			// gRPC Health Check
			fmt.Println("Calling gRPC Health Check")
			res, err = callGRPC(t.Context(), endpoints.GRPCHealthCheckURL, "grpc.health.v1.Health/Check", `{}`, nil)
			assert.NoError(t, err)
			fmt.Println(res)
		})
	}

	// Call all participants, starting with the SV
	testParticipant(t, "sv", bc.NetworkSpecificData.CantonEndpoints.SuperValidator)
	for i := 1; i <= in.BlockchainA.NumberOfCantonValidators; i++ {
		testParticipant(t, fmt.Sprintf("participant%d", i), bc.NetworkSpecificData.CantonEndpoints.Participants[i-1])
	}
}

// callGRPC makes a gRPC call to the given URL and method with the provided JSON request and headers.
func callGRPC(ctx context.Context, url string, method string, jsonRequest string, headers []string) (string, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", fmt.Errorf("failed to create grpc client: %w", err)
	}
	defer conn.Close()

	options := grpcurl.FormatOptions{EmitJSONDefaultFields: true}
	jsonRequestReader := strings.NewReader(jsonRequest)
	var output bytes.Buffer

	reflectClient := grpcreflect.NewClientAuto(ctx, conn)
	defer reflectClient.Reset()
	descriptorSource := grpcurl.DescriptorSourceFromServer(ctx, reflectClient)

	requestParser, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, descriptorSource, jsonRequestReader, options)
	if err != nil {
		return "", fmt.Errorf("failed to create request parser and formatter: %w", err)
	}
	eventHandler := &grpcurl.DefaultEventHandler{
		Out:            &output,
		Formatter:      formatter,
		VerbosityLevel: 0,
	}

	err = grpcurl.InvokeRPC(ctx, descriptorSource, conn, method, headers, eventHandler, requestParser.Next)
	if err != nil {
		return "", fmt.Errorf("rpc call failed: %w", err)
	}
	return output.String(), nil
}
