package blockchain

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain/canton"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/pods"
)

const (
	// DefaultCantonPort is a default Canton container port
	DefaultCantonPort = "8080"
	// TokenExpiry is JWT token expiry
	TokenExpiry = time.Hour * 24 * 365 * 10 // 10 years
)

type CantonEndpoints struct {
	// ScanAPIURL https://docs.sync.global/app_dev/scan_api/index.html
	ScanAPIURL string `toml:"scan_api_url" comment:"https://docs.sync.global/app_dev/scan_api/index.html"`
	// RegistryAPIURL https://docs.sync.global/app_dev/token_standard/index.html#api-references
	RegistryAPIURL string `toml:"registry_api_url" comment:"https://docs.sync.global/app_dev/token_standard/index.html#api-references"`

	// SuperValidator The endpoints for the super validator
	SuperValidator CantonParticipantEndpoints `toml:"super_validator" comment:"Canton network super validator"`
	// Participants The endpoints for the participants, in order from participant1 to participantN - depending on the number of validators requested
	Participants []CantonParticipantEndpoints `toml:"participants" comment:"Canton participant endpoints"`
}

type CantonParticipantEndpoints struct {
	// JSONLedgerAPIURL https://docs.digitalasset.com/build/3.5/reference/json-api/json-api.html
	JSONLedgerAPIURL string `toml:"json_ledger_api_url" comment:"https://docs.digitalasset.com/build/3.5/reference/json-api/json-api.html"`
	// GRPCLedgerAPIURL https://docs.digitalasset.com/build/3.5/reference/lapi-proto-docs.html
	GRPCLedgerAPIURL string `toml:"grpc_ledger_api_url" comment:"https://docs.digitalasset.com/build/3.5/reference/lapi-proto-docs.html"`
	// AdminAPIURL https://docs.digitalasset.com/operate/3.5/howtos/configure/apis/admin_api.html
	AdminAPIURL string `toml:"admin_api_url" comment:"https://docs.digitalasset.com/operate/3.5/howtos/configure/apis/admin_api.html"`
	// ValidatorAPIURL https://docs.sync.global/app_dev/validator_api/index.html
	ValidatorAPIURL string `toml:"validator_api_url" comment:"https://docs.sync.global/app_dev/validator_api/index.html"`

	// HTTPHealthCheckURL responds on GET /health
	HTTPHealthCheckURL string `toml:"http_health_check_url" comment:"HTTP health check endpoint, responds on GET /health"`
	// GRPCHealthCheckURL grpc.health.v1.Health/Check
	GRPCHealthCheckURL string `toml:"grpc_health_check_url" comment:"GRPC health check endpoint, responds to grpc.health.v1.Health/Check"`

	// JWT JSON Web Token for this participant
	JWT string `toml:"jwt" comment:"JSON Web Token for this participant"`
}

// newCanton sets up a Canton blockchain network with the specified number of validators.
// It creates a Docker network and starts the necessary containers for Postgres, Canton, Splice, and an Nginx reverse proxy.
//
// Startup timeout: note spinning up a Canton network can take several minutes due to the initialization of  the Splice service.
// tests utilizing this function should set an appropriate timeout to accommodate for this. CTF will time out after 1 hour by default.
//
// The reverse proxy is used to allow access to all validator participants through a single HTTP endpoint.
// The following routes are configured for each participant and the Super Validator (SV):
//   - http://[PARTICIPANT].json-ledger-api.localhost:[PORT] 	-> JSON Ledger API		=> https://docs.digitalasset.com/build/3.3/reference/json-api/json-api.html
//   - grpc://[PARTICIPANT].grpc-ledger-api.localhost:[PORT] 	-> gRPC Ledger API		=> https://docs.digitalasset.com/build/3.3/reference/lapi-proto-docs.html
//   - grpc://[PARTICIPANT].admin-api.localhost:[PORT] 			-> gRPC Admin API		=> https://docs.digitalasset.com/operate/3.5/howtos/configure/apis/admin_api.html
//   - http://[PARTICIPANT].validator-api.localhost:[PORT] 		-> Validator API		=> https://docs.sync.global/app_dev/validator_api/index.html
//   - http://[PARTICIPANT].http-health-check.localhost:[PORT] 	-> HTTP Health Check	=> responds on GET /health
//   - grpc://[PARTICIPANT].grpc-health-check.localhost:[PORT] 	-> gRPC Health Check	=> grpc.health.v1.Health/Check
//
// To access a participant's endpoints, replace [PARTICIPANT] with the participant's identifier, i.e. `sv`, `participant1`, `participant2`, ...
//
// Additionally, the global Scan service is accessible via:
//   - http://scan.localhost:[PORT]/api/scan 					-> Scan API				=> https://docs.sync.global/app_dev/scan_api/index.html
//   - http://scan.localhost:[PORT]/registry 					-> Token Standard API	=> https://docs.sync.global/app_dev/token_standard/index.html#api-references
//
// The PORT is the same for all routes and is specified in the input parameters, defaulting to 8080.
//
// Note: The maximum number of validators supported is 99, participants are numbered starting from `participant1` through `participant99`.
func newCanton(ctx context.Context, in *Input) (*Output, error) {
	if in.NumberOfCantonValidators >= 100 {
		return nil, fmt.Errorf("number of validators too high: %d, valid range is 0-99", in.NumberOfCantonValidators)
	}
	if in.Port == "" {
		in.Port = DefaultCantonPort
	}

	if pods.K8sEnabled() {
		return nil, fmt.Errorf("K8s support is not yet implemented")
	}

	// Set up Postgres container
	postgresReq := canton.PostgresContainerRequest(in.NumberOfCantonValidators)
	_, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: postgresReq,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	// Set up Canton container
	cantonReq := canton.ContainerRequest(in.NumberOfCantonValidators, in.Image, postgresReq.Name)
	_, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: cantonReq,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	// Set up Splice container
	spliceReq := canton.SpliceContainerRequest(in.NumberOfCantonValidators, in.Image, postgresReq.Name, cantonReq.Name)
	_, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: spliceReq,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	// Set up Nginx container
	nginxReq := canton.NginxContainerRequest(in.NumberOfCantonValidators, in.Port, cantonReq.Name, spliceReq.Name)
	nginxContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: nginxReq,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := nginxContainer.Host(ctx)
	if err != nil {
		return nil, err
	}

	svToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "",
		Subject:   "user-sv",
		Audience:  []string{canton.AuthProviderAudience},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiry)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        "",
	}).SignedString([]byte(canton.AuthProviderSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to create token for sv: %w", err)
	}
	endpoints := &CantonEndpoints{
		ScanAPIURL:     fmt.Sprintf("http://scan.%s:%s/api/scan", host, in.Port),
		RegistryAPIURL: fmt.Sprintf("http://scan.%s:%s/registry", host, in.Port),
		SuperValidator: CantonParticipantEndpoints{
			JSONLedgerAPIURL:   fmt.Sprintf("http://sv.json-ledger-api.%s:%s", host, in.Port),
			GRPCLedgerAPIURL:   fmt.Sprintf("sv.grpc-ledger-api.%s:%s", host, in.Port),
			AdminAPIURL:        fmt.Sprintf("sv.admin-api.%s:%s", host, in.Port),
			ValidatorAPIURL:    fmt.Sprintf("http://sv.validator-api.%s:%s/api/validator", host, in.Port),
			HTTPHealthCheckURL: fmt.Sprintf("http://sv.http-health-check.%s:%s", host, in.Port),
			GRPCHealthCheckURL: fmt.Sprintf("sv.grpc-health-check.%s:%s", host, in.Port),
			JWT:                svToken,
		},
		Participants: nil,
	}
	for i := 1; i <= in.NumberOfCantonValidators; i++ {
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Issuer:    "",
			Subject:   fmt.Sprintf("user-participant%v", i),
			Audience:  []string{canton.AuthProviderAudience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiry)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        "",
		}).SignedString([]byte(canton.AuthProviderSecret))
		if err != nil {
			return nil, fmt.Errorf("failed to create token for participant%v: %w", i, err)
		}
		participantEndpoints := CantonParticipantEndpoints{
			JSONLedgerAPIURL:   fmt.Sprintf("http://participant%d.json-ledger-api.%s:%s", i, host, in.Port),
			GRPCLedgerAPIURL:   fmt.Sprintf("participant%d.grpc-ledger-api.%s:%s", i, host, in.Port),
			AdminAPIURL:        fmt.Sprintf("participant%d.admin-api.%s:%s", i, host, in.Port),
			ValidatorAPIURL:    fmt.Sprintf("http://participant%d.validator-api.%s:%s/api/validator", i, host, in.Port),
			HTTPHealthCheckURL: fmt.Sprintf("http://participant%d.http-health-check.%s:%s", i, host, in.Port),
			GRPCHealthCheckURL: fmt.Sprintf("participant%d.grpc-health-check.%s:%s", i, host, in.Port),
			JWT:                token,
		}
		endpoints.Participants = append(endpoints.Participants, participantEndpoints)
	}

	return &Output{
		UseCache:      false,
		Type:          in.Type,
		Family:        FamilyCanton,
		ChainID:       in.ChainID,
		ContainerName: nginxReq.Name,
		NetworkSpecificData: &NetworkSpecificData{
			CantonEndpoints: endpoints,
		},
	}, nil
}
