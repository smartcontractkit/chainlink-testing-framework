package canton

import (
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	DefaultNginxImage        = "nginx:1.29.5"
	DefaultNginxInternalPort = 8080
)

const nginxConfig = `
events {
  worker_connections  64;
}

http {
	server_names_hash_bucket_size 128;
	include mime.types;
	default_type application/octet-stream;
	client_max_body_size 100M;
	
	# Logging
	log_format json_combined escape=json
	'{'
		'"time_local":"$time_local",'
		'"remote_addr":"$remote_addr",'
		'"remote_user":"$remote_user",'
		'"request":"$request",'
		'"status": "$status",'
		'"body_bytes_sent":"$body_bytes_sent",'
		'"request_time":"$request_time",'
		'"http_referrer":"$http_referer",'
		'"http_user_agent":"$http_user_agent"'
	'}';
	access_log /var/log/nginx/access.log json_combined;
	error_log /var/log/nginx/error.log;
	
	include /etc/nginx/conf.d/participants.conf;
	
	server {
		server_name localhost;
		location /readyz {
			add_header Content-Type text/plain;
			return 200 'OK';
		}
	}
}
`

func getNginxTemplate(nginxContainerName string, nginxInternalPort int, numberOfValidators int) (template string, internalHostnames []string) {
	template = fmt.Sprintf(`
# SV
server {
    listen 		%[1]d;
    server_name sv.json-ledger-api.*;
    location / {
        proxy_pass http://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_JSON_API_PORT_PREFIX}00;
		add_header Access-Control-Allow-Origin *;
		add_header Access-Control-Allow-Methods 'GET, POST, OPTIONS';
		add_header Access-Control-Allow-Headers 'Origin, Content-Type, Accept';
    }
}

server {
    listen 		%[1]d http2;
    server_name sv.grpc-ledger-api.*;
    location / {
        grpc_pass grpc://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_LEDGER_API_PORT_PREFIX}00;
    }
}

server {
    listen 		%[1]d;
    server_name sv.http-health-check.*;
    location / {
        proxy_pass http://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}00;
    }
}

server {
    listen 		%[1]d http2;
    server_name sv.grpc-health-check.*;
    location / {
        grpc_pass grpc://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}00;
    }
}

server {
    listen 		%[1]d http2;
    server_name sv.admin-api.*;
    location / {
        grpc_pass grpc://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_ADMIN_API_PORT_PREFIX}00;
    }
}

server {
    listen 		%[1]d;
    server_name sv.validator-api.*;
    location /api/validator {
        rewrite ^\/(.*) /$1 break;
        proxy_pass http://${SPLICE_CONTAINER_NAME}:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}00/api/validator;
    }
}

server {
	listen 		%[1]d;
	server_name scan.*;
	
	location /api/scan {
		rewrite ^\/(.*) /$1 break;
		proxy_pass http://${SPLICE_CONTAINER_NAME}:5012/api/scan;
	}
	location /registry {
		rewrite ^\/(.*) /$1 break;
		proxy_pass http://${SPLICE_CONTAINER_NAME}:5012/registry;
	}
}
	`, nginxInternalPort)
	internalHostnames = append(internalHostnames,
		fmt.Sprintf("sv.json-ledger-api.%s", nginxContainerName),
		fmt.Sprintf("sv.grpc-ledger-api.%s", nginxContainerName),
		fmt.Sprintf("sv.http-health-check.%s", nginxContainerName),
		fmt.Sprintf("sv.grpc-health-check.%s", nginxContainerName),
		fmt.Sprintf("sv.admin-api.%s", nginxContainerName),
		fmt.Sprintf("sv.validator-api.%s", nginxContainerName),
		fmt.Sprintf("scan.%s", nginxContainerName),
	)

	// Add additional validators
	for i := 1; i <= numberOfValidators; i++ {
		template += fmt.Sprintf(`
# Participant %[2]d
	server {
		listen      %[1]d;
		server_name participant%[2]d.json-ledger-api.*;
		location / {
			proxy_pass http://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_JSON_API_PORT_PREFIX}%02[2]d;
			add_header Access-Control-Allow-Origin *;
			add_header Access-Control-Allow-Methods 'GET, POST, OPTIONS';
			add_header Access-Control-Allow-Headers 'Origin, Content-Type, Accept';
		}
	}
	
	server {
		listen 		%[1]d http2;
		server_name participant%[2]d.grpc-ledger-api.*;
		location / {
			grpc_pass grpc://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_LEDGER_API_PORT_PREFIX}%02[2]d;
		}
	}
	
	server {
		listen 		%[1]d;
		server_name participant%[2]d.http-health-check.*;
		location / {
			proxy_pass http://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}%02[2]d;
		}
	}
	
	server {
		listen 		%[1]d http2;
		server_name participant%[2]d.grpc-health-check.*;
		location / {
			grpc_pass grpc://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}%02[2]d;
		}
	}
	
	server {
		listen 		%[1]d http2;
		server_name participant%[2]d.admin-api.*;
		location / {
			grpc_pass grpc://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_ADMIN_API_PORT_PREFIX}%02[2]d;
		}
	}
	
	server {
		listen 		%[1]d;
		server_name participant%[2]d.validator-api.*;
		location /api/validator {
			rewrite ^\/(.*) /$1 break;
			proxy_pass http://${SPLICE_CONTAINER_NAME}:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}%02[2]d/api/validator;
		}
	}
		`, nginxInternalPort, i)
		internalHostnames = append(internalHostnames,
			fmt.Sprintf("participant%d.json-ledger-api.%s", i, nginxContainerName),
			fmt.Sprintf("participant%d.grpc-ledger-api.%s", i, nginxContainerName),
			fmt.Sprintf("participant%d.http-health-check.%s", i, nginxContainerName),
			fmt.Sprintf("participant%d.grpc-health-check.%s", i, nginxContainerName),
			fmt.Sprintf("participant%d.admin-api.%s", i, nginxContainerName),
			fmt.Sprintf("participant%d.validator-api.%s", i, nginxContainerName),
		)
	}

	return template, internalHostnames
}

func NginxContainerRequest(
	numberOfValidators int,
	port string,
	cantonContainerName string,
	spliceContainerName string,
) (testcontainers.ContainerRequest, string) {
	nginxContainerName := framework.DefaultTCName("canton-nginx")
	// Docker doesn't support DNS wildcards: https://github.com/moby/moby/issues/43442
	// In order to allow for another container to reach the Nginx container under all the defined hostnames,
	// they need to be explicitly set as network aliases.
	nginxTemplate, internalHostnames := getNginxTemplate(nginxContainerName, DefaultNginxInternalPort, numberOfValidators)
	nginxReq := testcontainers.ContainerRequest{
		Image:    DefaultNginxImage,
		Name:     nginxContainerName,
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: append([]string{nginxContainerName}, internalHostnames...),
		},
		WaitingFor:   wait.ForHTTP("/readyz").WithStartupTimeout(time.Second * 10),
		ExposedPorts: []string{fmt.Sprintf("%d/tcp", DefaultNginxInternalPort)},
		HostConfigModifier: func(h *container.HostConfig) {
			containerPort := network.MustParsePort(fmt.Sprintf("%d/tcp", DefaultNginxInternalPort))
			h.PortBindings = network.PortMap{
				containerPort: []network.PortBinding{
					{HostIP: netip.MustParseAddr("0.0.0.0"), HostPort: port},
				},
			}
		},
		Env: map[string]string{
			"CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX": DefaultHTTPHealthcheckPortPrefix,
			"CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX": DefaultGRPCHealthcheckPortPrefix,
			"CANTON_PARTICIPANT_JSON_API_PORT_PREFIX":         DefaultParticipantJsonApiPortPrefix,
			"CANTON_PARTICIPANT_ADMIN_API_PORT_PREFIX":        DefaultParticipantAdminApiPortPrefix,
			"CANTON_PARTICIPANT_LEDGER_API_PORT_PREFIX":       DefaultLedgerApiPortPrefix,
			"SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX":          DefaultSpliceValidatorAdminApiPortPrefix,

			"CANTON_CONTAINER_NAME": cantonContainerName,
			"SPLICE_CONTAINER_NAME": spliceContainerName,
		},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(nginxConfig),
				ContainerFilePath: "/etc/nginx/nginx.conf",
				FileMode:          0755,
			}, {
				Reader:            strings.NewReader(nginxTemplate),
				ContainerFilePath: "/etc/nginx/templates/participants.conf.template",
				FileMode:          0755,
			},
		},
		Labels: framework.DefaultTCLabels(),
	}

	return nginxReq, nginxContainerName
}
