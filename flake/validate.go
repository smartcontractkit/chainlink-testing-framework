package flake

import (
	"embed"
	"errors"
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed flaky_test_schema.json
var content embed.FS

func ValidateFileAgainstSchema(flakyTestFile string) error {
	documentLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", flakyTestFile))
	schema, err := content.ReadFile("flaky_test_schema.json")
	if err != nil {
		return err
	}
	schemaLoader := gojsonschema.NewStringLoader(string(schema))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		// Joining the slice of error details into a single string
		var errorDetails []string
		errorDetails = append(errorDetails, "The document is not valid. See errors:")
		for _, desc := range result.Errors() {
			errorDetails = append(errorDetails, "- "+desc.String())
		}
		return errors.New(strings.Join(errorDetails, "\n"))
	}
	return nil
}
