package framework

import (
	"bytes"
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-playground/validator/v10"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/network"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"text/template"
)

const (
	DefaultConfigDir = "."
)

const (
	EnvVarTestConfigs       = "CTF_CONFIGS"
	EnvVarLokiStream        = "LOKI_STREAM"
	EnvVarAWSSecretsManager = "AWS_SECRETS_MANAGER"
)

var (
	// Secrets is a singleton AWS Secrets Manager
	// Loaded once on start inside Load and is safe to call concurrently
	Secrets *AWSSecretsManager

	DefaultNetworkName string

	AllowedEmptyConfigurationFields = []string{"Out", "Overrides"}
)

type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

// mergeInputs merges all EnvVarTestConfigs filenames into one files, starting from the last and applying to the first
func mergeInputs[T any]() (*T, error) {
	var config T
	paths := strings.Split(os.Getenv(EnvVarTestConfigs), ",")
	_, err := getBaseConfigPath()
	if err != nil {
		return nil, err
	}
	for _, path := range paths {
		L.Info().Str("Path", path).Msg("Loading configuration input")
		file, err := os.ReadFile(filepath.Join(DefaultConfigDir, path))
		if err != nil {
			return nil, fmt.Errorf("error reading promtailConfig file %s: %w", path, err)
		}
		if L.GetLevel() == zerolog.DebugLevel {
			fmt.Println(string(file))
		}

		err = toml.Unmarshal(file, &config)
		if err != nil {
			return nil, fmt.Errorf("error parsing promtailConfig file %s: %w", path, err)
		}
	}
	if L.GetLevel() == zerolog.DebugLevel {
		L.Debug().Msg("Merged inputs")
		spew.Dump(config)
	}
	return &config, nil
}

// checkRequiredTag recursively checks if "required" exists in any validation tags
func checkRequiredTag(structValue interface{}, parent string) []ValidationError {
	var validationErrors []ValidationError

	v := reflect.ValueOf(structValue)

	// Handle pointer cases (e.g., *ServerConfig, *DatabaseConfig)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			validationErrors = append(validationErrors, ValidationError{
				Field:   parent,
				Message: fmt.Sprintf("field '%s' is a nil pointer", parent),
			})
			return validationErrors
		}
		v = v.Elem()
	}

	// Ensure we're working with a struct
	if v.Kind() != reflect.Struct {
		return validationErrors
	}

	// Check all fields in the struct
outer:
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldValue := v.Field(i)
		fieldName := parent + "." + field.Name

		for _, f := range AllowedEmptyConfigurationFields {
			if strings.Contains(fieldName, f) {
				continue outer
			}
		}

		tomlValue, ok := field.Tag.Lookup("toml")
		if !ok {
			validationErrors = append(validationErrors, ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("field '%s' doesn't have any 'toml' tag value", fieldName),
			})
		}

		// Check if the field has a "validate" tag with "required"
		validateTag, ok := field.Tag.Lookup("validate")
		if !ok || !strings.Contains(validateTag, "required") {
			validationErrors = append(validationErrors, ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("field '%s' is missing 'required' validation tag, TOML key: %s", fieldName, tomlValue),
			})
		}

		// Recurse into structs, slices, maps, and pointers
		// we ignore Chan, Func, Interface or UnsafePointer
		switch fieldValue.Kind() {
		case reflect.Struct:
			validationErrors = append(validationErrors, checkRequiredTag(fieldValue.Interface(), fieldName)...)
		case reflect.Ptr:
			validationErrors = append(validationErrors, checkRequiredTag(fieldValue.Interface(), fieldName)...)
		case reflect.Slice:
			for j := 0; j < fieldValue.Len(); j++ {
				validationErrors = append(validationErrors, checkRequiredTag(fieldValue.Index(j).Interface(), fmt.Sprintf("%s[%d]", fieldName, j))...)
			}
		case reflect.Map:
			for _, key := range fieldValue.MapKeys() {
				validationErrors = append(validationErrors, checkRequiredTag(fieldValue.MapIndex(key).Interface(), fmt.Sprintf("%s[%v]", fieldName, key))...)
			}
		default:
		}
	}

	return validationErrors
}

func noFieldsWithoutRequiredTag(cfg interface{}) []ValidationError {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(cfg)

	var validationErrors []ValidationError
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.StructNamespace(),
				Value:   err.Value(),
				Message: fmt.Sprintf("validation failed on '%s' with tag '%s'", err.Field(), err.Tag()),
			})
		}
	}

	// Combine missing "required" tag errors with validation errors
	validationErrors = append(validationErrors, checkRequiredTag(cfg, "")...)

	return validationErrors
}

func validateStruct(s interface{}) error {
	errs := noFieldsWithoutRequiredTag(s)
	for _, e := range errs {
		L.Error().Any("error", e).Send()
	}
	if len(errs) > 0 {
		return fmt.Errorf("promtailConfig validation failed, exiting")
	}
	return nil
}

func Load[X any](t *testing.T) (*X, error) {
	input, err := mergeInputs[X]()
	if err != nil {
		return input, err
	}
	if err := validateStruct(input); err != nil {
		return nil, err
	}
	t.Cleanup(func() {
		err := Store[X](input)
		require.NoError(t, err)
	})
	// TODO: not all the people have AWS access, sadly enough, uncomment when granted
	//if os.Getenv(EnvVarAWSSecretsManager) == "true" {
	//	Secrets, err = NewAWSSecretsManager(1 * time.Minute)
	//	if err != nil {
	//		return nil, fmt.Errorf("failed to connect AWSSecretsManager: %w", err)
	//	}
	//}
	// TODO: do we really need more than one environment locally? Do we need more than one in CI so network isolation make sense?
	net, err := network.New(
		context.Background(),
		network.WithLabels(map[string]string{"framework": "ctf"}),
	)
	if err != nil {
		return input, err
	}
	DefaultNetworkName = net.Name
	if os.Getenv(EnvVarLokiStream) == "true" {
		err = NewLokiStreamer()
		require.NoError(t, err)
	}
	return input, nil
}

func RenderTemplate(tmpl string, data interface{}) (string, error) {
	var buf bytes.Buffer
	err := template.Must(template.New("tmpl").Parse(tmpl)).Execute(&buf, data)
	return buf.String(), err
}
