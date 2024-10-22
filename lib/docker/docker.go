package docker

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/mirror"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"
)

const RetryAttempts = 3
const defaultRyukImage = "testcontainers/ryuk:0.5.1"

func CreateNetwork(l zerolog.Logger) (*tc.DockerNetwork, error) {
	uuidObj, _ := uuid.NewRandom()
	var networkName = fmt.Sprintf("network-%s", uuidObj.String())
	ryukImage := mirror.AddMirrorToImageIfSet(defaultRyukImage)
	// currently there's no way to use custom Ryuk image with testcontainers-go v0.28.0 :/
	// but we can go around it, by setting TESTCONTAINERS_HUB_IMAGE_NAME_PREFIX env var to
	// our custom registry and then using the default Ryuk image
	//nolint:staticcheck
	reaperCO := tc.WithImageName(ryukImage)
	f := false
	//nolint:staticcheck
	network, err := tc.GenericNetwork(testcontext.Get(nil), tc.GenericNetworkRequest{
		//nolint:staticcheck
		NetworkRequest: tc.NetworkRequest{
			Name:           networkName,
			CheckDuplicate: true,
			EnableIPv6:     &f, // disabling due to https://github.com/moby/moby/issues/42442
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

	oldName := req.Name
	req.Name = req.Name + "-naive-retry"

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

	req.Name = oldName

	l.Debug().
		Str("Original start error", startErr.Error()).
		Str("Current start error", err.Error()).
		Str("Retrier", "NaiveRetrier").
		Msgf("Failed to start %s container,", req.Name)

	return nil, startErr
}

var LinuxPlatformImageRetrier = func(l zerolog.Logger, startErr error, req tc.GenericContainerRequest) (tc.Container, error) {
	// if it's nil we don't know if we can handle it so we won't try
	if startErr == nil {
		return nil, startErr
	}
	req.Reuse = false // We need to force a new container to be created
	oldName := req.Name
	req.Name = req.Name + "-linux-retry"

	// a bit lame, but that's the lame error we get in case there's no specific image for our platform :facepalm:
	if !strings.Contains(startErr.Error(), "No such image") {
		l.Debug().
			Str("Start error", startErr.Error()).
			Str("Retrier", "PlatformImageRetrier").
			Msgf("Won't try to start %s container again, unsupported error", req.Name)
		return nil, startErr
	}

	l.Debug().
		Str("Start error", startErr.Error()).
		Str("Retrier", "PlatformImageRetrier").
		Msgf("Attempting to start %s container", req.Name)

	originalPlatform := req.ImagePlatform
	req.ImagePlatform = "linux/x86_64"

	ct, err := tc.GenericContainer(testcontext.Get(nil), req)
	if err == nil {
		l.Debug().
			Str("Retrier", "PlatformImageRetrier").
			Msgf("Successfully started %s container", req.Name)
		return ct, nil
	}

	req.ImagePlatform = originalPlatform
	req.Name = oldName

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
		Str("Retrier", "PlatformImageRetrier").
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
		retriers = append(retriers, LinuxPlatformImageRetrier, NaiveRetrier)
	}

	l.Info().Err(err).Msgf("Cannot start %s container, retrying", req.Name)

	req.Reuse = true // Try and see if we can reuse the container for a retry
	for _, retrier := range retriers {
		ct, err = retrier(l, err, req)
		if err == nil {
			return ct, nil
		}
	}

	return nil, err
}
