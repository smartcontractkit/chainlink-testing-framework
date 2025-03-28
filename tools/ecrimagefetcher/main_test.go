package main

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"
)

func mockFetchImageDetailsSuccess(_ string) ([]byte, error) {
	return []byte(`[
		{
			"imageTags": ["v1.0.0"]
		},
		{
			"imageTags": ["v1.1.2"]
		},
		{
			"imageTags": ["v1.1.11"]
		},
		{
			"imageTags": ["v1.2.0"]
		}
	]`), nil
}

func mockFetchImageMixedDetailsSuccess(_ string) ([]byte, error) {
	return []byte(`[
		{
			"imageTags": ["alltools-0.0.1"]
		},
		{
			"imageTags": ["1.0.0"]
		},
		{
			"imageTags": ["0.5.0-alpha1"]
		},
		{
			"imageTags": ["0.5.0-alpha2"]
		},
		{
			"imageTags": ["1.0.0-alpha"]
		},
		{
			"imageTags": ["1.0.0-alpha.1"]
		},
		{
			"imageTags": ["1.0.0-alpha.2"]
		},
		{
			"imageTags": ["1.0.0-beta"]
		},
		{
			"imageTags": ["1.0.0-rc2"]
		},
		{
			"imageTags": ["1.0.0-rc1"]
		},
		{
			"imageTags": ["v1.1.2-alpha"]
		}
	]`), nil
}

func mockFetchImageDetailsError(_ string) ([]byte, error) {
	return nil, fmt.Errorf("failed to describe images")
}

func TestGetLatestImages(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		images, err := getLatestImages(mockFetchImageDetailsSuccess, "test-repo", "v1.*", 2, nil)
		require.NoError(t, err)
		require.Equal(t, "test-repo:v1.2.0,test-repo:v1.1.11", images)
	})

	t.Run("SuccessWithMixedVersionsWithoutFilter", func(t *testing.T) {
		images, err := getLatestImages(mockFetchImageMixedDetailsSuccess, "test-repo", "", 11, nil)
		require.NoError(t, err)
		require.Equal(t, "test-repo:v1.1.2-alpha,test-repo:1.0.0,test-repo:1.0.0-rc2,test-repo:1.0.0-rc1,test-repo:1.0.0-beta,test-repo:1.0.0-alpha.2,test-repo:1.0.0-alpha.1,test-repo:1.0.0-alpha,test-repo:0.5.0-alpha2,test-repo:0.5.0-alpha1,test-repo:0.0.1", images)
	})

	t.Run("SuccessWithMixedVersionsWithFilter", func(t *testing.T) {
		images, err := getLatestImages(mockFetchImageMixedDetailsSuccess, "test-repo", "^1.*", 4, nil)
		require.NoError(t, err)
		require.Equal(t, "test-repo:1.0.0,test-repo:1.0.0-rc2,test-repo:1.0.0-rc1,test-repo:1.0.0-beta", images)
	})

	t.Run("ErrorFetchingDetails", func(t *testing.T) {
		_, err := getLatestImages(mockFetchImageDetailsError, "test-repo", "v1.*", 2, nil)
		require.Error(t, err)
		require.Equal(t, "failed to describe images: failed to describe images", err.Error())
	})

	t.Run("ErrorParsingTags", func(t *testing.T) {
		_, err := getLatestImages(mockFetchImageDetailsSuccess, "test-repo", "invalid[regex", 2, nil)
		require.Error(t, err)
		require.Equal(t, "failed to parse image tags: failed to compile regex: error parsing regexp: missing closing ]: `[regex`", err.Error())
	})

	t.Run("InsufficientTags", func(t *testing.T) {
		images, err := getLatestImages(mockFetchImageDetailsSuccess, "test-repo", "v1.*", 5, nil)
		require.NoError(t, err)
		require.Equal(t, "test-repo:v1.2.0,test-repo:v1.1.11,test-repo:v1.1.2,test-repo:v1.0.0", images)
	})

	t.Run("WithConstraint", func(t *testing.T) {
		constraints, err := semver.NewConstraint("<v1.2.0")
		require.NoError(t, err)
		images, err := getLatestImages(mockFetchImageDetailsSuccess, "test-repo", "v1.*", 2, constraints)
		require.NoError(t, err)
		require.Equal(t, "test-repo:v1.1.11,test-repo:v1.1.2", images)
	})
}

func TestValidateInputs(t *testing.T) {
	t.Run("MissingArguments", func(t *testing.T) {
		os.Args = []string{"main"}
		expectedError := errors.New("usage: <repository_name> <grep_string> <count> [<ignored_tags>]")
		require.EqualError(t, validateInputs(), expectedError.Error())
	})

	t.Run("EmptyRepositoryName", func(t *testing.T) {
		os.Args = []string{"main", "", "v1.*", "2"}
		expectedError := errors.New("error: repository_name cannot be empty")
		require.EqualError(t, validateInputs(), expectedError.Error())
	})

	t.Run("EmptyGrepString", func(t *testing.T) {
		os.Args = []string{"main", "test-repo", "", "2"}
		expectedError := errors.New("error: grep_string cannot be empty")
		require.EqualError(t, validateInputs(), expectedError.Error())
	})

	t.Run("InvalidGrepString", func(t *testing.T) {
		os.Args = []string{"main", "test-repo", "invalid[regex", "2"}
		expectedError := errors.New("error: grep_string is not a valid regex")
		require.EqualError(t, validateInputs(), expectedError.Error())
	})

	t.Run("NonIntegerCount", func(t *testing.T) {
		os.Args = []string{"main", "test-repo", "v1.*", "two"}
		expectedError := fmt.Errorf("error: count must be an integer, but %s is not an integer", "two")
		require.EqualError(t, validateInputs(), expectedError.Error())
	})

	t.Run("EmptyConstraint", func(t *testing.T) {
		os.Args = []string{"main", "test-repo", "v1.*", "2", ">=v1.0.0,"}
		expectedError := errors.New("error: semver constraint cannot be empty")
		require.EqualError(t, validateInputs(), expectedError.Error())
	})

	t.Run("InvalidConstraint", func(t *testing.T) {
		os.Args = []string{"main", "test-repo", "v1.*", "2", "asdasd87"}
		expectedError := errors.New("error: invalid semver constraint: improper constraint: asdasd87")
		require.EqualError(t, validateInputs(), expectedError.Error())
	})

	t.Run("Success", func(t *testing.T) {
		os.Args = []string{"main", "test-repo", "v1.*", "2"}
		require.NoError(t, validateInputs())
	})
}
