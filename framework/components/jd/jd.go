package jd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/pods"
)

const (
	TmpImageName            = "jd-local"
	GRPCPort         string = "14231"
	CSAEncryptionKey string = "!PASsword000!"
	WSRPCPort        string = "8080"
	WSRPCHealthPort  string = "8081"
)

type Input struct {
	Image            string `toml:"image"`
	GRPCPort         string `toml:"grpc_port"`
	WSRPCPort        string `toml:"wsrpc_port"`
	CSAEncryptionKey string `toml:"csa_encryption_key"`
	DockerFilePath   string `toml:"docker_file"`
	DockerContext    string `toml:"docker_ctx"`
	JDSQLDumpPath    string `toml:"jd_sql_dump_path"`
	// DisableDNSIsolation keeps Docker's embedded DNS (127.0.0.11).
	// Leave false (default) to preserve historical isolation behavior.
	// Set true when JD must resolve peer service names (for example jd-db)
	// on a shared Docker network.
	DisableDNSIsolation bool            `toml:"disable_dns_isolation"`
	DBInput             *postgres.Input `toml:"db"`
	Out                 *Output         `toml:"out"`
}

type Output struct {
	UseCache         bool   `toml:"use_cache"`
	ContainerName    string `toml:"container_name"`
	DBContainerName  string `toml:"db_container_name"`
	ExternalGRPCUrl  string `toml:"grpc_url"`
	InternalGRPCUrl  string `toml:"internal_grpc_url"`
	ExternalWSRPCUrl string `toml:"wsrpc_url"`
	InternalWSRPCUrl string `toml:"internal_wsrpc_url"`
}

func defaults(in *Input) {
	if in.GRPCPort == "" {
		in.GRPCPort = GRPCPort
	}
	if in.WSRPCPort == "" {
		in.WSRPCPort = WSRPCPort
	}
	if in.CSAEncryptionKey == "" {
		in.CSAEncryptionKey = CSAEncryptionKey
	}
}

func defaultJDDB() *postgres.Input {
	return &postgres.Input{
		Image:      "postgres:16",
		Port:       14000,
		Name:       "jd-db",
		VolumeName: "jd",
		JDDatabase: true,
	}
}

func NewJD(in *Input) (*Output, error) {
	return NewWithContext(context.Background(), in)
}

func NewWithContext(ctx context.Context, in *Input) (*Output, error) {
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}
	defaults(in)
	jdImg := os.Getenv("CTF_JD_IMAGE")
	if jdImg != "" {
		// unset docker build context and file path to avoid conflicts, image provided via env var takes precedence
		in.Image = jdImg
		in.DockerContext = ""
		in.DockerFilePath = ""
	}
	if in.WSRPCPort == WSRPCHealthPort {
		return nil, fmt.Errorf("wsrpc port cannot be the same as wsrpc health port")
	}
	if in.DBInput == nil {
		in.DBInput = defaultJDDB()
		suffix := fmt.Sprintf("%d", time.Now().UnixNano())
		in.DBInput.Name = fmt.Sprintf("%s-%s", in.DBInput.Name, suffix)
		in.DBInput.VolumeName = fmt.Sprintf("%s-%s", in.DBInput.VolumeName, suffix)
	}
	in.DBInput.JDSQLDumpPath = in.JDSQLDumpPath
	pgOut, err := postgres.NewWithContext(ctx, in.DBInput)
	if err != nil {
		return nil, err
	}
	containerName := framework.DefaultTCName("jd")
	grpcPort := fmt.Sprintf("%s/tcp", in.GRPCPort)
	wsrpcPort := fmt.Sprintf("%s/tcp", in.WSRPCPort)
	wsHealthPort := fmt.Sprintf("%s/tcp", WSRPCHealthPort)

	if pods.K8sEnabled() {
		return nil, fmt.Errorf("K8s support is not yet implemented")
	}
	if err := framework.DefaultNetwork(nil); err != nil {
		return nil, fmt.Errorf("failed to ensure default docker network %q: %w", framework.DefaultNetworkName, err)
	}

	req := tc.ContainerRequest{
		Name:     containerName,
		Image:    in.Image,
		Labels:   framework.DefaultTCLabels(),
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		ExposedPorts: []string{grpcPort, wsrpcPort, wsHealthPort},
		HostConfigModifier: func(h *container.HostConfig) {
			// Default behavior keeps DNS isolation enabled for backwards compatibility.
			// Disable only when JD needs Docker service-name resolution (for example jd-db).
			framework.NoDNS(!in.DisableDNSIsolation, h)
			h.PortBindings = framework.MapTheSamePort(grpcPort, wsrpcPort)
			h.ExtraHosts = append(h.ExtraHosts, "host.docker.internal:host-gateway")
		},
		Env: map[string]string{
			"DATABASE_URL":              pgOut.JDInternalURL,
			"PORT":                      in.GRPCPort,
			"NODE_RPC_PORT":             in.WSRPCPort,
			"CSA_KEY_ENCRYPTION_SECRET": in.CSAEncryptionKey,
		},
		WaitingFor: tcwait.ForAll(
			tcwait.ForListeningPort(nat.Port(fmt.Sprintf("%s/tcp", in.GRPCPort))),
			wait.ForHTTP("/healthz").
				WithPort(nat.Port(fmt.Sprintf("%s/tcp", WSRPCHealthPort))). // WSRPC health endpoint uses different port than WSRPC
				WithStartupTimeout(1*time.Minute).
				WithPollInterval(200*time.Millisecond),
			NewGRPCHealthStrategy(nat.Port(fmt.Sprintf("%s/tcp", in.GRPCPort))).
				WithTimeout(1*time.Minute).
				WithPollInterval(200*time.Millisecond),
		),
	}
	if req.Image == "" {
		req.Image = TmpImageName
		if err := framework.BuildImage(in.DockerContext, in.DockerFilePath, req.Image, nil); err != nil {
			return nil, err
		}
		req.KeepImage = false
	}
	c, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	host, err := framework.GetHostWithContext(ctx, c)
	if err != nil {
		return nil, err
	}
	out := &Output{
		UseCache:         true,
		ContainerName:    containerName,
		DBContainerName:  pgOut.ContainerName,
		ExternalGRPCUrl:  fmt.Sprintf("%s:%s", host, in.GRPCPort),
		InternalGRPCUrl:  fmt.Sprintf("%s:%s", containerName, in.GRPCPort),
		ExternalWSRPCUrl: fmt.Sprintf("%s:%s", host, in.WSRPCPort),
		InternalWSRPCUrl: fmt.Sprintf("%s:%s", containerName, in.WSRPCPort),
	}
	in.Out = out
	return out, nil
}
