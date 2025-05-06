package jd

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	tc "github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
)

const (
	TmpImageName            = "jd-local"
	GRPCPort         string = "42242"
	CSAEncryptionKey string = "!PASsword000!"
	WSRPCPort        string = "8080"
)

type Input struct {
	Image            string          `toml:"image"`
	GRPCPort         string          `toml:"grpc_port"`
	WSRPCPort        string          `toml:"wsrpc_port"`
	CSAEncryptionKey string          `toml:"csa_encryption_key"`
	DockerFilePath   string          `toml:"docker_file"`
	DockerContext    string          `toml:"docker_ctx"`
	JDSQLDumpPath    string          `toml:"jd_sql_dump_path"`
	DBInput          *postgres.Input `toml:"db"`
	Out              *Output         `toml:"out"`
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
	if in.Out != nil && in.Out.UseCache {
		return in.Out, nil
	}
	ctx := context.Background()
	defaults(in)
	jdImg := os.Getenv("CTF_JD_IMAGE")
	if jdImg != "" {
		in.Image = jdImg
	}
	if in.DBInput == nil {
		in.DBInput = defaultJDDB()
	}
	in.DBInput.JDSQLDumpPath = in.JDSQLDumpPath
	pgOut, err := postgres.NewPostgreSQL(in.DBInput)
	if err != nil {
		return nil, err
	}
	containerName := framework.DefaultTCName("jd")
	bindPort := fmt.Sprintf("%s/tcp", in.GRPCPort)
	req := tc.ContainerRequest{
		Name:     containerName,
		Image:    in.Image,
		Labels:   framework.DefaultTCLabels(),
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		ExposedPorts: []string{bindPort},
		HostConfigModifier: func(h *container.HostConfig) {
			// JobDistributor service is isolated from internet by default!
			framework.NoDNS(true, h)
			h.PortBindings = framework.MapTheSamePort(bindPort)
		},
		Env: map[string]string{
			"DATABASE_URL":              pgOut.JDInternalURL,
			"PORT":                      in.GRPCPort,
			"NODE_RPC_PORT":             in.WSRPCPort,
			"CSA_KEY_ENCRYPTION_SECRET": in.CSAEncryptionKey,
		},
		WaitingFor: tcwait.ForAll(
			tcwait.ForListeningPort(nat.Port(fmt.Sprintf("%s/tcp", in.GRPCPort))),
		),
	}
	if req.Image == "" {
		req.Image = TmpImageName
		if err := framework.BuildImage(in.DockerContext, in.DockerFilePath, req.Image); err != nil {
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
	host, err := framework.GetHost(c)
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
