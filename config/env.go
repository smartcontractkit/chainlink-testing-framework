package config

import "github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"

const (
	E2E_TEST_LOKI_TENANT_ID_ENV          = "E2E_TEST_LOKI_TENANT_ID"
	E2E_TEST_LOKI_ENDPOINT_ENV           = "E2E_TEST_LOKI_ENDPOINT"
	E2E_TEST_LOKI_BASIC_AUTH_ENV         = "E2E_TEST_LOKI_BASIC_AUTH"
	E2E_TEST_LOKI_BEARER_TOKEN_ENV       = "E2E_TEST_LOKI_BEARER_TOKEN"
	E2E_TEST_GRAFANA_BASE_URL_ENV        = "E2E_TEST_GRAFANA_BASE_URL"
	E2E_TEST_GRAFANA_DASHBOARD_URL_ENV   = "E2E_TEST_GRAFANA_DASHBOARD_URL"
	E2E_TEST_GRAFANA_BEARER_TOKEN_ENV    = "E2E_TEST_GRAFANA_BEARER_TOKEN"
	E2E_TEST_PYROSCOPE_ENABLED_ENV       = "E2E_TEST_PYROSCOPE_ENABLED"
	E2E_TEST_PYROSCOPE_SERVER_URL_ENV    = "E2E_TEST_PYROSCOPE_SERVER_URL"
	E2E_TEST_PYROSCOPE_KEY_ENV           = "E2E_TEST_PYROSCOPE_KEY"
	E2E_TEST_PYROSCOPE_ENVIRONMENT_ENV   = "E2E_TEST_PYROSCOPE_ENVIRONMENT"
	E2E_TEST_CHAINLINK_IMAGE_ENV         = "E2E_TEST_CHAINLINK_IMAGE"
	E2E_TEST_CHAINLINK_UPGRADE_IMAGE_ENV = "E2E_TEST_CHAINLINK_UPGRADE_IMAGE"
	E2E_TEST_WALLET_KEY_ENV              = `E2E_TEST_(.+)_WALLET_KEY$`
	E2E_TEST_WALLET_KEYS_ENV             = `E2E_TEST_(.+)_WALLET_KEY_(\d+)$`
	E2E_TEST_RPC_HTTP_URL_ENV            = `E2E_TEST_(.+)_RPC_HTTP_URL$`
	E2E_TEST_RPC_HTTP_URLS_ENV           = `E2E_TEST_(.+)_RPC_HTTP_URL_(\d+)$`
	E2E_TEST_RPC_WS_URL_ENV              = `E2E_TEST_(.+)_RPC_WS_URL$`
	E2E_TEST_RPC_WS_URLS_ENV             = `E2E_TEST_(.+)_RPC_WS_URL_(\d+)$`
)

func MustReadEnvVar_String(name string) string {
	value, err := readEnvVarValue(name, String)
	if err != nil {
		panic(err)
	}
	if value == nil {
		return ""
	}
	return value.(string)
}

func MustReadEnvVar_Boolean(name string) *bool {
	value, err := readEnvVarValue(name, Boolean)
	if err != nil {
		panic(err)
	}
	if value == nil {
		return nil
	}
	return ptr.Ptr(value.(bool))
}

func ReadEnvVarGroupedMap(singleEnvPattern, groupEnvPattern string) map[string][]string {
	return mergeMaps(loadEnvVarSingleMap(singleEnvPattern), loadEnvVarGroupedMap(groupEnvPattern))
}
