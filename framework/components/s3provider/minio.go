package s3provider

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"
	"math/rand"
	"strconv"
)

const (
	DefaultImage       = "minio/minio"
	DefaultName        = "minio"
	DefaultBucket      = "test-bucket"
	DefaultRegion      = "us-east-1"
	DefaultPort        = 9000
	DefaultConsolePort = 9001
)

type Minio struct {
	host        string
	port        int
	consolePort int
	accessKey   string
	secretKey   string
	bucket      string
	region      string
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

func (m Minio) GetURL() string {
	return fmt.Sprintf("http://%s:%d", m.host, m.port)
}

func (m Minio) GetConsoleURL() string {
	return fmt.Sprintf("http://%s:%d", m.host, m.consolePort)
}

func (m Minio) GetEndpoint() string {
	return fmt.Sprintf("%s:%d", m.host, m.consolePort)
}

func (m Minio) GetRegion() string {
	return m.region
}

type Option func(*Minio)

type MinioFactory struct{}

func NewMinioFactory() ProviderFactory {
	return MinioFactory{}
}

func (mf MinioFactory) NewProvider(options ...Option) (Provider, error) {
	m := &Minio{
		port:        DefaultPort,
		consolePort: DefaultConsolePort,
		accessKey:   randomStr(20),
		secretKey:   randomStr(40),
		bucket:      DefaultBucket,
		region:      DefaultRegion,
	}

	for _, opt := range options {
		opt(m)
	}

	ctx := context.Background()
	containerName := framework.DefaultTCName(DefaultName)
	bindPort := fmt.Sprintf("%d/tcp", m.port)
	bindConsolePort := fmt.Sprintf("%d/tcp", m.consolePort)

	req := tc.ContainerRequest{
		Name:     containerName,
		Image:    DefaultImage,
		Labels:   framework.DefaultTCLabels(),
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		ExposedPorts: []string{bindPort, bindConsolePort},
		Env: map[string]string{
			"MINIO_ROOT_USER":     m.accessKey,
			"MINIO_ROOT_PASSWORD": m.secretKey,
			"MINIO_BUCKET":        DefaultBucket,
		},
		Entrypoint: []string{
			"minio",
			"server",
			"data",
			"--address",
			fmt.Sprintf(":%d", m.port),
			"--console-address",
			fmt.Sprintf(":%d", m.consolePort),
		},
		HostConfigModifier: func(h *container.HostConfig) {
			framework.NoDNS(true, h)
			h.PortBindings = framework.MapTheSamePort(bindPort)
		},
		WaitingFor: tcwait.ForAll(
			tcwait.ForListeningPort(nat.Port(fmt.Sprintf("%d/tcp", m.port))),
		),
	}
	req.HostConfigModifier = func(h *container.HostConfig) {
		h.PortBindings = nat.PortMap{
			nat.Port(bindPort): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: strconv.Itoa(m.port),
				},
			},
			nat.Port(bindPort): []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: strconv.Itoa(m.consolePort),
				},
			},
		}
	}

	c, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := framework.GetHost(c)
	if err != nil {
		return nil, err
	}
	m.host = host

	// Initialize minio client object.
	minioClient, err := minio.New(m.GetEndpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(m.GetAccessKey(), m.GetSecretKey(), ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	// Initialize default bucket
	err = minioClient.MakeBucket(ctx, m.GetBucket(), minio.MakeBucketOptions{Region: m.GetRegion()})
	if err != nil {
		return nil, err
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
