package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	var filePath string
	var customSecretID string

	var setCmd = &cobra.Command{
		Use:   "set",
		Short: "Set test secrets in GitHub",
		Run: func(cmd *cobra.Command, args []string) {
			if !isGHInstalled() {
				fmt.Println("GitHub CLI not found. Please go to https://cli.github.com/ and install it to use this tool.")
				return
			}

			if err := validateFile(filePath); err != nil {
				fmt.Println(err)
				return
			}

			secretID, err := getSecretID(customSecretID)
			if err != nil {
				log.Fatalf("Failed to obtain secret ID: %s", err)
			}

			setSecret(filePath, secretID)
		},
	}

	var rootCmd = &cobra.Command{
		Use:   "ghsecrets",
		Short: "A tool for managing GitHub test secrets",
	}

	rootCmd.AddCommand(setCmd)
	setCmd.PersistentFlags().StringVarP(&filePath, "file", "f", defaultSecretsPath(), "path to dotenv file with test secrets")
	setCmd.PersistentFlags().StringVarP(&customSecretID, "secret-id", "s", "", "custom secret ID. Do not use unless you know what you are doing")

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

func getSecretID(customID string) (string, error) {
	if customID != "" {
		return customID, nil
	}
	usernameCmd := exec.Command("gh", "api", "user", "--jq", ".login")
	usernameOutput, err := usernameCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %s, output: %s", err, usernameOutput)
	}
	trimmedUsername := strings.TrimSpace(string(usernameOutput))
	secretID := fmt.Sprintf("BASE64_TESTSECRETS_%s", trimmedUsername)
	return strings.ToUpper(secretID), nil
}

func setSecret(filePath, secretID string) {
	// Read the file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	// Base64 encode the file content
	encoded := base64.StdEncoding.EncodeToString(data)

	// Construct the GitHub CLI command to set the secret
	setSecretCmd := exec.Command("gh", "secret", "set", secretID, "--body", encoded)
	setSecretCmd.Stdin = strings.NewReader(encoded)

	// Execute the command to set the secret
	output, err := setSecretCmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to set secret: %s\nOutput: %s", err, string(output))
	}

	fmt.Printf(
		"Test secret set successfully in Github with key: %s\n\n"+
			"To run a Github workflow with the test secrets, use the 'test_secrets_override_key' flag.\n"+
			"Example: gh workflow run ${workflow_name} -f test_secrets_override_key=%s\n",
		secretID, secretID,
	)
}
