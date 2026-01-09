# Canton Blockchain Client

This supports spinning up a Canton LocalNet instance. It is heavily based on
the [Splice LocalNet Docker Compose setup](https://github.com/hyperledger-labs/splice/blob/3aede18a641bb657e25eea240adfb869d5c12503/cluster/compose/localnet/compose.yaml).

The LocalNet consists of one Super Validator and a variable number of additional participants/validators that can
be configured using the `number_of_canton_validators` parameter.

## Configuration

```toml
[blockchain_a]
  type = "canton"
  number_of_canton_validators = 5 # Controls the number of validators in the LocalNet
  image = "0.5.3"                 # Optional, can be used to override the default Canton image tag
  port = "8088"                   # Optional, defaults to 8080
```

## Endpoints

A reverse proxy is set up to route requests to the appropriate Canton node based on the URL path.

| Endpoint Path                                             | Description         | Documentation                                                                  |
|-----------------------------------------------------------|---------------------|--------------------------------------------------------------------------------|
| `http://scan.localhost:[PORT]/api/scan`                   | Scan API            | https://docs.sync.global/app_dev/scan_api/index.html                           |
| `http://scan.localhost:[PORT]/registry`                   | Token Standard APIs | https://docs.sync.global/app_dev/token_standard/index.html#api-references      |
|                                                           |                     |                                                                                |
| `http://[PARTICIPANT].json-ledger-api.localhost:[PORT]`   | JSON Ledger API     | https://docs.digitalasset.com/build/3.3/reference/json-api/json-api.html       |
| `grpc://[PARTICIPANT].grpc-ledger-api.localhost:[PORT]`   | gRPC Ledger API     | https://docs.digitalasset.com/build/3.3/reference/lapi-proto-docs.html         |
| `grpc://[PARTICIPANT].admin-api.localhost:[PORT]`         | gRPC Admin API      | https://docs.digitalasset.com/operate/3.5/howtos/configure/apis/admin_api.html |
| `http://[PARTICIPANT].validator-api.localhost:[PORT]`     | Validator API       | https://docs.sync.global/app_dev/validator_api/index.html                      |
| `http://[PARTICIPANT].http-health-check.localhost:[PORT]` | HTTP Health Check   | responds on GET /health                                                        |
| `grpc://[PARTICIPANT].grpc-health-check.localhost:[PORT]` | gRPC Health Check   | https://grpc.io/docs/guides/health-checking/                                   |

To access a participant's endpoint, replace `[PARTICIPANT]` with the participant's name (e.g., `sv`, `participant1`,
`participant2`, etc.).

> [!NOTE]
> The maximum number of participants is 99.

## Authentication

The following endpoints require authentication:

- JSON Ledger API
- gRPC Ledger API
- gRPC Admin API
- Validator API

To authenticate, create a JWT bearer using the following claims:

- `aud`: Set to the exported const `AuthProviderAudience`
- `sub`: Set to `user-[PARTICIPANT]` replacing `[PARTICIPANT]` with the participant's name like `sv`, `participant1`,
  etc.

Sign the JWT using the HMAC SHA256 algorithm with the secret of the exported const `AuthProviderSecret`.

## Usage

```golang
package examples

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/fullstorydev/grpcurl"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain/canton"
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
			j, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
				Issuer:    "",
				Subject:   fmt.Sprintf("user-%s", name),
				Audience:  []string{canton.AuthProviderAudience},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ID:        "",
			}).SignedString([]byte("unsafe"))

			// JSON Ledger API
			fmt.Println("Calling JSON Ledger API")
			resp, err := resty.New().SetBaseURL(endpoints.JSONLedgerAPIURL).SetAuthToken(j).R().
				Get("/v2/packages")
			assert.NoError(t, err)
			fmt.Println(resp)

			// gRPC Ledger API - use reflection
			fmt.Println("Calling gRPC Ledger API")
			res, err := callGRPC(t.Context(), endpoints.GRPCLedgerAPIURL, "com.daml.ledger.api.v2.admin.PartyManagementService/GetParties", `{}`, []string{fmt.Sprintf("Authorization: Bearer %s", j)})
			assert.NoError(t, err)
			fmt.Println(res)

			// gRPC Admin API - use reflection
			fmt.Println("Calling gRPC Admin API")
			res, err = callGRPC(t.Context(), endpoints.AdminAPIURL, "com.digitalasset.canton.admin.participant.v30.PackageService/ListDars", `{}`, []string{fmt.Sprintf("Authorization: Bearer %s", j)})
			assert.NoError(t, err)
			fmt.Println(res)

			// Validator API
			fmt.Println("Calling Validator API")
			resp, err = resty.New().SetBaseURL(endpoints.ValidatorAPIURL).SetAuthToken(j).R().
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

```