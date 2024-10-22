package toml

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
)

func TestItSavesStructToTomlAndReadsIt(t *testing.T) {
	type someStruct struct {
		Hello string            `toml:"hello"`
		Yolo  map[string]string `toml:"yolo"`
		Laika *string           `toml:"laika"`
		Mir   *string           `toml:"mir"`
	}
	testStruct := someStruct{
		Hello: "World",
		Yolo: map[string]string{
			"yolo": "swag",
		},
		Mir: ptr.Ptr("Mir"),
	}

	t.Cleanup(func() {
		os.Remove("test/test.toml")
		os.Remove("test")
		os.Remove("test.toml")
	})

	filePath, err := SaveStructAsToml(testStruct, "test", "test")
	require.NoError(t, err, "failed to save struct as toml")

	var readStruct someStruct
	err = OpenTomlFileAsStruct(filePath, &readStruct)
	require.NoError(t, err, "failed to read toml file")

	require.Equal(t, testStruct, readStruct, "read struct is not equal to original struct")

	//save to current folder
	filePath, err = SaveStructAsToml(testStruct, "", "test")
	require.NoError(t, err, "failed to save struct as toml")

	err = OpenTomlFileAsStruct(filePath, &readStruct)
	require.NoError(t, err, "failed to read toml file")

	require.Equal(t, testStruct, readStruct, "read struct is not equal to original struct")
}

func TestItOverridesExistingFile(t *testing.T) {
	type someStruct struct {
		Hello string            `toml:"hello"`
		Yolo  map[string]string `toml:"yolo"`
		Laika *string           `toml:"laika"`
		Mir   *string           `toml:"mir"`
	}
	testStruct := someStruct{
		Hello: "World",
		Yolo: map[string]string{
			"yolo": "swag",
		},
		Mir: ptr.Ptr("Mir"),
	}

	t.Cleanup(func() {
		os.Remove("test.toml")
	})

	err := os.WriteFile("test.toml", []byte("mietek"), 0600)
	require.NoError(t, err, "failed to create test file")

	filePath, err := SaveStructAsToml(testStruct, "", "test")
	require.NoError(t, err, "failed to save struct as toml")

	var readStruct someStruct
	err = OpenTomlFileAsStruct(filePath, &readStruct)
	require.NoError(t, err, "failed to read toml file")

	require.Equal(t, testStruct, readStruct, "read struct is not equal to original struct")
}

func TestItReturnsErrorWhenOpeningNonExistent(t *testing.T) {
	var readStruct struct{}
	err := OpenTomlFileAsStruct("nonexistent.toml", &readStruct)
	require.Error(t, err, "should return error when opening nonexistent file")
	require.Contains(t, err.Error(), "no such file or directory", "should return correct error when opening nonexistent file")
}
