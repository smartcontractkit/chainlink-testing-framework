package test_env

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/mirror"
)

type MockServer struct {
	EnvComponent
	Client           *ctfClient.MockserverClient
	Endpoint         string
	InternalEndpoint string
	EAMockUrls       []*url.URL
	t                *testing.T
	l                zerolog.Logger
}

func NewMockServer(networks []string, opts ...EnvComponentOption) *MockServer {
	ms := &MockServer{
		EnvComponent: EnvComponent{
			ContainerName: fmt.Sprintf("%s-%s", "mockserver", uuid.NewString()[0:8]),
			Networks:      networks,
		},
		l: log.Logger,
	}
	for _, opt := range opts {
		opt(&ms.EnvComponent)
	}
	return ms
}

func (ms *MockServer) WithTestLogger(t *testing.T) *MockServer {
	ms.l = logging.GetTestLogger(t)
	ms.t = t
	return ms
}

func (ms *MockServer) SetExternalAdapterMocks(count int) error {
	for i := 0; i < count; i++ {
		path := fmt.Sprintf("/ea-%d", i)
		err := ms.Client.SetRandomValuePath(path)
		if err != nil {
			return err
		}
		cName, err := ms.Container.Name(context.Background())
		if err != nil {
			return err
		}
		cName = strings.Replace(cName, "/", "", -1)
		eaUrl, err := url.Parse(fmt.Sprintf("http://%s:%s%s",
			cName, "1080", path))
		if err != nil {
			return err
		}
		ms.EAMockUrls = append(ms.EAMockUrls, eaUrl)
	}
	return nil
}

func (ms *MockServer) StartContainer() error {
	l := tc.Logger
	if ms.t != nil {
		l = logging.CustomT{
			T: ms.t,
			L: ms.l,
		}
	}
	cr, err := ms.getContainerRequest()
	if err != nil {
		return err
	}
	c, err := docker.StartContainerWithRetry(ms.l, tc.GenericContainerRequest{
		ContainerRequest: cr,
		Reuse:            true,
		Started:          true,
		Logger:           l,
	})
	if err != nil {
		return fmt.Errorf("cannot start MockServer container: %w", err)
	}
	ms.Container = c
	endpoint, err := GetEndpoint(context.Background(), c, "http")
	if err != nil {
		return err
	}
	ms.l.Info().Any("endpoint", endpoint).Str("containerName", ms.ContainerName).
		Msgf("Started MockServer container")
	ms.Endpoint = endpoint
	ms.InternalEndpoint = fmt.Sprintf("http://%s:%s", ms.ContainerName, "1080")

	client := ctfClient.NewMockserverClient(&ctfClient.MockserverConfig{
		LocalURL:   endpoint,
		ClusterURL: ms.InternalEndpoint,
	})
	if err != nil {
		return fmt.Errorf("cannot create MockServer client: %w", err)
	}
	ms.Client = client

	return nil
}

func (ms *MockServer) getContainerRequest() (tc.ContainerRequest, error) {
	msImage, err := mirror.GetImage("mockserver/mockserver")
	if err != nil {
		return tc.ContainerRequest{}, err
	}
	return tc.ContainerRequest{
		Name:         ms.ContainerName,
		Image:        msImage,
		ExposedPorts: []string{"1080/tcp"},
		Env: map[string]string{
			"SERVER_PORT": "1080",
		},
		Networks: ms.Networks,
		WaitingFor: tcwait.ForLog("INFO 1080 started on port: 1080").
			WithStartupTimeout(30 * time.Second).
			WithPollInterval(100 * time.Millisecond),
	}, nil
}
