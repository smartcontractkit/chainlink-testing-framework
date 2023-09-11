package docker

import (
	"context"
	"fmt"
	"time"

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
	ct, err := tc.GenericContainer(context.Background(), req)
	if err == nil {
		return ct, nil
	}
	for i := 0; i < RetryAttempts; i++ {
		log.Info().Err(err).Msgf("Cannot start %s container, restarting %d/%d", req.Name, i+1, RetryAttempts)
		timeout := 10 * time.Second
		err := ct.Stop(context.Background(), &timeout)
		if err != nil {
			log.Info().Err(err).Msgf("Cannot stop %s container", req.Name)
			continue
		}
		err = ct.Start(context.Background())
		if err == nil {
			break
		}
	}
	return ct, err
}
