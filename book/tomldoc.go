package main

import (
	"flag"
	"log"
	"os"
	"reflect"

	toml "github.com/pelletier/go-toml/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

// AllStructDocs is just an example that is unmarshalled to TOML
// with all "comment" tags so we can have a good configuration examples
// structs are initialized with initializeStruct
// when you add a new component just add a new type here
// and docs will be automatically generated
type AllStructDocs struct {
	Fake        *fake.Input              `toml:"fakes" comment:"Fake HTTP server representing 3rd party dependencies"`
	Blockchains []*blockchain.Input      `toml:"blockchains" comment:"Various blockchains, see 'type' field comments"`
	NodeSets    []*simple_node_set.Input `toml:"nodesets" comment:"Chainlink Node Set including multiple Chainlink nodes forming a DON"`
}

// initializeAny initializes any object so we can call Unmarshal method for all the fields
func initializeAny(obj any) {
	v := reflect.ValueOf(obj).Elem()
	// Check if it's a struct before we iterate the fields
	if v.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		// Skip unexported fields
		if !fieldType.IsExported() {
			continue
		}
		switch field.Kind() {
		case reflect.Ptr:
			if field.IsNil() {
				elemType := field.Type().Elem()
				// Check if it's a pointer to a struct
				if elemType.Kind() == reflect.Struct {
					newVal := reflect.New(elemType)
					field.Set(newVal)
					initializeAny(newVal.Interface())
				} else {
					// For pointers to basic types, just create zero value
					field.Set(reflect.New(elemType))
				}
			} else {
				// Already initialized, but might point to a struct
				if field.Elem().Kind() == reflect.Struct {
					initializeAny(field.Interface())
				}
			}

		case reflect.Slice:
			if field.Len() == 0 {
				elemType := field.Type().Elem()

				// For slice of pointers to structs
				if elemType.Kind() == reflect.Ptr {
					ptrElemType := elemType.Elem()
					if ptrElemType.Kind() == reflect.Struct {
						newElem := reflect.New(ptrElemType)
						// Initialize the struct element
						initializeAny(newElem.Interface())
						// Create slice with one element
						slice := reflect.MakeSlice(field.Type(), 1, 1)
						slice.Index(0).Set(newElem)
						field.Set(slice)
					}
				} else if elemType.Kind() == reflect.Struct {
					// For slice of structs (not pointers)
					newElem := reflect.New(elemType).Elem()
					initializeAny(newElem.Addr().Interface())
					slice := reflect.MakeSlice(field.Type(), 1, 1)
					slice.Index(0).Set(newElem)
					field.Set(slice)
				}
				// For slices of basic types ([]string, []int, etc.) do nothing
			}

		case reflect.Struct:
			// For nested structs (not pointers)
			initializeAny(field.Addr().Interface())
		case reflect.Map:
			if field.IsNil() {
				field.Set(reflect.MakeMap(field.Type()))
			}
			// For all other types (string, int, bool, etc.), do nothing, already initialized
		}
	}
}

func main() {
	outputFile := flag.String("output", "toml-docs.toml", "Output file to write all the struct examples to")
	flag.Parse()
	f, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	// this also writes "comment" lines from fields as # comment on top of each field
	encoder.SetIndentTables(true)

	// initialize all the framework structs
	all := &AllStructDocs{}
	initializeAny(all)

	// encode and write examples
	if err := encoder.Encode(&all); err != nil {
		log.Fatalf("Error encoding to file: %v", err)
	}
}
