package internal

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
	Create.MarkFlagRequired("image-registry-url")
	Create.Flags().String("image-tag", "", "Test name (e.g. mercury-load-test)")
	Create.MarkFlagRequired("image-registry-url")
}

var (
	K8sTestRunnerImageName = "k8s-test-runner"
)

func createRunnerImageRunE(cmd *cobra.Command, args []string) error {
	testPackage, err := filepath.Abs(args[0])
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}
	if _, err := os.Stat(testPackage); os.IsNotExist(err) {
		return fmt.Errorf("folder with the tests does not exist: %s", testPackage)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	// Get the file path of the current file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("error getting filepath of the current file")
	}
	// Get the root directory of the k8s runner
	k8sRunnerDir := filepath.Clean(filepath.Join(filepath.Dir(filename), "../../"))

	// Build the test binary
	testBinPath := filepath.Join(k8sRunnerDir, "testbin")
	err = buildTestBinary(ctx, testPackage, testBinPath)
	if err != nil {
		return err
	}

	// Execute ls -l command
	lsCmd := exec.Command("ls", "-lh", testBinPath)
	lsOutput, err := lsCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing ls -l: %v", err)
	}
	fmt.Print(string(lsOutput))

	imageTag, err := cmd.Flags().GetString("image-tag")
	if err != nil {
		return fmt.Errorf("error getting test name: %v", err)
	}

	fmt.Print("Creating docker image for the test binary..\n")

	buildTestBinDockerImageCmd := exec.CommandContext(ctx, "docker", "build", "--platform", "linux/amd64", "-f", "Dockerfile.testbin", "--build-arg", "TEST_BINARY=testbin", "-t", fmt.Sprintf("%s:%s", K8sTestRunnerImageName, imageTag), ".")
	buildTestBinDockerImageCmd.Dir = k8sRunnerDir

	fmt.Printf("Running command from %s: %s\n", k8sRunnerDir, buildTestBinDockerImageCmd.String())

	buildTestBinDockerImageOutput, err := buildTestBinDockerImageCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error building test binary image: %v, output: %s", err, buildTestBinDockerImageOutput)
	}

	fmt.Printf("Done. Created docker image %s:%s\n", K8sTestRunnerImageName, imageTag)

	registryURL, err := cmd.Flags().GetString("image-registry-url")
	if err != nil {
		return fmt.Errorf("error getting ECR registry name: %v", err)
	}

	err = pushDockerImageToECR(ctx, "us-west-2", registryURL, K8sTestRunnerImageName, imageTag)
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
	cmdDockerTag := exec.CommandContext(ctx, "docker", "tag", fmt.Sprintf("%s:%s", imageName, imageTag), fmt.Sprintf("%s/%s:%s", ecrURL, imageName, imageTag))
	if output, err := cmdDockerTag.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to tag Docker image: %s, %w", string(output), err)
	}

	fmt.Printf("Tagged Docker image: %s/%s:%s\n", ecrURL, imageName, imageTag)

	// Push Docker image to ECR

	fmt.Printf("Running command: docker push %s/%s:%s\n", ecrURL, imageName, imageTag)

	cmdDockerPush := exec.CommandContext(ctx, "docker", "push", fmt.Sprintf("%s/%s:%s", ecrURL, imageName, imageTag))
	if output, err := cmdDockerPush.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to push Docker image to ECR: %s, %w", string(output), err)
	}

	fmt.Printf("Pushed Docker image to ECR: %s/%s:%s\n", ecrURL, imageName, imageTag)

	return nil
}
