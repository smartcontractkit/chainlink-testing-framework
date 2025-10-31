package s3provider

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net"
	"strconv"

	"dario.cat/mergo"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	DefaultImage       = "minio/minio"
	DefaultName        = "minio"
	DefaultBucket      = "test-bucket"
	DefaultRegion      = "us-east-1"
	DefaultHost        = "minio"
	DefaultPort        = 9000
	DefaultConsolePort = 9001

	accessKeyLength = 20
	secretKeyLength = 40
)

type Minio struct {
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	ConsolePort int    `toml:"console_port"`
	AccessKey   string `toml:"access_key"`
	SecretKey   string `toml:"secret_key"`
	Bucket      string `toml:"bucket"`
	Region      string `toml:"region"`
}

type Input = Minio

type Output struct {
	SecretKey      string `toml:"secret_key"`
	AccessKey      string `toml:"access_key"`
	Bucket         string `toml:"bucket"`
	ConsoleURL     string `toml:"console_url"`
	ConsoleBaseURL string `toml:"console_base_url"`
	Endpoint       string `toml:"endpoint"`
	BaseEndpoint   string `toml:"base_endpoint"`
	Region         string `toml:"region"`
	UseCache       bool   `toml:"use_cache"`
}

func (m Minio) Output() *Output {
	return &Output{
		AccessKey:      m.GetAccessKey(),
		SecretKey:      m.GetSecretKey(),
		Bucket:         m.GetBucket(),
		ConsoleURL:     m.GetConsoleURL(),
		ConsoleBaseURL: m.GetConsoleBaseURL(),
		Endpoint:       m.GetEndpoint(),
		BaseEndpoint:   m.GetBaseEndpoint(),
		Region:         m.GetRegion(),
	}
}

func (m Minio) GetSecretKey() string {
	return m.SecretKey
}

func (m Minio) GetAccessKey() string {
	return m.AccessKey
}

func (m Minio) GetBucket() string {
	return m.Bucket
}

func (m Minio) GetConsoleURL() string {
	return fmt.Sprintf("http://%s", net.JoinHostPort(m.Host, strconv.Itoa(m.ConsolePort)))
}

func (m Minio) GetConsoleBaseURL() string {
	return fmt.Sprintf("http://%s", net.JoinHostPort(DefaultHost, strconv.Itoa(m.ConsolePort)))
}

func (m Minio) GetEndpoint() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

func (m Minio) GetBaseEndpoint() string {
	return fmt.Sprintf("%s:%d", DefaultHost, m.Port)
}

func (m Minio) GetRegion() string {
	return m.Region
}

type Option func(*Minio)

type MinioFactory struct{}

func NewMinioFactory() ProviderFactory {
	return MinioFactory{}
}

func (mf MinioFactory) NewFrom(input *Input) (*Output, error) {
	return mf.NewWithContextFrom(context.Background(), input)
}

func (mf MinioFactory) NewWithContextFrom(ctx context.Context, input *Input) (*Output, error) {
	// Fill in defaults on empty
	err := mergo.Merge(input, DefaultMinio())
	if err != nil {
		return nil, err
	}

	provider, err := mf.run(ctx, input)
	if err != nil {
		return nil, err
	}
	return provider.Output(), nil
}

func DefaultMinio() *Minio {
	return &Minio{
		Host:        DefaultHost,
		Port:        DefaultPort,
		ConsolePort: DefaultConsolePort,
		AccessKey:   randomStr(accessKeyLength),
		SecretKey:   randomStr(secretKeyLength),
		Bucket:      DefaultBucket,
		Region:      DefaultRegion,
	}
}

func (mf MinioFactory) New(options ...Option) (Provider, error) {
	return mf.NewWithContext(context.Background(), options...)
}

func (mf MinioFactory) NewWithContext(ctx context.Context, options ...Option) (Provider, error) {
	m := DefaultMinio()

	for _, opt := range options {
		opt(m)
	}

	return mf.run(ctx, m)
}

func (mf MinioFactory) run(ctx context.Context, m *Minio) (Provider, error) {
	var err error

	containerName := framework.DefaultTCName(DefaultName)
	bindPort := fmt.Sprintf("%d/tcp", m.Port)
	bindConsolePort := fmt.Sprintf("%d/tcp", m.ConsolePort)
	networks := []string{"compose_default"}
	networkAliases := map[string][]string{
		"compose_default": {DefaultName},
	}

	if len(framework.DefaultNetworkName) > 0 {
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
			"MINIO_ROOT_USER":     m.AccessKey,
			"MINIO_ROOT_PASSWORD": m.SecretKey,
			"MINIO_BUCKET":        m.Bucket,
		},
		Entrypoint: []string{
			"minio",
			"server",
			"/data",
			"--address",
			fmt.Sprintf(":%d", m.Port),
			"--console-address",
			fmt.Sprintf(":%d", m.ConsolePort),
		},
		HostConfigModifier: func(h *container.HostConfig) {
			framework.NoDNS(true, h)
			h.PortBindings = nat.PortMap{
				nat.Port(bindPort): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: strconv.Itoa(m.Port),
					},
				},
				nat.Port(bindConsolePort): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: strconv.Itoa(m.ConsolePort),
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
	})
	if err != nil {
		return nil, err
	}

	m.Host, err = framework.GetHostWithContext(ctx, c)
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

	myPolicy := fmt.Sprintf(`
		{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Action": [
						"s3:GetObject"
					],
					"Effect": "Allow",
					"Principal": {
						"AWS": [
							"*"
						]
					},
					"Resource": [
						"arn:aws:s3:::%s/*"
					],
					"Sid": ""
				}
			]
		}
	`, m.GetBucket())

	err = minioClient.SetBucketPolicy(ctx, m.GetBucket(), myPolicy)
	if err != nil {
		framework.L.Warn().Str("error", err.Error()).Msg("failed to set public policy to minio bucket")

		return nil, err
	}

	return m, nil
}

func WithPort(port int) Option {
	return func(m *Minio) {
		m.Port = port
	}
}

func WithConsolePort(consolePort int) Option {
	return func(m *Minio) {
		m.ConsolePort = consolePort
	}
}

func WithAccessKey(accessKey string) Option {
	return func(m *Minio) {
		m.AccessKey = accessKey
	}
}

func WithSecretKey(secretKey string) Option {
	return func(m *Minio) {
		m.SecretKey = secretKey
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.IntN(len(letterBytes))] // #nosec G404
	}

	return string(b)
}
