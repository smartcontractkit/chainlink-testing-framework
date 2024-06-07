package config_test

import (
	"encoding/json"
	"fmt"

	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/config"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/ptr"
)

func TestBytesToAnyTomlStruct(t *testing.T) {
	type DeepNestedConfig struct {
		Labels *[]string
	}

	type NestedConfig struct {
		Tags        []string
		Deep        *DeepNestedConfig
		Text        string
		AnotherText *string
	}

	type Value struct {
		Name   *string
		BigInt *big.Int
	}

	type MapEntry struct {
		Name  string
		Value *Value
	}

	type Config struct {
		Name     string
		Numbers  *[]int
		Nested   *NestedConfig
		Map      map[string]MapEntry
		Duration blockchain.StrDuration
	}

	logger := zerolog.New(nil)

	testCases := []struct {
		name           string
		initialConfig  *Config
		tomlContent    []byte
		expectedConfig *Config
	}{
		{
			name: "Partial update",
			initialConfig: &Config{
				Name:    "Original",
				Numbers: &[]int{1, 2, 3},
				Nested: &NestedConfig{
					Tags:        []string{"a", "b", "c"},
					Deep:        &DeepNestedConfig{Labels: &[]string{"label1", "label2"}},
					Text:        "Original Text",
					AnotherText: ptr.Ptr("Original AnotherText"),
				},
				Map: map[string]MapEntry{
					"entry1": {
						Name: "entry12",
						Value: &Value{
							Name:   ptr.Ptr("value1"),
							BigInt: big.NewInt(123),
						},
					},
				},
			},
			tomlContent: []byte(`
Name = "Partially Updated"
[Map.entry1.Value]
BigInt = 99999
`),
			expectedConfig: &Config{
				Name:    "Partially Updated",
				Numbers: &[]int{1, 2, 3}, // unchanged
				Nested: &NestedConfig{
					Tags:        []string{"a", "b", "c"},                                  // unchanged
					Deep:        &DeepNestedConfig{Labels: &[]string{"label1", "label2"}}, // unchanged
					Text:        "Original Text",                                          // unchanged
					AnotherText: ptr.Ptr("Original AnotherText"),                          // unchanged
				},
				Map: map[string]MapEntry{ // who map replaced
					"entry1": {
						Name: "",
						Value: &Value{
							Name:   nil,
							BigInt: big.NewInt(99999),
						},
					},
				},
			},
		},
		{
			name: "Full update",
			initialConfig: &Config{
				Name:     "Original",
				Numbers:  &[]int{1, 2, 3},
				Duration: blockchain.StrDuration{Duration: 1},
				Nested: &NestedConfig{
					Tags:        []string{"a", "b", "c"},
					Deep:        &DeepNestedConfig{Labels: &[]string{"label1", "label2"}},
					Text:        "Original Text",
					AnotherText: ptr.Ptr("Original AnotherText"),
				},
				Map: map[string]MapEntry{
					"entry1": {
						Name: "entry12",
						Value: &Value{
							Name:   ptr.Ptr("value1"),
							BigInt: big.NewInt(123),
						},
					},
				},
			},
			tomlContent: []byte(`
Name = "Fully Updated"
Numbers = [4, 5, 6]
Duration = '5m'
[Nested]
Tags = ["x", "y", "z"]
Text = "Updated Text"
AnotherText = "Updated AnotherText"
[Nested.Deep]
Labels = ["newLabel1", "newLabel2"]
[Map.entry1]
Name = "entry12"
[Map.entry1.Value]
Name = "value updated"
BigInt = 58172
`),
			expectedConfig: &Config{
				Name:     "Fully Updated",
				Numbers:  &[]int{4, 5, 6},
				Duration: blockchain.StrDuration{Duration: time.Duration(5 * time.Minute)},
				Nested: &NestedConfig{
					Tags:        []string{"x", "y", "z"},
					Deep:        &DeepNestedConfig{Labels: &[]string{"newLabel1", "newLabel2"}},
					Text:        "Updated Text",
					AnotherText: ptr.Ptr("Updated AnotherText"),
				},
				Map: map[string]MapEntry{
					"entry1": {
						Name: "entry12",
						Value: &Value{
							Name:   ptr.Ptr("value updated"),
							BigInt: big.NewInt(58172),
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := config.BytesToAnyTomlStruct(logger, "config.toml", "config", tc.initialConfig, tc.tomlContent)
			require.NoError(t, err, "Error unmarshalling TOML content")

			if !reflect.DeepEqual(tc.initialConfig, tc.expectedConfig) {
				actual, _ := json.MarshalIndent(tc.initialConfig, "", "  ")
				expected, _ := json.MarshalIndent(tc.expectedConfig, "", "  ")
				fmt.Printf("Actual:\n%s\n", actual)
				fmt.Printf("Expected:\n%s\n", expected)
				t.Fatal("Expected and actual structs do not match")
			}
		})
	}
}
