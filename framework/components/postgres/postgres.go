package postgres

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	tcwait "github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
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
	Image              string                        `toml:"image" validate:"required"`
	Port               int                           `toml:"port"`
	Name               string                        `toml:"name"`
	VolumeName         string                        `toml:"volume_name"`
	Databases          int                           `toml:"databases"`
	JDDatabase         bool                          `toml:"jd_database"`
	JDSQLDumpPath      string                        `toml:"jd_sql_dump_path"`
	PullImage          bool                          `toml:"pull_image"`
	ContainerResources *framework.ContainerResources `toml:"resources"`
	Out                *Output                       `toml:"out"`
}

type Output struct {
	Url           string `toml:"url"`
	ContainerName string `toml:"container_name"`
	InternalURL   string `toml:"internal_url"`
	JDUrl         string `toml:"jd_url"`
	JDInternalURL string `toml:"jd_internal_url"`
}

func NewPostgreSQL(in *Input) (*Output, error) {
	ctx := context.Background()

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
				FileMode:          0644,
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
		WaitingFor: tcwait.ForExec([]string{"psql", "-h", "127.0.0.1",
			"-U", User, "-p", Port, "-c", "select", "1", "-d", Database}).
			WithStartupTimeout(40 * time.Second).
			WithPollInterval(200 * time.Millisecond),
	}
	var portToExpose int
	if in.Port != 0 {
		portToExpose = in.Port
	} else {
		portToExpose = ExposedStaticPort
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
	host, err := framework.GetHost(c)
	if err != nil {
		return nil, err
	}
	o := &Output{
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
	return o, nil
}
