package blockchain

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain/canton"
)

// newCanton sets up a Canton blockchain network with the specified number of validators.
// It creates a Docker network and starts the necessary containers for Postgres, Canton, Splice, and an Nginx reverse proxy.
//
// The reverse proxy is used to allow access to all validator participants through a single HTTP endpoint.
// The following routes are configured for each participant and the Super Validator (SV):
//   - http://[PARTICIPANT].json-ledger-api.localhost:[PORT] 	-> JSON Ledger API
//   - grpc://[PARTICIPANT].grpc-ledger-api.localhost:[PORT] 	-> gRPC Ledger API
//   - http://[PARTICIPANT].admin-api.localhost:[PORT] 			-> Admin API
//   - http://[PARTICIPANT].wallet.localhost:[PORT] 			-> Wallet API
//   - http://[PARTICIPANT].http-health-check.localhost:[PORT] 	-> HTTP Health Check
//   - grpc://[PARTICIPANT].grpc-health-check.localhost:[PORT] 	-> gRPC Health Check
// To access a participant's endpoints, replace [PARTICIPANT] with the participant's identifier, i.e. `sv`, `participant01`, `participant02`, ...
//
// Additionally, the global Scan service is accessible via:
//   - http://scan.localhost:[PORT]/api/scan 					-> Scan API
//   - http://scan.localhost:[PORT]/registry 					-> Scan Registry
// The PORT is the same for all routes and is specified in the input parameters.
//
// Note: The maximum number of validators supported is 99, participants are numbered starting from `participant01` through `participant99`.
func newCanton(ctx context.Context, in *Input) (*Output, error) {
	if in.NumberOfCantonValidators >= 100 {
		return nil, fmt.Errorf("number of validators too high: %d, max is 99", in.NumberOfCantonValidators)
	}

	// Create separate Docker network for Canton stack
	dockerNetwork, err := network.New(ctx, network.WithAttachable())
	if err != nil {
		return nil, err
	}

	// Set up Postgres container
	postgresReq := canton.PostgresContainerRequest(in.NumberOfCantonValidators, dockerNetwork.Name)
	_, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: postgresReq,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	// Set up Canton container
	cantonReq := canton.CantonContainerRequest(dockerNetwork.Name, in.NumberOfCantonValidators, in.Image)
	_, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: cantonReq,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	// Set up Splice container
	spliceReq := canton.SpliceContainerRequest(dockerNetwork.Name, in.NumberOfCantonValidators, in.Image)
	_, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: spliceReq,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	// Set up Nginx container
	nginxReq := canton.NginxContainerRequest(dockerNetwork.Name, in.NumberOfCantonValidators, in.Port)
	nginxContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: nginxReq,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := nginxContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	return &Output{
		UseCache:      false,
		Type:          in.Type,
		Family:        FamilyCanton,
		ContainerName: nginxReq.Name,
		Nodes: []*Node{
			{
				ExternalHTTPUrl: fmt.Sprintf("http://%s:%s", host, in.Port),
			},
		},
	}, nil
}
