package internal

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/spf13/cobra"
)

var ECR = &cobra.Command{
	Use:   "ecr",
	Short: "ECR commands",
}

var dalUntaggedImages = &cobra.Command{
	Use:   "delete-untagged-images",
	RunE:  deleteUntaggedImagesRunE,
	Short: "",
}

func init() {
	ECR.AddCommand(dalUntaggedImages)

	dalUntaggedImages.Flags().String("registry-id", "", "Image registry ID")
	if err := dalUntaggedImages.MarkFlagRequired("registry-id"); err != nil {
		log.Fatalf("Failed to mark 'registry-id' flag as required: %v", err)
	}
	dalUntaggedImages.Flags().String("repository-name", "", "Image repository name")
	if err := dalUntaggedImages.MarkFlagRequired("repository-name"); err != nil {
		log.Fatalf("Failed to mark 'repository-name' flag as required: %v", err)
	}
}

func deleteUntaggedImagesRunE(cmd *cobra.Command, args []string) error {
	registryId, err := cmd.Flags().GetString("registry-id")
	if err != nil {
		return fmt.Errorf("error getting registry ID: %v", err)
	}
	repositoryName, err := cmd.Flags().GetString("repository-name")
	if err != nil {
		return fmt.Errorf("error getting repository name: %v", err)
	}

	err = checkAWSCredentials()
	if err != nil {
		return err
	}

	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"),
	)
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	// Create an ECR client
	client := ecr.NewFromConfig(cfg)

	// Call the DescribeImages API
	images, err := client.DescribeImages(context.TODO(), &ecr.DescribeImagesInput{
		RepositoryName: aws.String(repositoryName),
		RegistryId:     &registryId,
		Filter: &types.DescribeImagesFilter{
			TagStatus: types.TagStatusUntagged,
		},
	})
	if err != nil {
		log.Fatalf("Unable to describe images, %v", err)
	}

	// Prepare the list of image IDs for deletion.
	var imageIDs []types.ImageIdentifier
	for _, image := range images.ImageDetails {
		imageIDs = append(imageIDs, types.ImageIdentifier{
			ImageDigest: image.ImageDigest,
		})
	}

	// Delete untagged images.
	if len(imageIDs) > 0 {
		_, err = client.BatchDeleteImage(context.TODO(), &ecr.BatchDeleteImageInput{
			RepositoryName: aws.String(repositoryName),
			ImageIds:       imageIDs,
		})
		if err != nil {
			log.Fatalf("Unable to delete images, %v", err)
		}
		fmt.Println("Untagged images deleted successfully.")
	} else {
		fmt.Println("No untagged images found to delete.")
	}

	return nil
}

func checkAWSCredentials() error {
	// List of required AWS environment variables.
	envVars := []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_SESSION_TOKEN"}

	for _, envVar := range envVars {
		if _, exists := os.LookupEnv(envVar); !exists {
			// Return an updated error message pointing to AWS Command line or programmatic access for credentials.
			return fmt.Errorf("environment variable %s is not set. Ensure your AWS credentials are configured correctly. You can find your credentials via AWS -> Command line or programmatic access", envVar)
		}
	}

	// All environment variables are set.
	return nil
}
