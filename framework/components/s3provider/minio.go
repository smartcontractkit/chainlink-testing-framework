package s3provider

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
)

const (
	DefaultImage       = "minio/minio"
	DefaultName        = "minio"
	DefaultBucket      = "test-bucket"
	DefaultRegion      = "us-east-1"
	DefaultPort        = 9000
	DefaultConsolePort = 9001

	accessKeyLength = 20
	secretKeyLength = 40
)

type Minio struct {
	host        string
	port        int
	consolePort int
	accessKey   string
	secretKey   string
	bucket      string
	region      string
	keep        bool
}

func (m Minio) GetSecretKey() string {
	return m.secretKey
}

func (m Minio) GetAccessKey() string {
	return m.accessKey
}

func (m Minio) GetBucket() string {
	return m.bucket
}

func (m Minio) GetConsoleURL() string {
	return fmt.Sprintf("http://%s", net.JoinHostPort(m.host, strconv.Itoa(m.consolePort)))
}

func (m Minio) GetEndpoint() string {
	return fmt.Sprintf("%s:%d", m.host, m.port)
}

func (m Minio) GetRegion() string {
	return m.region
}

type Option func(*Minio)

type MinioFactory struct{}

func NewMinioFactory() ProviderFactory {
	return MinioFactory{}
}

func (mf MinioFactory) New(options ...Option) (Provider, error) {
	m := &Minio{
		port:        DefaultPort,
		consolePort: DefaultConsolePort,
		accessKey:   randomStr(accessKeyLength),
		secretKey:   randomStr(secretKeyLength),
		bucket:      DefaultBucket,
		region:      DefaultRegion,
		keep:        false,
	}

	for _, opt := range options {
		opt(m)
	}

	var (
		tcRyukDisabled string
		err            error
	)

	if m.keep {
		// store original env var to value
		tcRyukDisabled = os.Getenv("TESTCONTAINERS_RYUK_DISABLED")
		err = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

		if err != nil {
			return nil, err
		}
	}

	ctx := context.Background()
	containerName := framework.DefaultTCName(DefaultName)
	bindPort := fmt.Sprintf("%d/tcp", m.port)
	bindConsolePort := fmt.Sprintf("%d/tcp", m.consolePort)
	networks := []string{"compose_default"}
	networkAliases := map[string][]string{
		"compose_default": {DefaultName},
	}

	if len(framework.DefaultNetworkName) == 0 {
		// attach default ctf network if initiated
		networks = append(networks, framework.DefaultNetworkName)
		networkAliases[framework.DefaultNetworkName] = []string{
			containerName,
			DefaultName,
		}
	}

	req := tc.ContainerRequest{
		Name:           containerName,
		Image:          DefaultImage,
		Labels:         framework.DefaultTCLabels(),
		Networks:       networks,
		NetworkAliases: networkAliases,
		ExposedPorts: []string{
			bindPort,
			bindConsolePort,
		},
		Env: map[string]string{
			"MINIO_ROOT_USER":     m.accessKey,
			"MINIO_ROOT_PASSWORD": m.secretKey,
			"MINIO_BUCKET":        DefaultBucket,
		},
		Entrypoint: []string{
			"minio",
			"server",
			"/data",
			"--address",
			fmt.Sprintf(":%d", m.port),
			"--console-address",
			fmt.Sprintf(":%d", m.consolePort),
		},
		HostConfigModifier: func(h *container.HostConfig) {
			framework.NoDNS(true, h)
			h.PortBindings = nat.PortMap{
				nat.Port(bindPort): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: strconv.Itoa(m.port),
					},
				},
				nat.Port(bindConsolePort): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: strconv.Itoa(m.consolePort),
					},
				},
			}
		},
		WaitingFor: tcwait.ForAll(
			tcwait.ForListeningPort(nat.Port(bindPort)),
			tcwait.ForListeningPort(nat.Port(bindConsolePort)),
		),
	}

	c, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            m.keep,
	})
	if err != nil {
		return nil, err
	}

	m.host, err = framework.GetHost(c)
	if err != nil {
		return nil, err
	}

	// Initialize minio client object.
	minioClient, err := minio.New(m.GetEndpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(m.GetAccessKey(), m.GetSecretKey(), ""),
		Secure: false,
	})
	if err != nil {
		framework.L.Warn().Str("error", err.Error()).Msg("failed to create minio client")

		return nil, err
	}

	// Initialize default bucket
	err = minioClient.MakeBucket(ctx, m.GetBucket(), minio.MakeBucketOptions{Region: m.GetRegion()})
	if err != nil {
		framework.L.Warn().Str("error", err.Error()).Msg("failed to create minio bucket")

		return nil, err
	}

	if m.keep {
		// reverse env var to prev. value
		err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", tcRyukDisabled)
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func WithPort(port int) Option {
	return func(m *Minio) {
		m.port = port
	}
}

func WithConsolePort(consolePort int) Option {
	return func(m *Minio) {
		m.consolePort = consolePort
	}
}

func WithKeep() Option {
	return func(m *Minio) {
		m.keep = true
	}
}

func WithAccessKey(accessKey string) Option {
	return func(m *Minio) {
		m.accessKey = accessKey
	}
}

func WithSecretKey(secretKey string) Option {
	return func(m *Minio) {
		m.secretKey = secretKey
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}

	return string(b)
}
