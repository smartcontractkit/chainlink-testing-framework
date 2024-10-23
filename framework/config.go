package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
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
	"time"
)

const (
	DefaultConfigDir = "."
)

const (
	EnvVarTestConfigs       = "CTF_CONFIGS"
	EnvVarLokiStream        = "CTF_LOKI_STREAM"
	EnvVarAWSSecretsManager = "CTF_AWS_SECRETS_MANAGER"
	// EnvVarCI this is a default env variable many CI runners use so code can detect we run in CI
	EnvVarCI = "CI"
)

const (
	OutputFieldNameTOML = "out"
	OutputFieldName     = "Out"
	OverridesFieldName  = "Overrides"
)

var (
	// Secrets is a singleton AWS Secrets Manager
	// Loaded once on start inside Load and is safe to call concurrently
	Secrets *AWSSecretsManager

	DefaultNetworkName string

	AllowedEmptyConfigurationFields = []string{OutputFieldName, OverridesFieldName}
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
		data, err := os.ReadFile(filepath.Join(DefaultConfigDir, path))
		if err != nil {
			return nil, fmt.Errorf("error reading promtailConfig file %s: %w", path, err)
		}
		if L.GetLevel() == zerolog.DebugLevel {
			fmt.Println(string(data))
		}

		decoder := toml.NewDecoder(strings.NewReader(string(data)))
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&config); err != nil {
			var details *toml.StrictMissingError
			if errors.As(err, &details) {
				fmt.Println(details.String())
			}
			return nil, fmt.Errorf("failed to decode TOML config, strict mode: %s", err)
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
		return fmt.Errorf("config validation failed\nwe are using 'go-playground/validator', please read more here: https://github.com/go-playground/validator?tab=readme-ov-file#usage-and-documentation")
	}
	return nil
}

func Load[X any](t *testing.T) (*X, error) {
	input, err := mergeInputs[X]()
	if err != nil {
		return input, err
	}
	if err := applyEnvConfig("", input); err != nil {
		return nil, fmt.Errorf("error overriding config using envconfig: %s", err)
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

// applyEnvConfig recursively processes environment variables for structs, slices, and maps.
func applyEnvConfig(prefix string, input interface{}) error {
	// Get the value of the input
	val := reflect.ValueOf(input)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("input must be a pointer to a struct")
	}

	// Get the element (dereference) of the pointer
	val = val.Elem()

	// Process the current struct with envconfig
	if err := envconfig.Process(prefix, input); err != nil {
		return err
	}

	// Iterate over the fields of the struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)
		// Handle struct fields
		if field.Kind() == reflect.Struct {
			if err := applyEnvConfig(prefix, field.Addr().Interface()); err != nil {
				return fmt.Errorf("failed to process struct field %s: %w", fieldType.Name, err)
			}
		}
		// Slice fields
		if field.Kind() == reflect.Slice {
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				if elem.Kind() == reflect.Struct {
					if err := applyEnvConfig(prefix, elem.Addr().Interface()); err != nil {
						return fmt.Errorf("failed to process slice element %s[%d]: %w", fieldType.Name, j, err)
					}
				}
			}
		}
		// Recursively handle map fields
		if field.Kind() == reflect.Map {
			for _, key := range field.MapKeys() {
				elem := field.MapIndex(key)
				// If the map value is a struct, create a copy to process it
				if elem.Kind() == reflect.Struct {
					// Create a new variable to hold the value
					elemCopy := reflect.New(elem.Type()).Elem()
					elemCopy.Set(elem)
					if err := applyEnvConfig(prefix, elemCopy.Addr().Interface()); err != nil {
						return fmt.Errorf("failed to process map element %s[%v]: %w", fieldType.Name, key, err)
					}
					// Update the map with the processed value
					field.SetMapIndex(key, elemCopy)
				}
			}
		}
	}
	return nil
}

func UseCache() bool {
	return os.Getenv("CTF_USE_CACHED_OUTPUTS") == "true"
}

func getBaseConfigPath() (string, error) {
	configs := os.Getenv("CTF_CONFIGS")
	if configs == "" {
		return "", fmt.Errorf("no %s env var is provided, you should provide at least one test promtailConfig in TOML", EnvVarTestConfigs)
	}
	return strings.Split(configs, ",")[0], nil
}

func Store[T any](cfg *T) error {
	if UseCache() {
		return nil
	}
	baseConfigPath, err := getBaseConfigPath()
	if err != nil {
		return err
	}
	cachedOutName := fmt.Sprintf("%s-cache.toml", strings.Replace(baseConfigPath, ".toml", "", -1))
	L.Info().Str("OutputFile", cachedOutName).Msg("Storing configuration output")
	d, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(DefaultConfigDir, cachedOutName), d, os.ModePerm)
}

// JSONStrDuration is JSON friendly duration that can be parsed from "1h2m0s" Go format
type JSONStrDuration struct {
	time.Duration
}

func (d *JSONStrDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *JSONStrDuration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
