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
	DefaultNginxImage = "nginx:1.27.0"
)

const nginxConfig = `
events {
  worker_connections  64;
}

http {
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

func getNginxTemplate(numberOfValidators int, enableSplice bool) string {
	template := `
# SV
server {
    listen 		8080;
    server_name sv.json-ledger-api.localhost;
    location / {
        proxy_pass http://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_JSON_API_PORT_PREFIX}00;
		add_header Access-Control-Allow-Origin *;
		add_header Access-Control-Allow-Methods 'GET, POST, OPTIONS';
		add_header Access-Control-Allow-Headers 'Origin, Content-Type, Accept';
    }
}

server {
    listen 		8080 http2;
    server_name sv.grpc-ledger-api.localhost;
    location / {
        grpc_pass grpc://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_LEDGER_API_PORT_PREFIX}00;
    }
}

server {
    listen 		8080;
    server_name sv.http-health-check.localhost;
    location / {
        proxy_pass http://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}00;
    }
}

server {
    listen 		8080 http2;
    server_name sv.grpc-health-check.localhost;
    location / {
        grpc_pass grpc://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}00;
    }
}

server {
    listen 		8080 http2;
    server_name sv.admin-api.localhost;
    location / {
        grpc_pass grpc://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_ADMIN_API_PORT_PREFIX}00;
	}
}
`
	if enableSplice {
		template += `
server {
    listen 		8080;
    server_name sv.validator-api.localhost;
    location /api/validator {
        rewrite ^\/(.*) /$1 break;
        proxy_pass http://${SPLICE_CONTAINER_NAME}:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}00/api/validator;
    }
}

server {
	listen 		8080;
	server_name scan.localhost;
	
	location /api/scan {
		rewrite ^\/(.*) /$1 break;
		proxy_pass http://${SPLICE_CONTAINER_NAME}:5012/api/scan;
	}
	location /registry {
		rewrite ^\/(.*) /$1 break;
		proxy_pass http://${SPLICE_CONTAINER_NAME}:5012/registry;
	}
}
`
	}

	// Add additional validators
	for i := 1; i <= numberOfValidators; i++ {
		template += fmt.Sprintf(`
# Participant %[1]d
	server {
		listen      8080;
		server_name participant%[1]d.json-ledger-api.localhost;
		location / {
			proxy_pass http://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_JSON_API_PORT_PREFIX}%02[1]d;
			add_header Access-Control-Allow-Origin *;
			add_header Access-Control-Allow-Methods 'GET, POST, OPTIONS';
			add_header Access-Control-Allow-Headers 'Origin, Content-Type, Accept';
		}
	}
	
	server {
		listen 		8080 http2;
		server_name participant%[1]d.grpc-ledger-api.localhost;
		location / {
			grpc_pass grpc://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_LEDGER_API_PORT_PREFIX}%02[1]d;
		}
	}
	
	server {
		listen 		8080;
		server_name participant%[1]d.http-health-check.localhost;
		location / {
			proxy_pass http://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}%02[1]d;
		}
	}
	
	server {
		listen 		8080 http2;
		server_name participant%[1]d.grpc-health-check.localhost;
		location / {
			grpc_pass grpc://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}%02[1]d;
		}
	}
	
	server {
		listen 		8080 http2;
		server_name participant%[1]d.admin-api.localhost;
		location / {
			grpc_pass grpc://${CANTON_CONTAINER_NAME}:${CANTON_PARTICIPANT_ADMIN_API_PORT_PREFIX}%02[1]d;
		}
	}
`, i)
		if enableSplice {
			template += fmt.Sprintf(`
	server {
		listen 		8080;
		server_name participant%[1]d.validator-api.localhost;
		location /api/validator {
			rewrite ^\/(.*) /$1 break;
			proxy_pass http://${SPLICE_CONTAINER_NAME}:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}%02[1]d/api/validator;
		}
	}
		`, i)
		}
	}

	// Ensure template ends with a newline
	if !strings.HasSuffix(template, "\n") {
		template += "\n"
	}

	return template
}

func NginxContainerRequest(
	numberOfValidators int,
	port string,
	cantonContainerName string,
	spliceContainerName string,
	enableSplice bool,
) testcontainers.ContainerRequest {
	nginxContainerName := framework.DefaultTCName("nginx")
	nginxReq := testcontainers.ContainerRequest{
		Image:    DefaultNginxImage,
		Name:     nginxContainerName,
		Networks: []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {nginxContainerName},
		},
		WaitingFor:   wait.ForHTTP("/readyz").WithStartupTimeout(time.Second * 10),
		ExposedPorts: []string{fmt.Sprintf("%s:8080", port)},
		Env: func() map[string]string {
			env := map[string]string{
				"CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX": DefaultHTTPHealthcheckPortPrefix,
				"CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX": DefaultGRPCHealthcheckPortPrefix,
				"CANTON_PARTICIPANT_JSON_API_PORT_PREFIX":         DefaultParticipantJsonApiPortPrefix,
				"CANTON_PARTICIPANT_ADMIN_API_PORT_PREFIX":        DefaultParticipantAdminApiPortPrefix,
				"CANTON_PARTICIPANT_LEDGER_API_PORT_PREFIX":       DefaultLedgerApiPortPrefix,
				"CANTON_CONTAINER_NAME":                           cantonContainerName,
			}
			// Always set Splice variables to avoid envsubst issues, even if Splice is disabled
			// (they won't be used in the template when enableSplice is false)
			if enableSplice {
				env["SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX"] = DefaultSpliceValidatorAdminApiPortPrefix
				env["SPLICE_CONTAINER_NAME"] = spliceContainerName
			} else {
				env["SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX"] = ""
				env["SPLICE_CONTAINER_NAME"] = ""
			}
			return env
		}(),
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(nginxConfig),
				ContainerFilePath: "/etc/nginx/nginx.conf",
				FileMode:          0755,
			}, {
				Reader:            strings.NewReader(getNginxTemplate(numberOfValidators, enableSplice)),
				ContainerFilePath: "/etc/nginx/templates/participants.conf.template",
				FileMode:          0755,
			},
		},
		Labels: framework.DefaultTCLabels(),
	}

	return nginxReq
}
