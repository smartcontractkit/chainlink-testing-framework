package docker

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/testcontext"
)

const RetryAttempts = 3
const defaultRyukImage = "testcontainers/ryuk:0.5.1"

func CreateNetwork(l zerolog.Logger) (*tc.DockerNetwork, error) {
	uuidObj, _ := uuid.NewRandom()
	var networkName = fmt.Sprintf("network-%s", uuidObj.String())
	ryukImage := mirror.AddMirrorToImageIfSet(defaultRyukImage)
	reaperCO := tc.WithImageName(ryukImage)
	network, err := tc.GenericNetwork(testcontext.Get(nil), tc.GenericNetworkRequest{
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

type StartContainerRetrier func(l zerolog.Logger, startErr error, req tc.GenericContainerRequest) (tc.Container, error)

var NaiveRetrier = func(l zerolog.Logger, startErr error, req tc.GenericContainerRequest) (tc.Container, error) {
	l.Debug().
		Str("Start error", startErr.Error()).
		Str("Retrier", "NaiveRetrier").
		Msgf("Attempting to start %s container", req.Name)

	ct, err := tc.GenericContainer(testcontext.Get(nil), req)
	if err == nil {
		l.Debug().
			Str("Retrier", "NaiveRetrier").
			Msgf("Successfully started %s container", req.Name)
		return ct, nil
	}
	if ct != nil {
		err := ct.Terminate(testcontext.Get(nil))
		if err != nil {
			l.Error().
				Err(err).
				Msgf("Cannot terminate %s container to initiate restart", req.Name)
			return nil, err
		}
	}
	req.Reuse = false

	l.Debug().
		Str("Original start error", startErr.Error()).
		Str("Current start error", err.Error()).
		Str("Retrier", "NaiveRetrier").
		Msgf("Failed to start %s container,", req.Name)

	return nil, startErr
}

var LinuxPlatoformImageRetrier = func(l zerolog.Logger, startErr error, req tc.GenericContainerRequest) (tc.Container, error) {
	// if it's nil we don't know if we can handle it so we won't try
	if startErr == nil {
		return nil, startErr
	}

	// a bit lame, but that's the lame error we get in case there's no specific image for our platform :facepalm:
	if !strings.Contains(startErr.Error(), "No such image") {
		l.Debug().
			Str("Start error", startErr.Error()).
			Str("Retrier", "PlatoformImageRetrier").
			Msgf("Won't try to start %s container again, unsupported error", req.Name)
		return nil, startErr
	}

	l.Debug().
		Str("Start error", startErr.Error()).
		Str("Retrier", "PlatoformImageRetrier").
		Msgf("Attempting to start %s container", req.Name)

	originalPlatform := req.ImagePlatform
	req.ImagePlatform = "linux/x86_64"

	ct, err := tc.GenericContainer(testcontext.Get(nil), req)
	if err == nil {
		l.Debug().
			Str("Retrier", "PlatoformImageRetrier").
			Msgf("Successfully started %s container", req.Name)
		return ct, nil
	}

	req.ImagePlatform = originalPlatform

	if ct != nil {
		err := ct.Terminate(testcontext.Get(nil))
		if err != nil {
			l.Error().Err(err).Msgf("Cannot terminate %s container to initiate restart", req.Name)
			return nil, err
		}
	}

	l.Debug().
		Str("Original start error", startErr.Error()).
		Str("Current start error", err.Error()).
		Str("Retrier", "PlatoformImageRetrier").
		Msgf("Failed to start %s container,", req.Name)

	return nil, startErr
}

// StartContainerWithRetry attempts to start a container with 3 retry attempts.
// It will try to start the container with the provided retriers, if none are provided it will use the default retriers.
// Default being: 1. tries to download image for "linux/x86_64" platform 2. simply starts again without changing anything
func StartContainerWithRetry(l zerolog.Logger, req tc.GenericContainerRequest, retriers ...StartContainerRetrier) (tc.Container, error) {
	var ct tc.Container
	var err error

	ct, err = tc.GenericContainer(testcontext.Get(nil), req)
	if err == nil {
		return ct, nil
	}

	if len(retriers) == 0 {
		retriers = append(retriers, LinuxPlatoformImageRetrier, NaiveRetrier)
	}

	for i := 0; i < RetryAttempts; i++ {
		l.Info().Err(err).Msgf("Cannot start %s container, retrying %d/%d", req.Name, i+1, RetryAttempts)

		for _, retrier := range retriers {
			ct, err = retrier(l, err, req)
			if err == nil {
				return ct, nil
			}
		}
	}

	return nil, err
}
