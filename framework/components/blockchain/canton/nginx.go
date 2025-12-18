package canton

import (
	"fmt"
	"strings"

	"github.com/testcontainers/testcontainers-go"

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
}
`

func getNginxTemplate(numberOfValidators int) string {
	template := `
# SV
server {
    listen      8080;
    server_name sv.json-ledger-api.localhost;
    location / {
        proxy_pass http://canton:${CANTON_PARTICIPANT_JSON_API_PORT_PREFIX}00;
		add_header Access-Control-Allow-Origin *;
		add_header Access-Control-Allow-Methods 'GET, POST, OPTIONS';
		add_header Access-Control-Allow-Headers 'Origin, Content-Type, Accept';
    }
}

server {
    listen      8080 http2;
    server_name sv.grpc-ledger-api.localhost;
    location / {
        grpc_pass grpc://canton:${CANTON_PARTICIPANT_LEDGER_API_PORT_PREFIX}00;
    }
}

server {
    listen 8080;
    server_name sv.http-health-check.localhost;
    location / {
        proxy_pass http://canton:${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}00;
    }
}

server {
    listen 8080;
    server_name sv.grpc-health-check.localhost;
    location / {
        proxy_pass http://canton:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}00;
    }
}

server {
    listen 8080 http2;
    server_name sv.admin-api.localhost;
    location / {
        grpc_pass grpc://canton:${CANTON_PARTICIPANT_ADMIN_API_PORT_PREFIX}00;
    }
}

server {
    listen 8080;
    server_name sv.wallet.localhost;
    location /api/validator {
        rewrite ^\/(.*) /$1 break;
        proxy_pass http://splice:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}00/api/validator;
    }
}

server {
	listen 8080;
	server_name scan.localhost;
	
	location /api/scan {
		rewrite ^\/(.*) /$1 break;
		proxy_pass http://splice:5012/api/scan;
	}
	location /registry {
		rewrite ^\/(.*) /$1 break;
		proxy_pass http://splice:5012/registry;
	}
}
	`

	// Add additional validators
	for i := range numberOfValidators {
		i += 1 // start from 1 since SV is 0
		template += fmt.Sprintf(`
# Participant %[1]d
	server {
		listen      8080;
		server_name participant%[1]d.json-ledger-api.localhost;
		location / {
			proxy_pass http://canton:${CANTON_PARTICIPANT_JSON_API_PORT_PREFIX}%02[1]d;
			add_header Access-Control-Allow-Origin *;
			add_header Access-Control-Allow-Methods 'GET, POST, OPTIONS';
			add_header Access-Control-Allow-Headers 'Origin, Content-Type, Accept';
		}
	}
	
	server {
		listen      8080 http2;
		server_name participant%[1]d.grpc-ledger-api.localhost;
		location / {
			grpc_pass grpc://canton:${CANTON_PARTICIPANT_LEDGER_API_PORT_PREFIX}%02[1]d;
		}
	}
	
	server {
		listen 8080;
		server_name participant%[1]d.http-health-check.localhost;
		location / {
			proxy_pass http://canton:${CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX}%02[1]d;
		}
	}
	
	server {
		listen 8080;
		server_name participant%[1]d.grpc-health-check.localhost;
		location / {
			proxy_pass http://canton:${CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX}%02[1]d;
		}
	}
	
	server {
		listen 8080 http2;
		server_name participant%[1]d.admin-api.localhost;
		location / {
			grpc_pass grpc://canton:${CANTON_PARTICIPANT_ADMIN_API_PORT_PREFIX}%02[1]d;
		}
	}
	
	server {
		listen 8080;
		server_name participant%[1]d.wallet.localhost;
		location /api/validator {
			rewrite ^\/(.*) /$1 break;
			proxy_pass http://splice:${SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX}%02[1]d/api/validator;
		}
	}
		`, i)
	}

	return template
}

func NginxContainerRequest(
	networkName string,
	numberOfValidators int,
	port string,
) testcontainers.ContainerRequest {
	nginxContainerName := framework.DefaultTCName("nginx")
	nginxReq := testcontainers.ContainerRequest{
		Image:    DefaultNginxImage,
		Name:     nginxContainerName,
		Networks: []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {"nginx"},
		},
		ExposedPorts: []string{fmt.Sprintf("%s:8080", port)},
		Env: map[string]string{
			"CANTON_PARTICIPANT_HTTP_HEALTHCHECK_PORT_PREFIX": DefaultHTTPHealthcheckPortPrefix,
			"CANTON_PARTICIPANT_GRPC_HEALTHCHECK_PORT_PREFIX": DefaultGRPCHealthcheckPortPrefix,
			"CANTON_PARTICIPANT_JSON_API_PORT_PREFIX":         DefaultParticipantJsonApiPortPrefix,
			"CANTON_PARTICIPANT_ADMIN_API_PORT_PREFIX":        DefaultParticipantAdminApiPortPrefix,
			"CANTON_PARTICIPANT_LEDGER_API_PORT_PREFIX":       DefaultLedgerApiPortPrefix,
			"SPLICE_VALIDATOR_ADMIN_API_PORT_PREFIX":          DefaultSpliceValidatorAdminApiPortPrefix,
		},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(nginxConfig),
				ContainerFilePath: "/etc/nginx/nginx.conf",
				FileMode:          0755,
			}, {
				Reader:            strings.NewReader(getNginxTemplate(numberOfValidators)),
				ContainerFilePath: "/etc/nginx/templates/participants.conf.template",
				FileMode:          0755,
			},
		},
	}

	return nginxReq
}
