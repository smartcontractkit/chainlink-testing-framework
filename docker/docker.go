package docker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
)

const RetryAttempts = 3

func CreateNetwork() (*tc.DockerNetwork, error) {
	uuidObj, _ := uuid.NewRandom()
	var networkName = fmt.Sprintf("network-%s", uuidObj.String())
	network, err := tc.GenericNetwork(context.Background(), tc.GenericNetworkRequest{
		NetworkRequest: tc.NetworkRequest{
			Name:           networkName,
			CheckDuplicate: true,
		},
	})
	if err != nil {
		return nil, err
	}
	dockerNetwork, ok := network.(*tc.DockerNetwork)
	if !ok {
		return nil, fmt.Errorf("failed to cast network to *dockertest.Network")
	}
	log.Trace().Any("network", dockerNetwork).Msgf("created network")
	return dockerNetwork, nil
}

func StartContainerWithRetry(req tc.GenericContainerRequest) (tc.Container, error) {
	var ct tc.Container
	var err error
	for i := 0; i < RetryAttempts; i++ {
		ct, err = tc.GenericContainer(context.Background(), req)
		if err == nil {
			break
		}
		log.Info().Err(err).Msgf("Cannot start %s container, retrying %d/%d", req.Name, i+1, RetryAttempts)
		req.Started = false
	}
	return ct, err
}
