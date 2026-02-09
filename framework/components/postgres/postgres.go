package postgres

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/smartcontractkit/pods"
	"github.com/smartcontractkit/pods/imports/k8s"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components"
)

const (
	User              = "chainlink"
	Password          = "thispasswordislongenough"
	Port              = "5432"
	ExposedStaticPort = 13000
	Database          = "chainlink"
	JDDatabase        = "job-distributor-db"
	DBVolumeName      = "postgresql_data"
)

type Input struct {
	// Image PostgreSQL Docker image in format: $registry:$tag
	Image string `toml:"image" validate:"required" comment:"PostgreSQL Docker image in format: $registry:$tag"`
	// Port PostgreSQL connection port
	Port int `toml:"port" comment:"PostgreSQL connection port"`
	// Name PostgreSQL container name
	Name string `toml:"name" comment:"PostgreSQL container name"`
	// VolumeName PostgreSQL Docker volume name
	VolumeName string `toml:"volume_name" comment:"PostgreSQL docker volume name"`
	// Databases number of pre-created databases for Chainlink nodes
	Databases int `toml:"databases" comment:"Number of pre-created databases for Chainlink nodes"`
	// JDDatabase whether to create JobDistributor database or not
	JDDatabase bool `toml:"jd_database" comment:"Whether to create JobDistributor database or not"`
	// JDSQLDumpPath JobDistributor SQL dump path to load
	JDSQLDumpPath string `toml:"jd_sql_dump_path" comment:"JobDistributor database dump path to load"`
	// PullImage whether to pull PostgreSQL image or not
	PullImage bool `toml:"pull_image" comment:"Whether to pull PostgreSQL image or not"`
	// ContainerResources Docker container resources
	ContainerResources *framework.ContainerResources `toml:"resources" comment:"Docker container resources"`
	// Out PostgreSQL config output
	Out *Output `toml:"out" comment:"PostgreSQL config output"`
}

type Output struct {
	// Url PostgreSQL connection Url
	Url string `toml:"url" comment:"PostgreSQL connection URL"`
	// ContainerName PostgreSQL Docker container name
	ContainerName string `toml:"container_name" comment:"Docker container name"`
	// InternalURL PostgreSQL internal connection URL
	InternalURL string `toml:"internal_url" comment:"PostgreSQL internal connection URL"`
	// JDUrl PostgreSQL external connection URL to JobDistributor database
	JDUrl string `toml:"jd_url" comment:"PostgreSQL internal connection URL to JobDistributor database"`
	// JDInternalURL PostgreSQL internal connection URL to JobDistributor database
	JDInternalURL string `toml:"jd_internal_url" comment:"PostgreSQL internal connection URL to JobDistributor database"`
}

func NewPostgreSQL(in *Input) (*Output, error) {
	return NewWithContext(context.Background(), in)
}

func NewWithContext(ctx context.Context, in *Input) (*Output, error) {
	bindPort := fmt.Sprintf("%s/tcp", Port)
	var containerName string
	if in.Name == "" {
		containerName = "ns-postgresql"
	} else {
		containerName = in.Name
	}

	var sqlCommands []string
	for i := 0; i <= in.Databases; i++ {
		sqlCommands = append(sqlCommands,
			fmt.Sprintf("CREATE DATABASE db_%d;", i),
			fmt.Sprintf("\\c db_%d", i),
			"CREATE EXTENSION pg_stat_statements;",
		)
	}
	sqlCommands = append(sqlCommands, "ALTER USER chainlink WITH SUPERUSER;")
	if in.JDDatabase {
		if in.JDSQLDumpPath != "" {
			// if we have a full dump we replace RDS specific commands and apply it creating db and filling the tables
			d, err := os.ReadFile(in.JDSQLDumpPath)
			if err != nil {
				return nil, fmt.Errorf("error reading JD dump file '%s': %v", in.JDSQLDumpPath, err)
			}
			// transaction_timeout is a custom RDS instruction, we must replace it
			sqlMigration := strings.Replace(string(d), "SET transaction_timeout = 0;", "", -1)
			sqlCommands = append(sqlCommands, sqlMigration)
			sqlCommands = append(sqlCommands, "DELETE FROM public.csa_keypairs where id = 1;")
		} else {
			// if we don't have a dump we create an empty DB
			sqlCommands = append(sqlCommands, fmt.Sprintf("CREATE DATABASE \"%s\";", JDDatabase))
		}
	}
	initSQL := strings.Join(sqlCommands, "\n")
	initFile, err := os.CreateTemp("", "init-*.sql")
	if err != nil {
		return nil, err
	}
	if _, err := initFile.WriteString(initSQL); err != nil {
		return nil, err
	}
	if err := initFile.Close(); err != nil {
		return nil, err
	}

	var portToExpose int
	if in.Port != 0 {
		portToExpose = in.Port
	} else {
		portToExpose = ExposedStaticPort
	}

	var o *Output

	ns := os.Getenv(components.K8sNamespaceEnvVar)
	// k8s deployment
	if ns != "" {
		_, err := pods.Run(&pods.Config{
			Namespace: pods.S(ns),
			Pods: []*pods.PodConfig{
				{
					Name:  pods.S(in.Name),
					Image: pods.S(in.Image),
					Ports: []string{fmt.Sprintf("%d:%s", portToExpose, Port)},
					Env: &[]*k8s.EnvVar{
						{
							Name:  pods.S("POSTGRES_USER"),
							Value: pods.S("chainlink"),
						},
						{
							Name:  pods.S("POSTGRES_PASSWORD"),
							Value: pods.S("thispasswordislongenough"),
						},
						{
							Name:  pods.S("POSTGRES_DB"),
							Value: pods.S("chainlink"),
						},
					},
					Limits: pods.ResourcesMedium(),
					ContainerSecurityContext: &k8s.SecurityContext{
						RunAsUser:  pods.I(999),
						RunAsGroup: pods.I(999),
					},
					PodSecurityContext: &k8s.PodSecurityContext{
						FsGroup: pods.I(999),
					},
					ConfigMap: map[string]*string{
						"init.sql": pods.S(initSQL),
					},
					ConfigMapMountPath: map[string]*string{
						"init.sql": pods.S("/docker-entrypoint-initdb.d/init.sql"),
					},
					VolumeClaimTemplates: pods.SizedVolumeClaim(pods.S("4Gi")),
				},
			},
		})
		if err != nil {
			return nil, err
		}
		o = &Output{
			ContainerName: containerName,
			InternalURL: fmt.Sprintf(
				"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
				User,
				Password,
				fmt.Sprintf("%s-svc", in.Name),
				// use svc internally too
				portToExpose,
				Database,
			),
			Url: fmt.Sprintf(
				"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
				User,
				Password,
				fmt.Sprintf("%s-svc", in.Name),
				portToExpose,
				Database,
			),
		}
	} else {
		// local deployment
		req := testcontainers.ContainerRequest{
			AlwaysPullImage: in.PullImage,
			Image:           in.Image,
			Name:            containerName,
			Labels:          framework.DefaultTCLabels(),
			ExposedPorts:    []string{bindPort},
			Networks:        []string{framework.DefaultNetworkName},
			NetworkAliases: map[string][]string{
				framework.DefaultNetworkName: {containerName},
			},
			Env: map[string]string{
				"POSTGRES_USER":     User,
				"POSTGRES_PASSWORD": Password,
				"POSTGRES_DB":       Database,
			},
			Cmd: []string{
				"postgres", "-c",
				fmt.Sprintf("port=%s", Port),
				"-c", "shared_preload_libraries=pg_stat_statements",
				"-c", "pg_stat_statements.track=all",
			},
			Files: []testcontainers.ContainerFile{
				{
					HostFilePath:      initFile.Name(),
					ContainerFilePath: "/docker-entrypoint-initdb.d/init.sql",
					FileMode:          0o644,
				},
			},
			Mounts: testcontainers.ContainerMounts{
				{
					Source: testcontainers.GenericVolumeMountSource{
						Name: fmt.Sprintf("%s%s", DBVolumeName, in.VolumeName),
					},
					Target: "/var/lib/postgresql/data",
				},
			},
			WaitingFor: tcwait.ForExec([]string{
				"psql", "-h", "127.0.0.1",
				"-U", User, "-p", Port, "-c", "select", "1", "-d", Database,
			}).
				WithStartupTimeout(3 * time.Minute).
				WithPollInterval(200 * time.Millisecond),
		}
		req.HostConfigModifier = func(h *container.HostConfig) {
			h.PortBindings = nat.PortMap{
				nat.Port(bindPort): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: strconv.Itoa(portToExpose),
					},
				},
			}
			framework.ResourceLimitsFunc(h, in.ContainerResources)
		}
		c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Reuse:            true,
		})
		if err != nil {
			return nil, err
		}
		host, err := framework.GetHostWithContext(ctx, c)
		if err != nil {
			return nil, err
		}
		o = &Output{
			ContainerName: containerName,
			InternalURL: fmt.Sprintf(
				"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
				User,
				Password,
				containerName,
				Port,
				Database,
			),
			Url: fmt.Sprintf(
				"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
				User,
				Password,
				host,
				portToExpose,
				Database,
			),
		}
		if in.JDDatabase {
			o.JDInternalURL = fmt.Sprintf(
				"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
				User,
				Password,
				containerName,
				Port,
				JDDatabase,
			)
			o.JDUrl = fmt.Sprintf(
				"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
				User,
				Password,
				host,
				portToExpose,
				JDDatabase,
			)
		}
	}
	return o, nil
}
