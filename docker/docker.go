package docker

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
)

const RetryAttempts = 3

func CreateNetwork(l zerolog.Logger) (*tc.DockerNetwork, error) {
	uuidObj, _ := uuid.NewRandom()
	var networkName = fmt.Sprintf("network-%s", uuidObj.String())
	ryukImage, err := mirror.GetImage("testcontainers/ryuk")
	if err != nil {
		return nil, err
	}
	reaperCO := tc.WithImageName(ryukImage)
	network, err := tc.GenericNetwork(context.Background(), tc.GenericNetworkRequest{
		NetworkRequest: tc.NetworkRequest{
			Name:           networkName,
			CheckDuplicate: true,
			EnableIPv6:     false, // disabling due to https://github.com/moby/moby/issues/42442
			ReaperOptions: []tc.ContainerOption{
				reaperCO,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	dockerNetwork, ok := network.(*tc.DockerNetwork)
	if !ok {
		return nil, fmt.Errorf("failed to cast network to *dockertest.Network")
	}
	l.Trace().Any("network", dockerNetwork).Msgf("created network")
	return dockerNetwork, nil
}

func StartContainerWithRetry(l zerolog.Logger, req tc.GenericContainerRequest) (tc.Container, error) {
	var ct tc.Container
	var err error
	for i := 0; i < RetryAttempts; i++ {
		ct, err = tc.GenericContainer(context.Background(), req)
		if err == nil {
			break
		}
		l.Info().Err(err).Msgf("Cannot start %s container, retrying %d/%d", req.Name, i+1, RetryAttempts)
		if ct != nil {
			err := ct.Terminate(context.Background())
			if err != nil {
				l.Error().Err(err).Msgf("Cannot terminate %s container to initiate restart", req.Name)
				return nil, err
			}
		}
		// TODO we should dynamically check if container was removed, otherwise retry will fail anyway, just becasuse
		// container with the same name already exists | that or regenerate container name
		// But regenerating container name will not help to remove temporary data created on host machine and that can be
		// a problem if data is stateful
		req.Reuse = false
	}
	return ct, err
}
