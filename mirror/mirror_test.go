package mirror

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetImage(t *testing.T) {
	i, err := GetImage("testcontainers/ryuk")
	require.NoError(t, err)
	pattern := `^.+:\d+\.\d+\.\d+$`
	regex, err := regexp.Compile(pattern)
	require.NoError(t, err, "pattern did not compile")
	require.True(t, regex.MatchString(i), "should return an image with a version tag")
}
