package config

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
)

const (
	E2E_TEST_LOG_COLLECT_ENV                          = "E2E_TEST_LOG_COLLECT"
	E2E_TEST_LOGGING_RUN_ID_ENV                       = "E2E_TEST_LOGGING_RUN_ID"
	E2E_TEST_LOG_STREAM_LOG_TARGETS_ENV               = "E2E_TEST_LOG_STREAM_LOG_TARGETS"
	E2E_TEST_LOKI_TENANT_ID_ENV                       = "E2E_TEST_LOKI_TENANT_ID"
	E2E_TEST_LOKI_ENDPOINT_ENV                        = "E2E_TEST_LOKI_ENDPOINT"
	E2E_TEST_LOKI_BASIC_AUTH_ENV                      = "E2E_TEST_LOKI_BASIC_AUTH"
	E2E_TEST_LOKI_BEARER_TOKEN_ENV                    = "E2E_TEST_LOKI_BEARER_TOKEN" // #nosec G101
	E2E_TEST_GRAFANA_BASE_URL_ENV                     = "E2E_TEST_GRAFANA_BASE_URL"
	E2E_TEST_GRAFANA_DASHBOARD_URL_ENV                = "E2E_TEST_GRAFANA_DASHBOARD_URL"
	E2E_TEST_GRAFANA_BEARER_TOKEN_ENV                 = "E2E_TEST_GRAFANA_BEARER_TOKEN" // #nosec G101
	E2E_TEST_PYROSCOPE_ENABLED_ENV                    = "E2E_TEST_PYROSCOPE_ENABLED"
	E2E_TEST_PYROSCOPE_SERVER_URL_ENV                 = "E2E_TEST_PYROSCOPE_SERVER_URL"
	E2E_TEST_PYROSCOPE_KEY_ENV                        = "E2E_TEST_PYROSCOPE_KEY"
	E2E_TEST_PYROSCOPE_ENVIRONMENT_ENV                = "E2E_TEST_PYROSCOPE_ENVIRONMENT"
	E2E_TEST_CHAINLINK_IMAGE_ENV                      = "E2E_TEST_CHAINLINK_IMAGE"
	E2E_TEST_CHAINLINK_VERSION_ENV                    = "E2E_TEST_CHAINLINK_VERSION"
	E2E_TEST_CHAINLINK_POSTGRES_VERSION_ENV           = "E2E_TEST_CHAINLINK_POSTGRES_VERSION"
	E2E_TEST_CHAINLINK_UPGRADE_IMAGE_ENV              = "E2E_TEST_CHAINLINK_UPGRADE_IMAGE"
	E2E_TEST_CHAINLINK_UPGRADE_VERSION_ENV            = "E2E_TEST_CHAINLINK_UPGRADE_VERSION"
	E2E_TEST_SELECTED_NETWORK_ENV                     = `E2E_TEST_SELECTED_NETWORK`
	E2E_TEST_WALLET_KEY_ENV                           = `E2E_TEST_(.+)_WALLET_KEY$`
	E2E_TEST_WALLET_KEYS_ENV                          = `E2E_TEST_(.+)_WALLET_KEY_(\d+)$`
	E2E_TEST_RPC_HTTP_URL_ENV                         = `E2E_TEST_(.+)_RPC_HTTP_URL$`
	E2E_TEST_RPC_HTTP_URLS_ENV                        = `E2E_TEST_(.+)_RPC_HTTP_URL_(\d+)$`
	E2E_TEST_RPC_WS_URL_ENV                           = `E2E_TEST_(.+)_RPC_WS_URL$`
	E2E_TEST_RPC_WS_URLS_ENV                          = `E2E_TEST_(.+)_RPC_WS_URL_(\d+)$`
	E2E_TEST_PRIVATE_ETHEREUM_EXECUTION_LAYER_ENV     = "E2E_TEST_PRIVATE_ETHEREUM_EXECUTION_LAYER"
	E2E_TEST_PRIVATE_ETHEREUM_ETHEREUM_VERSION_ENV    = "E2E_TEST_PRIVATE_ETHEREUM_ETHEREUM_VERSION"
	E2E_TEST_PRIVATE_ETHEREUM_CUSTOM_DOCKER_IMAGE_ENV = "E2E_TEST_PRIVATE_ETHEREUM_CUSTOM_DOCKER_IMAGE"
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

func MustReadEnvVar_Strings(name, sep string) []string {
	value, err := readEnvVarValue(name, String)
	if err != nil {
		panic(err)
	}
	if value == nil {
		return nil
	}
	strVal := value.(string)
	return strings.Split(strVal, sep)
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

// ReadEnvVarSlice_String reads all environment variables matching the specified pattern and returns a slice of strings.
func ReadEnvVarSlice_String(pattern string) []string {
	re := regexp.MustCompile(pattern)
	var values []string

	for _, env := range getSortedEnvs() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			continue
		}
		key, value := pair[0], pair[1]
		if re.MatchString(key) && value != "" {
			values = append(values, utils.MustResolveEnvPlaceholder(value))
		}
	}
	return values
}

// ReadEnvVarGroupedMap combines environment variables into a map where keys map to slices of strings.
// It accepts `singleEnvPattern` for single variables and `groupEnvPattern` for grouped variables.
// Returns a map combining values from both patterns, with single values wrapped in slices.
func ReadEnvVarGroupedMap(singleEnvPattern, groupEnvPattern string) map[string][]string {
	var singleMap map[string]string
	if singleEnvPattern != "" {
		singleMap = readEnvVarSingleMap(singleEnvPattern)
	}
	return mergeMaps(singleMap, readEnvVarGroupedMap(groupEnvPattern))
}

// readEnvVarValue reads an environment variable and returns the value parsed according to the specified type.
func readEnvVarValue(envVarName string, valueType EnvValueType) (interface{}, error) {
	// Get the environment variable value
	value, isSet := os.LookupEnv(envVarName)
	if !isSet {
		return nil, nil
	}
	if isSet && value == "" {
		return nil, nil
	}
	value = utils.MustResolveEnvPlaceholder(value)

	// Parse the value according to the specified type
	switch valueType {
	case Integer:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("error converting value to integer: %v", err)
		}
		return intVal, nil
	case Boolean:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("error converting value to boolean: %v", err)
		}
		return boolVal, nil
	case Float:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("error converting value to float: %v", err)
		}
		return floatVal, nil
	default: // String or unrecognized type
		return value, nil
	}
}

// readEnvVarGroupedMap scans all environment variables, matches them against
// a specified pattern, and returns a map of grouped values based on the pattern.
// The grouping is defined by the first capture group of the regex.
func readEnvVarGroupedMap(pattern string) map[string][]string {
	re := regexp.MustCompile(pattern)
	groupedVars := make(map[string][]string)
	for _, env := range getSortedEnvs() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			continue
		}
		key, value := pair[0], pair[1]
		matches := re.FindStringSubmatch(key)
		if len(matches) > 1 && value != "" {
			group := matches[1] // Use the first capture group for grouping
			groupedVars[group] = append(groupedVars[group], utils.MustResolveEnvPlaceholder(value))
		}
	}
	return groupedVars
}

func readEnvVarSingleMap(pattern string) map[string]string {
	re := regexp.MustCompile(pattern)
	singleVars := make(map[string]string)
	for _, env := range getSortedEnvs() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			continue
		}
		key, value := pair[0], pair[1]
		matches := re.FindStringSubmatch(key)
		if len(matches) > 1 && value != "" {
			group := matches[1] // Use the first capture group for grouping
			singleVars[group] = utils.MustResolveEnvPlaceholder(value)
		}
	}
	return singleVars
}

// Merges a map[string]string with a map[string][]string and returns a new map[string][]string.
// Elements from the single map are inserted at index 0 in the slice of the new map.
func mergeMaps(single map[string]string, multi map[string][]string) map[string][]string {
	newMap := make(map[string][]string)

	// First, copy all elements from the multi map to the new map
	for key, values := range multi {
		newMap[key] = make([]string, len(values))
		copy(newMap[key], values)
	}

	// Next, insert or prepend the elements from the single map
	for key, value := range single {
		if existingValues, exists := newMap[key]; exists {
			// Prepend the value from the single map
			newMap[key] = append([]string{value}, existingValues...)
		} else {
			// Initialize a new slice if the key does not exist
			newMap[key] = []string{value}
		}
	}

	return newMap
}

type EnvValueType int

const (
	String EnvValueType = iota
	Integer
	Boolean
	Float
)

// getSortedEnvs returns a sorted slice of environment variables
func getSortedEnvs() []string {
	envs := os.Environ()
	// Sort environment variables by key
	sort.Slice(envs, func(i, j int) bool {
		return strings.SplitN(envs[i], "=", 2)[0] < strings.SplitN(envs[j], "=", 2)[0]
	})
	return envs
}
