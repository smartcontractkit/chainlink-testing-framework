package blockchain

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain/canton"
)

func newCanton(ctx context.Context, in *Input) (*Output, error) {
	if in.NumberOfCantonValidators >= 100 {
		return nil, fmt.Errorf("number of validators too high: %d, max is 99", in.NumberOfCantonValidators)
	}

	// TODO - remove debug prints
	fmt.Println("Starting Canton blockchain node...")
	fmt.Println("Creating network...")
	dockerNetwork, err := network.New(ctx, network.WithAttachable())
	if err != nil {
		return nil, err
	}
	fmt.Println("Network created:", dockerNetwork.Name)

	// Set up Postgres container
	postgresReq := canton.PostgresContainerRequest(in.NumberOfCantonValidators, dockerNetwork.Name)
	fmt.Printf("Starting postgres container %s...\n", postgresReq.Name)
	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: postgresReq,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	_ = c
	fmt.Println("Postgres container started")

	// Set up Canton container
	cantonReq := canton.CantonContainerRequest(dockerNetwork.Name, in.NumberOfCantonValidators, in.Image)
	fmt.Printf("Starting canton container %s...\n", cantonReq.Name)
	cantonContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: cantonReq,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	_ = cantonContainer
	fmt.Println("Canton container started")

	// Set up Splice container
	spliceReq := canton.SpliceContainerRequest(dockerNetwork.Name, in.NumberOfCantonValidators, in.Image)
	fmt.Printf("Starting splice container %s...\n", spliceReq.Name)
	spliceContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: spliceReq,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	_ = spliceContainer
	fmt.Println("Splice container started")

	// Set up Nginx container
	nginxReq := canton.NginxContainerRequest(dockerNetwork.Name, in.NumberOfCantonValidators, in.Port)
	fmt.Printf("Starting nginx container %s...\n", nginxReq.Name)
	nginxContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: nginxReq,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("Nginx container started")

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
				InternalHTTPUrl: fmt.Sprintf("http://%s:%s", nginxReq.Name, in.Port), // TODO - should be docker-internal port instead?
			},
		},
	}, nil
}
