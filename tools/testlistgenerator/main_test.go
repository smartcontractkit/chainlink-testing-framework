package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const OutputFile = "test_list.json"

func TestMainFunction(t *testing.T) {
	resetEnv := func() {
		os.Args = os.Args[:1]
		if _, err := os.Stat(OutputFile); err == nil {
			_ = os.Remove(OutputFile)
		}
	}

	t.Run("FileCreationAndWrite", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "-t", "test", "-o", OutputFile, "-p", "ocr", "-r", "TestOCR.*", "-f", "./smoke/ocr_test.go", "-e", "besu", "-d", "hyperledger/besu:21.0.0,hyperledger/besu:22.0.0"}
		require.NotPanics(t, func() { main() })

		require.FileExists(t, OutputFile)
		bytes, err := os.ReadFile(OutputFile)
		require.NoError(t, err)

		var output []OutputEntry
		err = json.Unmarshal(bytes, &output)
		require.NoError(t, err)
		require.Len(t, output, 2)
		require.Equal(t, "ocr", output[0].Product)
		require.Equal(t, "-run 'TestOCR.*' ./smoke/ocr_test.go", output[0].Run)
		require.Equal(t, "ubuntu-latest", output[0].OS)
		require.Equal(t, "besu", output[0].EthImplementation)
		require.Equal(t, "hyperledger/besu:21.0.0", output[0].DockerImage)
	})

	t.Run("FileCreationAndWriteWithChainIDs", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "-t", "test", "-o", OutputFile, "-p", "ocr", "-r", "TestOCR.*", "-f", "./smoke/ocr_test.go", "-e", "besu", "-d", "hyperledger/besu:21.0.0", "-c", "1337,2337", "-n", "ubuntu-latest"}
		require.NotPanics(t, func() { main() })

		require.FileExists(t, OutputFile)
		bytes, err := os.ReadFile(OutputFile)
		require.NoError(t, err)

		var output []OutputEntry
		err = json.Unmarshal(bytes, &output)
		require.NoError(t, err)
		require.Len(t, output, 1)
		require.Equal(t, "ocr", output[0].Product)
		require.Equal(t, "-run 'TestOCR.*' ./smoke/ocr_test.go", output[0].Run)
		require.Equal(t, "ubuntu-latest", output[0].OS)
		require.Equal(t, "besu", output[0].EthImplementation)
		require.Equal(t, "1337=hyperledger/besu:21.0.0,2337=hyperledger/besu:21.0.0", output[0].DockerImage)
	})

	t.Run("FileCreationAndWriteWithNetworksAndNodeLabel", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "-t", "test", "-o", OutputFile, "-p", "ocr", "-r", "TestOCR.*", "-f", "./smoke/ocr_test.go", "-e", "besu", "-d", "hyperledger/besu:21.0.0", "-w", "mainnet,ropsten", "-n", "ubuntu-latest-2core"}
		require.NotPanics(t, func() { main() })

		require.FileExists(t, OutputFile)
		bytes, err := os.ReadFile(OutputFile)
		require.NoError(t, err)

		var output []OutputEntry
		err = json.Unmarshal(bytes, &output)
		require.NoError(t, err)
		require.Len(t, output, 1)
		require.Equal(t, "ocr", output[0].Product)
		require.Equal(t, "-run 'TestOCR.*' ./smoke/ocr_test.go", output[0].Run)
		require.Equal(t, "ubuntu-latest-2core", output[0].OS)
		require.Equal(t, "besu", output[0].EthImplementation)
		require.Equal(t, "hyperledger/besu:21.0.0", output[0].DockerImage)
		require.Equal(t, "mainnet,ropsten", output[0].Networks)
	})

	t.Run("AppendToFile", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "-t", "test", "-o", OutputFile, "-p", "ocr", "-r", "TestOCR.*", "-f", "./smoke/ocr_test.go", "-e", "besu", "-d", "hyperledger/besu:21.0.0,hyperledger/besu:22.0.0"}
		require.NotPanics(t, func() { main() })

		os.Args = []string{"main", "-t", "test", "-o", OutputFile, "-p", "ocr", "-r", "TestOCR.*", "-f", "./smoke/ocr_test.go", "-e", "geth", "-d", "ethereum/client-go:1.10.0", "-n", "ubuntu-latest"}
		require.NotPanics(t, func() { main() })

		require.FileExists(t, OutputFile)
		bytes, err := os.ReadFile(OutputFile)
		require.NoError(t, err)

		var output []OutputEntry
		err = json.Unmarshal(bytes, &output)
		require.NoError(t, err)
		require.Len(t, output, 3)
	})

	t.Run("OverwriteFile", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "-t", "test", "-o", OutputFile, "-p", "ocr", "-r", "TestOCR.*", "-f", "./smoke/ocr_test.go", "-e", "besu", "-d", "hyperledger/besu:21.0.0,hyperledger/besu:22.0.0", "-n", "ubuntu-latest"}
		require.NotPanics(t, func() { main() })

		require.FileExists(t, OutputFile)
		bytes, err := os.ReadFile(OutputFile)
		require.NoError(t, err)
		var initialOutput []OutputEntry
		err = json.Unmarshal(bytes, &initialOutput)
		require.NoError(t, err)
		require.Len(t, initialOutput, 2)

		os.Args = []string{"main", "-t", "test", "-o", OutputFile, "-p", "ocr", "-r", "TestOCR.*", "-f", "./smoke/ocr_test.go", "-e", "besu", "-d", "hyperledger/besu:22.0.0,hyperledger/besu:23.0.0", "-n", "ubuntu-latest"}
		require.NotPanics(t, func() { main() })

		require.FileExists(t, OutputFile)
		bytes, err = os.ReadFile(OutputFile)
		require.NoError(t, err)

		var output []OutputEntry
		err = json.Unmarshal(bytes, &output)
		require.NoError(t, err)
		require.Len(t, output, 4)
		require.Equal(t, "hyperledger/besu:23.0.0", output[3].DockerImage)
	})

	t.Run("MissingArguments", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "arg1", "arg2", "arg3", "arg4"}
		require.Panics(t, func() { main() })
	})

	t.Run("InvalidDockerImageFormat", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "-t", "test", "-o", OutputFile, "-p", "ocr", "-r", "TestOCR.*", "-f", "./smoke/ocr_test.go", "-e", "besu", "-d", "hyperledger/besu", "-n", "ubuntu-latest"}
		require.PanicsWithError(t, fmt.Sprintf("docker image format is invalid: %s", "hyperledger/besu"), func() { main() })
	})

	t.Run("EmptyOutputFileName", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "-t", "test", "-o", "", "-p", "ocr", "-r", "TestOCR.*", "-f", "./smoke/ocr_test.go", "-e", "besu", "-d", "hyperledger/besu:21.0.0", "-n", "ubuntu-latest"}
		require.Panics(t, func() { main() })
	})

	t.Run("EmptyProduct", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "-t", "test", "-o", OutputFile, "-p", "", "-r", "TestOCR.*", "-f", "./smoke/ocr_test.go", "-e", "besu", "-d", "hyperledger/besu:21.0.0", "-n", "ubuntu-latest"}
		require.Panics(t, func() { main() })
	})

	t.Run("EmptyTestRegex", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "-t", "test", "-o", OutputFile, "-p", "ocr", "-r", "", "-f", "./smoke/ocr_test.go", "-e", "besu", "-d", "hyperledger/besu:21.0.0", "-n", "ubuntu-latest"}
		require.Panics(t, func() { main() })
	})

	t.Run("InvalidTestRegex", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "-t", "test", "-o", OutputFile, "-p", "ocr", "-r", "[invalid", "-f", "./smoke/ocr_test.go", "-e", "besu", "-d", "hyperledger/besu:21.0.0", "-n", "ubuntu-latest"}
		require.Panics(t, func() { main() })
	})

	t.Run("EmptyFile", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "-t", "test", "-o", OutputFile, "-p", "ocr", "-r", "TestOCR.*", "-f", "", "-e", "besu", "-d", "hyperledger/besu:21.0.0", "-n", "ubuntu-latest"}
		require.Panics(t, func() { main() })
	})

	t.Run("EmptyEthImplementation", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "-t", "test", "-o", OutputFile, "-p", "ocr", "-r", "TestOCR.*", "-f", "./smoke/ocr_test.go", "-e", "", "-d", "hyperledger/besu:21.0.0", "-n", "ubuntu-latest"}
		require.Panics(t, func() { main() })
	})

	t.Run("EmptyDockerImages", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "-t", "test", "-o", OutputFile, "-p", "ocr", "-r", "TestOCR.*", "-f", "./smoke/ocr_test.go", "-e", "besu", "-d", "", "-n", "ubuntu-latest"}
		require.Panics(t, func() { main() })
	})

	t.Run("ChainIdsNotInteger", func(t *testing.T) {
		resetEnv()
		os.Args = []string{"main", "-t", "test", "-o", OutputFile, "-p", "ocr", "-r", "TestOCR.*", "-f", "./smoke/ocr_test.go", "-e", "besu", "-d", "hyperledger/besu:21.0.0", "-n", "ubuntu-latest", "-c", "2,invalid"}
		require.Panics(t, func() { main() })
	})

	defer func() { resetEnv() }()
}
