package chipingressset

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jhump/protocompile"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/smartcontractkit/chainlink-common/pkg/chipingress/pb"
)

// code copied from: https://github.com/smartcontractkit/atlas/blob/master/chip-cli/config/config.go and https://github.com/smartcontractkit/atlas/blob/master/chip-cli/config/proto_validator.go
// reason: avoid dependency on the private atlas module
func convertToPbSchemas(schemas map[string]*schema, domain string) []*pb.Schema {
	pbSchemas := make([]*pb.Schema, len(schemas))

	for i, schema := range slices.Collect(maps.Values(schemas)) {

		pbReferences := make([]*pb.SchemaReference, len(schema.References))
		for j, reference := range schema.References {
			pbReferences[j] = &pb.SchemaReference{
				Subject: fmt.Sprintf("%s-%s", domain, reference.Entity),
				Name:    reference.Name,
				// Explicitly omit Version, this tells chip-config to use the latest version of the schema for this reference
			}
		}

		pbSchema := &pb.Schema{
			Subject:    fmt.Sprintf("%s-%s", domain, schema.Entity),
			Schema:     schema.SchemaContent,
			References: pbReferences,
		}

		// If the schema has metadata, we need to add pb metadata to the schema
		if schema.Metadata.Stores != nil {

			stores := make(map[string]*pb.Store, len(schema.Metadata.Stores))
			for key, store := range schema.Metadata.Stores {
				stores[key] = &pb.Store{
					Index:     store.Index,
					Partition: store.Partition,
				}
			}

			pbSchema.Metadata = &pb.MetaData{
				Stores: stores,
			}
		}

		pbSchemas[i] = pbSchema
	}

	return pbSchemas
}

type registrationConfig struct {
	Domain  string   `json:"domain"`
	Schemas []schema `json:"schemas"`
}

type schema struct {
	Entity        string            `json:"entity"`
	Path          string            `json:"path"`
	References    []schemaReference `json:"references,omitempty"`
	SchemaContent string
	Metadata      metadata `json:"metadata,omitempty"`
}

type metadata struct {
	Stores map[string]store `json:"stores"`
}

type store struct {
	Index     []string `json:"index"`
	Partition []string `json:"partition"`
}

type schemaReference struct {
	Name   string `json:"name"`
	Entity string `json:"entity"`
	Path   string `json:"path"`
}

func parseSchemaConfig(configFilePath, schemaDir string) (*registrationConfig, map[string]*schema, error) {
	cfg, err := readConfig(configFilePath)
	if err != nil {
		return nil, nil, err
	}

	if err := validateEntityNames(cfg, schemaDir); err != nil {
		return nil, nil, fmt.Errorf("entity name validation failed: %w", err)
	}

	// Our end goal is to generate a schema registration request to chip config
	// We will use a map to store the schemas by entity and path
	// this is because more than one schema may reference the same schema
	// technically, since SR is idempotent, this is not strictly necessary, as duplicate registrations are noop
	schemas := make(map[string]*schema)

	for _, s := range cfg.Schemas {
		// For each of the schemas, we need to get the references schema content
		for _, reference := range s.References {

			// read schema contents
			refSchemaContent, err := os.ReadFile(path.Join(schemaDir, reference.Path))
			if err != nil {
				return nil, nil, fmt.Errorf("error reading schema: %v", err)
			}

			// generate key with entity and path since other schemas may also reference this schema
			key := fmt.Sprintf("%s:%s", reference.Entity, reference.Path)

			// if the schema already exists, skip it
			if _, ok := schemas[key]; ok {
				continue
			}

			schemas[key] = &schema{
				Entity:        reference.Entity,
				Path:          reference.Path,
				SchemaContent: string(refSchemaContent),
			}
		}

		// add the root schema to the map
		schemaContent, err := os.ReadFile(path.Join(schemaDir, s.Path))
		if err != nil {
			return nil, nil, fmt.Errorf("error reading schema: %v", err)
		}

		key := fmt.Sprintf("%s:%s", s.Entity, s.Path)
		// if the schema already exists, that means it is referenced by another schema.
		// so we just need to add the references to the existing schema in the map
		if s, ok := schemas[key]; ok {
			s.References = append(s.References, s.References...)
			continue
		}

		schemas[key] = &schema{
			Entity:        s.Entity,
			Path:          s.Path,
			SchemaContent: string(schemaContent),
			References:    s.References,
		}

	}

	return cfg, schemas, nil
}

func readConfig(path string) (*registrationConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file '%s': %w", path, err)
	}
	defer f.Close()

	var cfg registrationConfig

	dErr := json.NewDecoder(f).Decode(&cfg)
	if dErr != nil {
		return nil, fmt.Errorf("failed to decode config: %w", dErr)
	}

	return &cfg, nil
}

// validateEntityNames validates that all entity names in the config match the fully qualified
// protobuf names (package.MessageName) from their corresponding proto files.
// It collects all validation errors and returns them together for better user experience.
func validateEntityNames(cfg *registrationConfig, schemaDir string) error {
	var errors []string

	for _, schema := range cfg.Schemas {
		if err := validateEntityName(schema.Entity, schema.Path, schemaDir); err != nil {
			errors = append(errors, fmt.Sprintf("  - schema '%s': %s", schema.Path, err))
		}

		for _, ref := range schema.References {
			if err := validateEntityName(ref.Entity, ref.Path, schemaDir); err != nil {
				errors = append(errors, fmt.Sprintf("  - referenced schema '%s': %s", ref.Path, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("entity name validation failed with %d error(s):\n%s", len(errors), strings.Join(errors, "\n"))
	}

	return nil
}

func validateEntityName(entityName, protoPath, schemaDir string) error {
	fullPath := path.Join(schemaDir, protoPath)

	// Find the message descriptor that matches the entity name
	msgDesc, err := findMessageDescriptor(fullPath, entityName)
	if err != nil {
		return fmt.Errorf("failed to find message descriptor in '%s': %w", protoPath, err)
	}

	// Extract the expected entity name from the message descriptor
	expectedEntity := string(msgDesc.FullName())
	if entityName != expectedEntity {
		return fmt.Errorf(
			"entity name mismatch in chip.json:\n"+
				"  Proto file: %s\n"+
				"  Expected:   %s\n"+
				"  Got:        %s\n"+
				"  \n"+
				"  The entity name must be the fully qualified protobuf name: {package}.{MessageName}",
			protoPath,
			expectedEntity,
			entityName,
		)
	}

	return nil
}

// findMessageDescriptor finds a message descriptor by name (either full name or short name)
// This matches the logic in chip-ingress/internal/serde/message.go
func findMessageDescriptor(filePath, targetMessageName string) (protoreflect.MessageDescriptor, error) {
	compiler := protocompile.Compiler{
		Resolver: &protocompile.SourceResolver{
			ImportPaths: getImportPaths(filePath, 3),
		},
	}

	filename := filepath.Base(filePath)
	fds, err := compiler.Compile(context.Background(), filename)
	if err != nil {
		return nil, fmt.Errorf("failed to compile proto file: %w", err)
	}

	if len(fds) == 0 {
		return nil, fmt.Errorf("no file descriptors found")
	}

	// Search through all file descriptors for the target message
	for _, fd := range fds {
		messages := fd.Messages()
		for i := range messages.Len() {
			msgDesc := messages.Get(i)

			// Match by full name (e.g., "package.MessageName") or short name (e.g., "MessageName")
			if string(msgDesc.FullName()) == targetMessageName || string(msgDesc.Name()) == targetMessageName {
				return msgDesc, nil
			}
		}
	}

	return nil, fmt.Errorf("message descriptor not found for name: %s", targetMessageName)
}

func getImportPaths(path string, depth int) []string {
	paths := make([]string, 0, depth+1)
	paths = append(paths, filepath.Dir(path))

	currentPath := path
	for i := 0; i < depth; i++ {
		currentPath = filepath.Dir(currentPath)
		paths = append(paths, currentPath)
	}
	return paths
}
