package canton

import (
	"fmt"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	DefaultPostgresImage = "postgres:14"
	DefaultPostgresUser  = "canton"
	DefaultPostgresPass  = "password"
	DefaultPostgresDB    = "canton"
)

// language=bash
const initDbScript = `
#!/usr/bin/env bash

set -Eeo pipefail

function create_database() {
    local database=$1
    echo "    Creating database: '$database'"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
        CREATE DATABASE "$database";
        GRANT ALL PRIVILEGES ON DATABASE "$database" TO $POSTGRES_USER;
EOSQL
}

if [ -n "$POSTGRES_INIT_DATABASES" ]; then
    echo "Creating multiple databases: $POSTGRES_INIT_DATABASES"
    for database in $(echo $POSTGRES_INIT_DATABASES | tr ',' ' '); do
        create_database $database
    done
    echo "All databases created"
fi
`

func PostgresContainerRequest(
	numberOfValidators int,
) testcontainers.ContainerRequest {
	postgresDatabases := []string{
		"sequencer",
		"mediator",
		"scan",
		"sv",
		"participant-sv",
		"validator-sv",
	}
	for i := range numberOfValidators {
		postgresDatabases = append(postgresDatabases, fmt.Sprintf("participant-%d", i+1))
		postgresDatabases = append(postgresDatabases, fmt.Sprintf("validator-%d", i+1))
	}
	postgresContainerName := framework.DefaultTCName("postgres")
	postgresReq := testcontainers.ContainerRequest{
		Image:    DefaultPostgresImage,
		Name:     postgresContainerName,
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {postgresContainerName},
		},
		WaitingFor: wait.ForExec([]string{
			"pg_isready",
			"-U", DefaultPostgresUser,
			"-d", DefaultPostgresDB,
		}),
		Env: map[string]string{
			"POSTGRES_USER":           DefaultPostgresUser,
			"POSTGRES_PASSWORD":       DefaultPostgresPass,
			"POSTGRES_DB":             DefaultPostgresDB,
			"POSTGRES_INIT_DATABASES": strings.Join(postgresDatabases, ","),
		},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(initDbScript),
				ContainerFilePath: "/docker-entrypoint-initdb.d/create-multiple-databases.sh",
				FileMode:          0755,
			},
		},
		Cmd: []string{
			"postgres",
			"-c", "max_connections=2000",
			"-c", "log_statement=all",
		},
	}

	return postgresReq
}
