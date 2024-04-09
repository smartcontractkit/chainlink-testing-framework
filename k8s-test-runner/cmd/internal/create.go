package internal

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var Create = &cobra.Command{
	Use:   "create-test-image [path-to-test-folder]",
	RunE:  createRunnerImageRunE,
	Short: "Create test image for K8s test runner",
}

func init() {
	Create.Flags().String("image-registry-url", "", "Image registry url (e.g. ECR url)")
	if err := Create.MarkFlagRequired("image-registry-url"); err != nil {
		log.Fatalf("Failed to mark 'image-registry-url' flag as required: %v", err)
	}
	Create.Flags().String("image-name", "", "Image name (e.g. k8s-test-runner-binary)")
	if err := Create.MarkFlagRequired("image-name"); err != nil {
		log.Fatalf("Failed to mark 'image-name' flag as required: %v", err)
	}
	Create.Flags().String("image-tag", "", "Image tag (e.g. mercury-load-tests)")
	if err := Create.MarkFlagRequired("image-tag"); err != nil {
		log.Fatalf("Failed to mark 'image-tag' flag as required: %v", err)
	}
	Create.Flags().String("test-runner-root-dir", "./", "Test runner root directory with default chart and Dockerfile.testbin")
	Create.Flags().String("timeout", "10m", "Timeout for the test binary build and image push")
}

func createRunnerImageRunE(cmd *cobra.Command, args []string) error {
	testPackage, err := filepath.Abs(args[0])
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}
	if _, err := os.Stat(testPackage); os.IsNotExist(err) {
		return fmt.Errorf("folder with the tests does not exist: %s", testPackage)
	}

	timeoutStr := cmd.Flag("timeout").Value.String()
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return fmt.Errorf("error parsing timeout: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rootDir, err := cmd.Flags().GetString("test-runner-root-dir")
	if err != nil {
		return fmt.Errorf("error getting test runner root directory: %v", err)
	}

	// Get the root directory of the k8s runner
	rootDirAbs, err := filepath.Abs(rootDir)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}

	fmt.Printf("Test runner root directory: %s\n", rootDirAbs)

	// Build the test binary
	testBinPath := filepath.Join(rootDirAbs, "testbin")
	err = buildTestBinary(ctx, testPackage, testBinPath)
	if err != nil {
		return err
	}

	fmt.Print("Creating docker image for the test binary..\n")

	imageName, err := cmd.Flags().GetString("image-name")
	if err != nil {
		return fmt.Errorf("error getting image name: %v", err)
	}
	imageTag, err := cmd.Flags().GetString("image-tag")
	if err != nil {
		return fmt.Errorf("error getting image tag: %v", err)
	}

	// #nosec G204
	buildTestBinDockerImageCmd := exec.CommandContext(ctx, "docker", "build", "--platform", "linux/amd64", "-f", "Dockerfile.testbin", "--build-arg", "TEST_BINARY=testbin", "-t", fmt.Sprintf("%s:%s", imageName, imageTag), ".")
	buildTestBinDockerImageCmd.Dir = rootDirAbs

	fmt.Printf("Running command from %s: %s\n", rootDirAbs, buildTestBinDockerImageCmd.String())

	buildTestBinDockerImageOutput, err := buildTestBinDockerImageCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error building test binary image: %v, output: %s", err, buildTestBinDockerImageOutput)
	}

	fmt.Printf("Done. Created docker image %s:%s\n", imageName, imageTag)

	registryURL, err := cmd.Flags().GetString("image-registry-url")
	if err != nil {
		return fmt.Errorf("error getting ECR registry name: %v", err)
	}

	err = pushDockerImageToECR(ctx, "us-west-2", registryURL, imageName, imageTag)
	if err != nil {
		return fmt.Errorf("error pushing docker image to ECR: %v", err)
	}

	return nil
}

func buildTestBinary(ctx context.Context, testPackage, testBinPath string) error {
	fmt.Printf("Creating test binary for %s\n", testPackage)

	// Create the command to build the test binary
	buildTestBinCmd := exec.CommandContext(ctx, "go", "test", "-c", "-o", testBinPath)
	buildTestBinCmd.Dir = testPackage

	// Set environment variables for the build process
	env := os.Environ()
	env = append(env, "GOOS=linux", "CGO_ENABLED=0", "GOARCH=amd64")
	buildTestBinCmd.Env = env

	// Print the command being executed
	fmt.Printf("Running command from %s: %s\n", testPackage, buildTestBinCmd.String())

	// Execute the build command and capture any output or errors
	output, err := buildTestBinCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error building test binary: %v, output: %s", err, string(output))
	}

	fmt.Printf("Done. Created binary %s\n", testBinPath)

	// If successful, return nil
	return nil
}

// pushDockerImageToECR authenticates with AWS ECR, tags, and pushes a Docker image.
func pushDockerImageToECR(ctx context.Context, region, ecrURL, imageName, imageTag string) error {
	// Authenticate Docker with ECR
	cmdGetLoginPassword := exec.CommandContext(ctx, "aws", "ecr", "get-login-password", "--region", region)
	var loginPassword bytes.Buffer
	cmdGetLoginPassword.Stdout = &loginPassword
	if err := cmdGetLoginPassword.Run(); err != nil {
		return fmt.Errorf("failed to get ECR login password: %w", err)
	}

	cmdDockerLogin := exec.CommandContext(ctx, "docker", "login", "--username", "AWS", "--password-stdin", ecrURL)
	cmdDockerLogin.Stdin = &loginPassword
	if output, err := cmdDockerLogin.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to login to Docker with ECR: %s, %w", string(output), err)
	}

	fmt.Printf("Authenticated Docker with ECR: %s\n", ecrURL)

	// Tag the Docker image with ECR registry name
	// #nosec G204
	cmdDockerTag := exec.CommandContext(ctx, "docker", "tag", fmt.Sprintf("%s:%s", imageName, imageTag), fmt.Sprintf("%s/%s:%s", ecrURL, imageName, imageTag))
	if output, err := cmdDockerTag.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to tag Docker image: %s, %w", string(output), err)
	}

	fmt.Printf("Tagged Docker image: %s/%s:%s\n", ecrURL, imageName, imageTag)

	// Push Docker image to ECR

	fmt.Printf("Running command: docker push %s/%s:%s\n", ecrURL, imageName, imageTag)

	// #nosec G204
	cmdDockerPush := exec.CommandContext(ctx, "docker", "push", fmt.Sprintf("%s/%s:%s", ecrURL, imageName, imageTag))
	if output, err := cmdDockerPush.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to push Docker image to ECR: %s, %w", string(output), err)
	}

	fmt.Printf("Pushed Docker image to ECR: %s/%s:%s\n", ecrURL, imageName, imageTag)

	return nil
}
