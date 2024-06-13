package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type Input struct {
	Product           string   `json:"product"`
	TestRegex         string   `json:"test_regex"`
	File              string   `json:"file"`
	EthImplementation string   `json:"eth_implementation"`
	DockerImages      []string `json:"docker_images"`
	ChainIDs          []int    `json:"chain_ids"`
	Networks          string   `json:"networks"`
}

type OutputEntry struct {
	Name              string `json:"name"`
	OS                string `json:"os"`
	Product           string `json:"product"`
	EthImplementation string `json:"eth_implementation"`
	DockerImage       string `json:"docker_image"`
	Run               string `json:"run"`
	Networks          string `json:"networks,omitempty"`
}

func main() {
	var ciTestId, outputFile, product, regex, file, ethImplementation, dockerImagesArg, nodeLabel, chainIDsArg, networksArg string
	var helpArg bool

	var rootCmd = &cobra.Command{
		Use:   "main",
		Short: "This script builds a JSON file with the tests to be run for a given product and Ethereum implementation",
		Long: `This script builds a JSON file with the tests to be run for a given product and Ethereum implementation. The JSON file can be used to generate a test matrix for CI.

When -c/chain_ids flag is provided, the output docker image will have following format: <chain_id>=<docker_image>.
For example if run with -c=1337,2337 and -d=hyperledger/besu:21.0.0 we would get: 1337=hyperledger/besu:21.0.0,2337=hyperledger/besu:21.0.0. This is useful for CCIP, which requires at least two chains per test.

The -w/networks flag is optional and can be used to specify the networks to run the tests on. If not provided, 'networks' will be omitted from the output.
`,
		Run: func(cmd *cobra.Command, args []string) {
			if helpArg {
				fmt.Println(cmd.Long)
				_ = cmd.Usage()
				return
			}
			if ciTestId == "" || outputFile == "" || product == "" || regex == "" || file == "" || ethImplementation == "" || dockerImagesArg == "" {
				fmt.Println(cmd.Short)
				_ = cmd.Usage()
				panic("All flags are required")
			}

			if nodeLabel == "" {
				nodeLabel = "ubuntu-latest"
				fmt.Fprintf(os.Stderr, "Node label not provided, using default: %s\n", nodeLabel)
			}

			dockerImages := strings.Split(dockerImagesArg, ",")

			var chainIDs []int
			if chainIDsArg != "" {
				for _, idStr := range strings.Split(chainIDsArg, ",") {
					id, err := strconv.Atoi(idStr)
					if err != nil {
						panic(fmt.Errorf("invalid chain ID: %s", idStr))
					}
					chainIDs = append(chainIDs, id)
				}
			}

			var networks string
			if networksArg != "" {
				spl := strings.Split(networksArg, ",")
				networks = strings.Join(spl, ",")
			}

			input := Input{
				Product:           product,
				TestRegex:         regex,
				File:              file,
				EthImplementation: ethImplementation,
				DockerImages:      dockerImages,
				ChainIDs:          chainIDs,
				Networks:          networks,
			}

			validateInput(input)

			var output []OutputEntry
			var file *os.File
			var err error
			counter := 1

			if _, err = os.Stat(outputFile); err == nil {
				file, err = os.OpenFile(outputFile, os.O_RDWR, 0644)
				if err != nil {
					panic(fmt.Errorf("error opening file: %v", err))
				}
				defer func() { _ = file.Close() }()

				bytes, err := io.ReadAll(file)
				if err != nil {
					panic(fmt.Errorf("error reading file: %v", err))
				}

				if len(bytes) > 0 {
					if err := json.Unmarshal(bytes, &output); err != nil {
						panic(fmt.Errorf("error unmarshalling JSON: %v", err))
					}
					counter = len(output) + 1
				}
			} else {
				file, err = os.Create(outputFile)
				if err != nil {
					panic(fmt.Errorf("error creating file: %v", err))
				}
			}
			defer func() { _ = file.Close() }()

			for _, image := range dockerImages {
				if !strings.Contains(image, ":") {
					panic(fmt.Errorf("docker image format is invalid: %s", image))
				}

				entry := OutputEntry{
					Name:              fmt.Sprintf("%s-%02d", ciTestId, counter),
					OS:                nodeLabel,
					Product:           input.Product,
					EthImplementation: input.EthImplementation,
					DockerImage:       parseImage(image, chainIDs),
					Run:               fmt.Sprintf("-run '%s' %s", input.TestRegex, input.File),
					Networks:          networks,
				}
				output = append(output, entry)
				counter++

			}

			newOutput, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				panic(fmt.Errorf("error marshalling JSON: %v", err))
			}

			if _, err := file.WriteAt(newOutput, 0); err != nil {
				panic(fmt.Errorf("error writing to file: %v", err))
			}

			fmt.Printf("%d test(s) for %s and %s added successfully!\n", len(dockerImages), input.Product, input.EthImplementation)
		},
	}

	rootCmd.Flags().StringVarP(&ciTestId, "ci_test_id", "t", "", "Test id to use in CI")
	rootCmd.Flags().StringVarP(&outputFile, "output_file", "o", "", "Output file name")
	rootCmd.Flags().StringVarP(&product, "product", "p", "", "Product")
	rootCmd.Flags().StringVarP(&regex, "regex", "r", "", "Test regex")
	rootCmd.Flags().StringVarP(&file, "file", "f", "", "File")
	rootCmd.Flags().StringVarP(&ethImplementation, "eth_implementation", "e", "", "Ethereum implementation (e.g. besu, geth, erigon)")
	rootCmd.Flags().StringVarP(&dockerImagesArg, "docker_images", "d", "", "Docker images (comma separated)")
	rootCmd.Flags().StringVarP(&nodeLabel, "node_label", "n", "", "Node label (runner to use)")
	rootCmd.Flags().StringVarP(&chainIDsArg, "chain_ids", "c", "", "Chain IDs (comma separated)")
	rootCmd.Flags().StringVarP(&networksArg, "networks", "w", "", "Networks (comma separated)")
	rootCmd.Flags().BoolVarP(&helpArg, "help", "h", false, "Display help")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseImage(image string, chainIDs []int) string {
	if len(chainIDs) == 0 {
		return image
	}

	var images []string
	for _, id := range chainIDs {
		images = append(images, fmt.Sprintf("%d=%s", id, image))
	}
	return strings.Join(images, ",")
}

func validateInput(input Input) {
	if input.Product == "" {
		panic(fmt.Errorf("parameter 'product' cannot be empty"))
	}
	if input.TestRegex == "" {
		panic(fmt.Errorf("parameter 'test_regex' cannot be empty"))
	}

	if _, err := regexp.Compile(input.TestRegex); err != nil {
		panic(fmt.Errorf("failed to compile regex: %v", err))
	}

	if input.File == "" {
		panic(fmt.Errorf("parameter 'file' cannot be empty"))
	}
	if input.EthImplementation == "" {
		panic(fmt.Errorf("parameter 'eth_implementation' cannot be empty"))
	}
	if len(input.DockerImages) == 0 || (len(input.DockerImages) == 1 && input.DockerImages[0] == "") {
		panic(fmt.Errorf("parameter 'docker_images' cannot be empty"))
	}
}
