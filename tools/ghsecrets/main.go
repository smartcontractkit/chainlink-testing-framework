package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	"github.com/spf13/cobra"
)

func main() {
	var filePath string
	var customSecretID string
	var backend string // Backend: GitHub or AWS

	var setCmd = &cobra.Command{
		Use:   "set",
		Short: "Set test secrets in GitHub or AWS",
		Run: func(cmd *cobra.Command, args []string) {
			// Validate file
			if err := validateFile(filePath); err != nil {
				fmt.Println(err)
				return
			}

			var secretID string
			var err error

			if customSecretID != "" {
				secretID = customSecretID
			} else {
				if !isGHInstalled() {
					fmt.Println("GitHub CLI not found. Please go to https://cli.github.com/ and install it to use this tool.")
					return
				}
				secretID, err = generateSecretIDFromGithubUsername()
				if err != nil {
					log.Fatalf("Failed to generate secret ID: %s", err)
				}
			}

			// Set secret according to backend
			switch strings.ToLower(backend) {
			case "github":
				if err := setGitHubSecret(filePath, secretID); err != nil {
					log.Fatalf("Failed to set GitHub secret: %s", err)
				}
			case "aws":
				if err := setAWSSecret(filePath, secretID); err != nil {
					log.Fatalf("Failed to set AWS secret: %s", err)
				}
			default:
				log.Fatalf("Unsupported backend: %s. Valid backends are 'github' or 'aws'.", backend)
			}
		},
	}

	var rootCmd = &cobra.Command{
		Use:   "ghsecrets",
		Short: "A tool for managing GitHub or AWS test secrets",
	}

	rootCmd.AddCommand(setCmd)

	setCmd.PersistentFlags().StringVarP(&filePath, "file", "f", defaultSecretsPath(), "path to file with test secrets")
	setCmd.PersistentFlags().StringVarP(&customSecretID, "secret-id", "s", "", "custom secret ID")
	setCmd.PersistentFlags().StringVarP(&backend, "backend", "b", "aws", "Backend to use for storing secrets. Options: github, aws")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func defaultSecretsPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %s", err)
	}
	return filepath.Join(homeDir, ".testsecrets")
}

func isGHInstalled() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

func validateFile(filePath string) error {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file '%s' does not exist", filePath)
	}
	if info.Size() == 0 {
		return fmt.Errorf("file '%s' is empty", filePath)
	}
	return nil
}

// generateSecretIDFromGithubUsername generates a secret ID based on the GitHub username
func generateSecretIDFromGithubUsername() (string, error) {
	usernameCmd := exec.Command("gh", "api", "user", "--jq", ".login")
	usernameOutput, err := usernameCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %s, output: %s", err, usernameOutput)
	}
	trimmedUsername := strings.TrimSpace(string(usernameOutput))
	secretID := fmt.Sprintf("BASE64_TESTSECRETS_%s", trimmedUsername)
	return strings.ToUpper(secretID), nil
}

// ===========================
// GitHub Secrets Logic
// ===========================
func setGitHubSecret(filePath, secretID string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Base64 encode the file content
	encoded := base64.StdEncoding.EncodeToString(data)

	// Construct the GitHub CLI command to set the secret
	setSecretCmd := exec.Command("gh", "secret", "set", secretID, "--body", encoded)
	setSecretCmd.Stdin = strings.NewReader(encoded)

	// Execute the command to set the secret
	output, err := setSecretCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set secret: %s\nOutput: %s", err, string(output))
	}

	fmt.Printf(
		"Test secret set successfully in GitHub with key: %s\n\n"+
			"To run a GitHub workflow with the test secrets, use the 'test_secrets_override_key' flag.\n"+
			"Example: gh workflow run ${workflow_name} -f test_secrets_override_key=%s\n",
		secretID, secretID,
	)

	return nil
}

// ===========================
// AWS Secrets Manager Logic
// ===========================
func setAWSSecret(filePath, secretID string) error {
	// 1. Read the file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// 2. Base64 encode the file content (or skip if you prefer raw)
	encoded := base64.StdEncoding.EncodeToString(data)

	// 3. Create a new AWS Secrets Manager client
	sm, err := framework.NewAWSSecretsManager(10 * time.Second)
	if err != nil {
		return fmt.Errorf("failed to initialize AWS Secrets Manager: %w", err)
	}

	// 4. Create (or override) the secret
	err = sm.CreateSecret(secretID, encoded, true)
	if err != nil {
		return fmt.Errorf("failed to create (or override) AWS secret: %w", err)
	}

	fmt.Printf("Test secret set successfully in AWS with key: %s\n", secretID)
	return nil
}
