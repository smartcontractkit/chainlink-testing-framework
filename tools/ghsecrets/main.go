package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/cobra"
)

func main() {
	var filePath string
	var secretID string
	var backend string      // Backend: GitHub or AWS
	var decode bool         // Decode flag for `get`
	var profile string      // AWS profile to use
	var sharedWith []string // List of ARNs to share the secret with
	var env string          // Environment name for GitHub secrets

	// Set Command
	var setCmd = &cobra.Command{
		Use:   "set",
		Short: "Set test secrets in GitHub or AWS",
		Run: func(cmd *cobra.Command, args []string) {
			if err := validateFile(filePath); err != nil {
				exitWithError(err, "Failed to validate file")
				return
			}

			if secretID == "" {
				if !isGHInstalled() {
					exitWithError(nil, "GitHub CLI not found. Please go to https://cli.github.com/ and install it to use this tool.")
					return
				}
				var err error
				secretID, err = generateSecretIDFromGithubUsername()
				if err != nil {
					exitWithError(err, "Failed to generate secret ID")
					return
				}
			}

			switch strings.ToLower(backend) {
			case "github":
				if err := setGitHubSecret(filePath, secretID, env); err != nil {
					exitWithError(err, "Failed to set GitHub secret")
					return
				}
			case "aws":
				if profile == "" {
					exitWithError(nil, "AWS profile is required when using the AWS backend. Use the --profile flag to specify it.")
					return
				}
				// Ensure AWS secretID starts with "testsecrets/" prefix
				secretID = ensurePrefix(secretID, "testsecrets/")
				if err := setAWSSecret(filePath, secretID, profile, sharedWith); err != nil {
					exitWithError(err, "Failed to set AWS secret")
					return
				}
			default:
				exitWithError(nil, "Unsupported backend. Valid backends are 'github' or 'aws'.")
				return
			}
		},
	}

	// Get Command
	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Retrieve a secret from AWS Secrets Manager",
		Run: func(cmd *cobra.Command, args []string) {
			if profile == "" {
				exitWithError(nil, "AWS profile is required when using the AWS backend. Use the --profile flag to specify it.")
				return
			}
			secretID = ensurePrefix(secretID, "testsecrets/")
			if err := getAWSSecret(secretID, decode, profile); err != nil {
				exitWithError(err, "Failed to retrieve AWS secret")
			}
		},
	}

	var rootCmd = &cobra.Command{
		Use:   "ghsecrets",
		Short: "A tool for managing GitHub or AWS test secrets",
	}

	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(getCmd)

	setCmd.PersistentFlags().StringVarP(&filePath, "file", "f", defaultSecretsPath(), "Path to file with test secrets")
	setCmd.PersistentFlags().StringVarP(&secretID, "secret-id", "s", "", "ID of the secret to set")
	setCmd.PersistentFlags().StringVarP(&backend, "backend", "b", "aws", "Backend to use for storing secrets. Options: github, aws")
	setCmd.PersistentFlags().StringVar(&profile, "profile", "", "AWS profile to use for credentials (required for AWS backend)")
	setCmd.PersistentFlags().StringSliceVar(&sharedWith, "shared-with", []string{}, "Comma-separated list of IAM ARNs to share the secret with")
	setCmd.PersistentFlags().StringVar(&env, "env", "", "Optional environment name (for GitHub Secrets)")

	getCmd.PersistentFlags().StringVarP(&secretID, "secret-id", "s", "", "ID of the secret to retrieve")
	getCmd.PersistentFlags().BoolVarP(&decode, "decode", "d", true, "Decode the Base64-encoded secret value")
	getCmd.PersistentFlags().StringVar(&profile, "profile", "", "AWS profile to use for credentials (required for AWS backend)")

	getCmd.MarkPersistentFlagRequired("secret-id")

	if err := rootCmd.Execute(); err != nil {
		exitWithError(err, "Failed to execute command")
	}
}

// setGitHubSecret creates or updates a secret in GitHub
func setGitHubSecret(filePath, secretID, env string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(data)

	// Build the gh command
	args := []string{"secret", "set", secretID, "--body", encoded}

	// If --env was provided, add the environment argument
	if env != "" {
		args = append(args, "--env", env)
	}

	setSecretCmd := exec.Command("gh", args...)
	setSecretCmd.Stdin = strings.NewReader(encoded)

	output, err := setSecretCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set secret: %s\nOutput: %s", err, string(output))
	}

	fmt.Printf(
		"Test secret set successfully in GitHub Secrets with key: %s\n\n"+
			"To run a GitHub workflow with the test secrets, use the 'test_secrets_override_key' flag.\n"+
			"Example: gh workflow run ${workflow_name} -f test_secrets_override_key=%s\n",
		secretID, secretID,
	)
	return nil
}

// setAWSSecret creates or updates a secret in AWS Secrets Manager
func setAWSSecret(filePath, secretID, profile string, sharedWith []string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	encoded := base64.StdEncoding.EncodeToString(data)

	cfg, err := loadAWSConfig(profile)
	if err != nil {
		return handleAWSSSOError(err)
	}

	smClient := secretsmanager.NewFromConfig(cfg)

	_, err = smClient.CreateSecret(context.TODO(), &secretsmanager.CreateSecretInput{
		Name:         aws.String(secretID),
		SecretString: aws.String(encoded),
		Description:  aws.String("Chainlink Test Secret created by CTF/ghsecrets CLI"),
	})
	if err != nil {
		var resourceExistsErr *types.ResourceExistsException
		if errors.As(err, &resourceExistsErr) {
			fmt.Printf("Secret %s already exists, updating its value...\n", secretID)
			_, err = smClient.UpdateSecret(context.TODO(), &secretsmanager.UpdateSecretInput{
				SecretId:     aws.String(secretID),
				SecretString: aws.String(encoded),
				Description:  aws.String("Secret updated by ghsecrets CLI"),
			})
			if err != nil {
				return handleAWSSSOError(err)
			}
		} else {
			return handleAWSSSOError(err)
		}
	}

	if len(sharedWith) > 0 {
		err = updateAWSSecretAccessPolicy(secretID, sharedWith, profile)
		if err != nil {
			return fmt.Errorf("failed to update secret sharing policy: %w", err)
		}
	}

	fmt.Printf(
		"Test secret set successfully in AWS Secrets Manager with key: %s\n\n"+
			"To use this secret in a GitHub workflow, set the 'test_secrets_override_key' flag with the 'aws:' prefix. Example:\n"+
			"gh workflow run ${workflow_name} -f test_secrets_override_key=aws:%s\n",
		secretID, secretID,
	)
	return nil
}

// getAWSSecret retrieves a test secret from AWS Secrets Manager
func getAWSSecret(secretID string, decode bool, profile string) error {
	cfg, err := loadAWSConfig(profile)
	if err != nil {
		return handleAWSSSOError(err)
	}

	smClient := secretsmanager.NewFromConfig(cfg)
	out, err := smClient.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	})
	if err != nil {
		return handleAWSSSOError(err)
	}

	value := aws.ToString(out.SecretString)
	if decode {
		decoded, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return fmt.Errorf("failed to decode secret value: %w", err)
		}
		value = string(decoded)
	}

	fmt.Printf("Retrieved secret value:\n%s\n", value)
	return nil
}

// updateAWSSecretAccessPolicy updates the sharing policy for a secret in AWS Secrets Manager
func updateAWSSecretAccessPolicy(secretID string, sharedWith []string, profile string) error {
	// 1) Load AWS config
	cfg, err := loadAWSConfig(profile)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	// 2) Create an STS client to find your AWS Account ID
	stsClient := sts.NewFromConfig(cfg)
	callerIdentity, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return fmt.Errorf("failed to get caller identity: %w", err)
	}
	accountID := aws.ToString(callerIdentity.Account)

	// 3) Build a list of principals to allow access
	principals := []string{
		fmt.Sprintf("arn:aws:iam::%s:root", accountID), // Your AWS account
	}
	principals = append(principals, sharedWith...) // Add additional ARNs

	// 4) Build the policy
	statements := []map[string]interface{}{
		{
			"Sid":    "AllowAccessToSpecificPrincipals",
			"Effect": "Allow",
			"Action": "secretsmanager:GetSecretValue",
			"Resource": fmt.Sprintf("arn:aws:secretsmanager:%s:%s:secret:%s",
				cfg.Region, accountID, secretID),
			"Principal": map[string]interface{}{
				"AWS": principals,
			},
		},
	}
	policyDoc := map[string]interface{}{
		"Version":   "2012-10-17",
		"Statement": statements,
	}

	policyBytes, err := json.Marshal(policyDoc)
	if err != nil {
		return fmt.Errorf("failed to marshal resource policy: %w", err)
	}

	// 5) Attach the resource policy to the secret
	smClient := secretsmanager.NewFromConfig(cfg)
	_, err = smClient.PutResourcePolicy(context.TODO(), &secretsmanager.PutResourcePolicyInput{
		SecretId:       aws.String(secretID),
		ResourcePolicy: aws.String(string(policyBytes)),
	})
	if err != nil {
		return fmt.Errorf("failed to attach resource policy: %w", err)
	}

	fmt.Printf("Updated sharing policy for secret: %s\n", secretID)
	return nil
}

// Utility Functions
func loadAWSConfig(profile string) (aws.Config, error) {
	return config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(profile),
	)
}

func ensurePrefix(secretID, prefix string) string {
	if !strings.HasPrefix(secretID, prefix) {
		return prefix + secretID
	}
	return secretID
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

func handleAWSSSOError(err error) error {
	if strings.Contains(err.Error(), "SSO session has expired") || strings.Contains(err.Error(), "InvalidGrantException") {
		return fmt.Errorf(
			"AWS SSO session has expired or is invalid. Please re-authenticate by running:\n\n"+
				"  aws sso login --profile <your-sdlc-profile>\n\n"+
				"Then try again with --profile <your-sdlc-profile> flag.\n\nOriginal error: %w",
			err,
		)
	}
	return fmt.Errorf("AWS operation failed: %w", err)
}

func exitWithError(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", msg)
	}
	os.Exit(1)
}

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
