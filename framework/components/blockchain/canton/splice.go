package canton

import (
	"fmt"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	SpliceImage = "ghcr.io/digital-asset/decentralized-canton-sync/docker/splice-app"
)

func getSpliceHealthCheckScript(numberOfValidators int) string {
	script := `
#!/bin/bash
# Copyright (c) 2024 Digital Asset (Switzerland) GmbH and/or its affiliates. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

set -eou pipefail

curl -f http://localhost:5012/api/scan/readyz
curl -f http://localhost:5014/api/sv/readyz

# SV
curl -f "http://localhost:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}00/api/validator/readyz"
`
	for i := 1; i <= numberOfValidators; i++ {
		script += fmt.Sprintf(`
# Participant %02[1]d
curl -f "http://localhost:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}%02[1]d/api/validator/readyz"
		`, i)
	}

	return script
}

func getSpliceConfig(numberOfValidators int) string {
	//language=hocon
	config := `
_storage {
  type = postgres
  config {
     dataSourceClass = "org.postgresql.ds.PGSimpleDataSource"
     properties = {
       serverName = ${DB_SERVER}
       portNumber = 5432
       databaseName = validator
       currentSchema = validator
       user =  ${DB_USER}
       password = ${DB_PASS}
       tcpKeepAlive = true
     }
   }
   parameters {
     max-connections = 32
     migrate-and-start = true
   }
 }

_validator_backend {
  latest-packages-only = true
  domain-migration-id = 0
  storage = ${_storage}
  admin-api = {
    address = "0.0.0.0"
    port = 5003
  }
  participant-client = {
    admin-api = {
      address = ${CANTON_CONTAINER_NAME}
      port = 5002
    }
    ledger-api.client-config = {
      address = ${CANTON_CONTAINER_NAME}
      port = 5001
    }
  }
  scan-client {
    type = "bft"
    seed-urls = []
    seed-urls.0 = "http://localhost:5012"
  }

  app-instances {
  }
  onboarding.sv-client.admin-api.url = "http://localhost:5014"
  domains.global.alias = "global"
  contact-point = "contact@local.host"
  canton-identifier-config.participant = participant
}

canton.features.enable-testing-commands = yes

# SV
_sv_participant_client = {
  admin-api {
    address = ${CANTON_CONTAINER_NAME}
    port = ${CANTON_PARTICIPANT_ADMIN_API_PORT_PREFIX}00
  }
  ledger-api {
    client-config {
      address = ${CANTON_CONTAINER_NAME}
      port = ${CANTON_PARTICIPANT_LEDGER_API_PORT_PREFIX}00
    }
    auth-config {
        type = "self-signed"
        user = "user-sv"
        audience = ${API_AUDIENCE}
        secret = ${API_SECRET}
    }
  }
}

_splice-instance-names {
  network-name = "Splice"
  network-favicon-url = "https://www.hyperledger.org/hubfs/hyperledgerfavicon.png"
  amulet-name = "Amulet"
  amulet-name-acronym = "AMT"
  name-service-name = "Amulet Name Service"
  name-service-name-acronym = "ANS"
}

canton {
  scan-apps.scan-app {
    is-first-sv = true
    domain-migration-id = 0
    storage = ${_storage} {
      config.properties {
        databaseName = scan
        currentSchema = scan
      }
    }

    admin-api = {
      address = "0.0.0.0"
      port = 5012
    }
    participant-client = ${_sv_participant_client}
    sequencer-admin-client = {
      address = ${CANTON_CONTAINER_NAME}
      port = 5009
    }
    mediator-admin-client = {
      address = ${CANTON_CONTAINER_NAME}
      port = 5007
    }
    sv-user = "user-sv"
    splice-instance-names = ${_splice-instance-names}
  }

  sv-apps.sv {
    latest-packages-only = true
    domain-migration-id = 0
    expected-validator-onboardings = [
    ]
    scan {
      public-url="http://localhost:5012"
      internal-url="http://localhost:5012"
    }
    local-synchronizer-node {
      sequencer {
        admin-api {
          address = ${CANTON_CONTAINER_NAME}
          port = 5009
        }
        internal-api {
          address = ${CANTON_CONTAINER_NAME}
          port = 5008
        }
        external-public-api-url = "http://"${CANTON_CONTAINER_NAME}":5008"
      }
      mediator.admin-api {
        address = ${CANTON_CONTAINER_NAME}
        port = 5007
      }
    }

    storage = ${_storage} {
      config.properties {
        databaseName = sv
        currentSchema = sv
      }
    }

    admin-api = {
      address = "0.0.0.0"
      port = 5014
    }
    participant-client = ${_sv_participant_client}

    domains {
      global {
        alias = "global"
        url = ${?SPLICE_APP_SV_GLOBAL_DOMAIN_URL}
      }
    }

    auth = {
        algorithm = "hs-256-unsafe"
        audience = ${API_AUDIENCE}
        secret = ${API_SECRET}
    }
    ledger-api-user = "user-sv"
    validator-ledger-api-user = "user-sv"

    automation {
      paused-triggers = [
        "org.lfdecentralizedtrust.splice.sv.automation.delegatebased.ExpiredAmuletTrigger",
        "org.lfdecentralizedtrust.splice.sv.automation.delegatebased.ExpiredLockedAmuletTrigger",
        "org.lfdecentralizedtrust.splice.sv.automation.delegatebased.ExpiredAnsSubscriptionTrigger",
        "org.lfdecentralizedtrust.splice.sv.automation.delegatebased.ExpiredAnsEntryTrigger",
        "org.lfdecentralizedtrust.splice.sv.automation.delegatebased.ExpireTransferPreapprovalsTrigger",
      ]
    }

    onboarding = {
      type = found-dso
      name = sv
      first-sv-reward-weight-bps = 10000
      round-zero-duration = ${?SPLICE_APP_SV_ROUND_ZERO_DURATION}
      initial-tick-duration = ${?SPLICE_APP_SV_INITIAL_TICK_DURATION}
      initial-holding-fee = ${?SPLICE_APP_SV_INITIAL_HOLDING_FEE}
      initial-amulet-price = ${?SPLICE_APP_SV_INITIAL_AMULET_PRICE}
      is-dev-net = true
      public-key = ${?SPLICE_APP_SV_PUBLIC_KEY}
      private-key = ${?SPLICE_APP_SV_PRIVATE_KEY}
      initial-round = ${?SPLICE_APP_SV_INITIAL_ROUND}
    }
    initial-amulet-price-vote = ${?SPLICE_APP_SV_INITIAL_AMULET_PRICE_VOTE}
    comet-bft-config = {
      enabled = false
      enabled = ${?SPLICE_APP_SV_COMETBFT_ENABLED}
      connection-uri = ""
      connection-uri = ${?SPLICE_APP_SV_COMETBFT_CONNECTION_URI}
    }
    contact-point = "contact@local.host"
    canton-identifier-config = {
      participant = sv
      sequencer = sv
      mediator = sv
    }

    splice-instance-names = ${_splice-instance-names}
  }
}

# SV
canton.validator-apps.sv-validator_backend = ${_validator_backend} {
	canton-identifier-config.participant = sv
	onboarding = null
	scan-client = null
	scan-client = {
		type = "trust-single"
		url="http://localhost:5012"
	}
	sv-user="user-sv"
	sv-validator=true
	storage.config.properties.databaseName = validator-sv
	admin-api.port = ${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}00
	participant-client = ${_sv_participant_client}
	auth = {
		algorithm = "hs-256-unsafe"
		audience = ${API_AUDIENCE}
		secret = ${API_SECRET}
	}
	ledger-api-user = "user-sv"
	validator-wallet-users.0 = "sv"
}

`
	// Add additional participants
	for i := 1; i <= numberOfValidators; i++ {
		config += fmt.Sprintf(`
# Participant %02[1]d
canton.validator-apps.participant%[1]d-validator_backend = ${_validator_backend} {
	onboarding.secret = "participant%[1]d-validator-onboarding-secret"
	validator-party-hint = "participant%[1]d-localparty-1"
	domain-migration-dump-path = "/domain-upgrade-dump/domain_migration_dump-participant%[1]d.json"
	storage.config.properties.databaseName = validator-%[1]d
	admin-api.port = ${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}%02[1]d
	participant-client {
		admin-api.port = ${CANTON_PARTICIPANT_ADMIN_API_PORT_PREFIX}%02[1]d

		ledger-api = {
			client-config.port = ${CANTON_PARTICIPANT_LEDGER_API_PORT_PREFIX}%02[1]d
			auth-config = {
				type = "self-signed"
				user = "user-participant%[1]d"
				audience = ${API_AUDIENCE}
				secret = ${API_SECRET}
			}
		}
	}
	auth = {
		algorithm = "hs-256-unsafe"
		audience = ${API_AUDIENCE}
		secret = ${API_SECRET}
	}
	ledger-api-user = "user-participant%[1]d"
	validator-wallet-users.0="participant%[1]d"

	domains.global.buy-extra-traffic {
		min-topup-interval = "1m"
		target-throughput = "20000"
	}
}

canton.sv-apps.sv.expected-validator-onboardings += { secret = "participant%[1]d-validator-onboarding-secret" }
		`, i)
	}

	return config
}

func SpliceContainerRequest(
	numberOfValidators int,
	spliceVersion string,
	postgresContainerName string,
	cantonContainerName string,
) testcontainers.ContainerRequest {
	if spliceVersion == "" {
		spliceVersion = SpliceVersion
	}
	spliceContainerName := framework.DefaultTCName("splice")
	spliceReq := testcontainers.ContainerRequest{
		Image:    fmt.Sprintf("%s:%s", SpliceImage, spliceVersion),
		Name:     spliceContainerName,
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {spliceContainerName},
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
			"SPLICE_APP_VALIDATOR_LEDGER_API_AUTH_AUDIENCE": AuthProviderAudience,

			"CANTON_CONTAINER_NAME": cantonContainerName,

			"CANTON_PARTICIPANT_ADMIN_API_PORT_PREFIX":  DefaultParticipantAdminApiPortPrefix,
			"CANTON_PARTICIPANT_LEDGER_API_PORT_PREFIX": DefaultLedgerApiPortPrefix,
			"SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX":    DefaultSpliceValidatorAdminApiPortPrefix,
		},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(getSpliceHealthCheckScript(numberOfValidators)),
				ContainerFilePath: "/app/health-check.sh",
				FileMode:          0755,
			}, {
				Reader:            strings.NewReader(getSpliceConfig(numberOfValidators)),
				ContainerFilePath: "/app/app.conf",
				FileMode:          0755,
			},
		},
	}

	return spliceReq
}
