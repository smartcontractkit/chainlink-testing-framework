package canton

import (
	"fmt"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// Canton Defaults
const (
	SpliceVersion = "0.5.3"
	Image         = "ghcr.io/digital-asset/decentralized-canton-sync/docker/canton"
)

// JWT Auth defaults
const (
	AuthProviderAudience = "https://chain.link"
	AuthProviderSecret   = "unsafe"
)

// Port prefixes for participants
const (
	DefaultParticipantJsonApiPortPrefix      = "11"
	DefaultParticipantAdminApiPortPrefix     = "12"
	DefaultLedgerApiPortPrefix               = "13"
	DefaultHTTPHealthcheckPortPrefix         = "15"
	DefaultGRPCHealthcheckPortPrefix         = "16"
	DefaultSpliceValidatorAdminApiPortPrefix = "22"
)

func getCantonHealthCheckScript(numberOfValidators int) string {
	script := `
#!/bin/bash
# Copyright (c) 2024 Digital Asset (Switzerland) GmbH and/or its affiliates. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

set -eou pipefail

# SV
echo "Checking ${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}00"
grpcurl -plaintext "localhost:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}00" grpc.health.v1.Health/Check

	`
	// Add additional participants
	for i := 1; i <= numberOfValidators; i++ {
		script += fmt.Sprintf(`
# Participant %02[1]d
echo "Checking ${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}%02[1]d"
grpcurl -plaintext "localhost:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}%02[1]d" grpc.health.v1.Health/Check
`, i)
	}

	return script
}

func getCantonConfig(numberOfValidators int) string {
	//language=hocon
	config := `
# re-used storage config block
_storage {
  type = postgres
  config {
     dataSourceClass = "org.postgresql.ds.PGSimpleDataSource"
     properties = {
       serverName = ${?DB_SERVER}
       portNumber = 5432
       databaseName = participant
       currentSchema = participant
       user =  ${?DB_USER}
       password = ${?DB_PASS}
       tcpKeepAlive = true
     }
   }
   parameters {
     max-connections = 32
     migrate-and-start = true
   }
 }

canton {
  features {
    enable-preview-commands = yes
    enable-testing-commands = yes
  }
  parameters {
    manual-start = no
    non-standard-config = yes
    # Bumping because our topology state can get very large due to
    # a large number of participants.
    timeouts.processing.verify-active = 40.seconds
    timeouts.processing.slow-future-warn = 20.seconds
  }

  # Bumping because our topology state can get very large due to
  # a large number of participants.
  monitoring.logging.delay-logging-threshold = 40.seconds
}


_participant {
  init {
    generate-topology-transactions-and-keys = false
    identity.type = manual
  }

  monitoring.grpc-health-server {
    address = "0.0.0.0"
    port = 5061
  }
  storage = ${_storage}

  admin-api {
    address = "0.0.0.0"
    port = 5002
  }

  init.ledger-api.max-deduplication-duration = 30s

  ledger-api {
    # TODO(DACH-NY/canton-network-internal#2347) Revisit this; we want to avoid users to have to set an exp field in their tokens
    max-token-lifetime = Inf
    # Required for pruning
    admin-token-config.admin-claim=true
    address = "0.0.0.0"
    port = 5001

    # We need to bump this because we run one stream per user +
    # polling for domain connections which can add up quite a bit
    # once you're around ~100 users.
    rate-limit.max-api-services-queue-size = 80000
    interactive-submission-service {
      enable-verbose-hashing = true
    }
  }

  http-ledger-api {
    port = 7575
    address = 0.0.0.0
    path-prefix = ${?CANTON_PARTICIPANT_JSON_API_SERVER_PATH_PREFIX}
  }

  parameters {
    initial-protocol-version = 34
    # tune the synchronisation protocols contract store cache
    caching {
      contract-store {
        maximum-size = 1000 # default 1e6
        expire-after-access = 120s # default 10 minutes
      }
    }
    # Bump ACS pruning interval to make sure ACS snapshots are available for longer
    journal-garbage-collection-delay = 24h
  }

  # TODO(DACH-NY/canton-network-node#8331) Tune cache sizes
  # from https://docs.daml.com/2.8.0/canton/usermanual/performance.html#configuration
  # tune caching configs of the ledger api server
  ledger-api {
    index-service {
      max-contract-state-cache-size = 1000 # default 1e4
      max-contract-key-state-cache-size = 1000 # default 1e4

      # The in-memory fan-out will serve the transaction streams from memory as they are finalized, rather than
      # using the database. Therefore, you should choose this buffer to be large enough such that the likeliness of
      # applications having to stream transactions from the database is low. Generally, having a 10s buffer is
      # sensible. Therefore, if you expect e.g. a throughput of 20 tx/s, then setting this number to 200 is sensible.
      # The default setting assumes 100 tx/s.
      max-transactions-in-memory-fan-out-buffer-size = 200 # default 1000
    }
    # Restrict the command submission rate (mainly for SV participants, since they are granted unlimited traffic)
    command-service.max-commands-in-flight = 30 # default = 256
  }

  monitoring.http-health-server {
    address="0.0.0.0"
    port=7000
  }

  topology.broadcast-batch-size = 1
}

# Sequencer
canton.sequencers.sequencer {
  init {
    generate-topology-transactions-and-keys = false
    identity.type = manual
  }

  storage = ${_storage} {
    config.properties {
      databaseName = "sequencer"
      currentSchema = "sequencer"
    }
  }

  public-api {
    address = "0.0.0.0"
    port = 5008
  }

  admin-api {
    address = "0.0.0.0"
    port = 5009
  }

  monitoring.grpc-health-server {
    address = "0.0.0.0"
    port = 5062
  }

  sequencer {
    config {
      storage = ${_storage} {
        config.properties {
          databaseName = "sequencer"
          currentSchema = "sequencer_driver"
        }
      }
    }
    type = reference
  }
}

# Mediator
canton.mediators.mediator {
  init {
    generate-topology-transactions-and-keys = false
    identity.type = manual
  }

  storage = ${_storage} {
    config.properties {
      databaseName = "mediator"
      currentSchema = "mediator"
    }
  }

  admin-api {
    address = "0.0.0.0"
    port = 5007
  }

  monitoring.grpc-health-server {
    address = "0.0.0.0"
    port = 5061
  }
}

################
# Participants #
################

# SV
canton.participants.sv = ${_participant} {
  storage.config.properties.databaseName = participant-sv
  monitoring {
    http-health-server.port = ${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}00
    grpc-health-server.port= ${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}00
  }
  http-ledger-api.port = ${CANTON_PARTICIPANT_JSON_API_PORT_PREFIX}00
  admin-api.port = ${CANTON_PARTICIPANT_ADMIN_API_PORT_PREFIX}00

  ledger-api{
    port = ${CANTON_PARTICIPANT_LEDGER_API_PORT_PREFIX}00
    auth-services = [{
      type = unsafe-jwt-hmac-256
      target-audience = ${API_AUDIENCE}
      secret = ${API_SECRET}
    }]

    user-management-service.additional-admin-user-id = "user-sv"
  }
}

	`

	// Add additional participants
	for i := 1; i <= numberOfValidators; i++ {
		config += fmt.Sprintf(`
# Participant %02[1]d
canton.participants.participant%[1]d = ${_participant} {
  storage.config.properties.databaseName = participant-%[1]d
  monitoring {
    http-health-server.port = ${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}%02[1]d
    grpc-health-server.port= ${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}%02[1]d
  }
  http-ledger-api.port = ${CANTON_PARTICIPANT_JSON_API_PORT_PREFIX}%02[1]d
  admin-api.port = ${CANTON_PARTICIPANT_ADMIN_API_PORT_PREFIX}%02[1]d

  ledger-api{
    port = ${CANTON_PARTICIPANT_LEDGER_API_PORT_PREFIX}%02[1]d
    auth-services = [{
      type = unsafe-jwt-hmac-256
      target-audience = ${API_AUDIENCE}
      secret = ${API_SECRET}
    }]

    user-management-service.additional-admin-user-id = "user-participant%[1]d"
  }
}

		`, i)
	}

	return config
}

func ContainerRequest(
	numberOfValidators int,
	spliceVersion string, // optional, will default to SpliceVersion if empty
	postgresContainerName string,
) testcontainers.ContainerRequest {
	if spliceVersion == "" {
		spliceVersion = SpliceVersion
	}
	cantonContainerName := framework.DefaultTCName("canton")
	cantonReq := testcontainers.ContainerRequest{
		Image:    fmt.Sprintf("%s:%s", Image, spliceVersion),
		Name:     cantonContainerName,
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {cantonContainerName},
		},
		WaitingFor: wait.ForExec([]string{
			"/bin/bash",
			"/app/health-check.sh",
		}).WithStartupTimeout(time.Minute * 5),
		Env: map[string]string{
			"DB_SERVER": postgresContainerName,
			"DB_USER":   DefaultPostgresUser,
			"DB_PASS":   DefaultPostgresPass,

			"API_AUDIENCE": AuthProviderAudience,
			"API_SECRET":   AuthProviderSecret,

			"CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX": DefaultHTTPHealthcheckPortPrefix,
			"CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX": DefaultGRPCHealthcheckPortPrefix,
			"CANTON_PARTICIPANT_JSON_API_PORT_PREFIX":         DefaultParticipantJsonApiPortPrefix,
			"CANTON_PARTICIPANT_ADMIN_API_PORT_PREFIX":        DefaultParticipantAdminApiPortPrefix,
			"CANTON_PARTICIPANT_LEDGER_API_PORT_PREFIX":       DefaultLedgerApiPortPrefix,
		},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(getCantonHealthCheckScript(numberOfValidators)),
				ContainerFilePath: "/app/health-check.sh",
				FileMode:          0755,
			}, {
				Reader:            strings.NewReader(getCantonConfig(numberOfValidators)),
				ContainerFilePath: "/app/app.conf",
				FileMode:          0755,
			},
		},
		Labels: framework.DefaultTCLabels(),
	}

	return cantonReq
}
