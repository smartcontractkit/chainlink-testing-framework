package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContextWithTest(t *testing.T) {
	ctx := TestContext(t)
	require.NotEqual(t, context.Background(), ctx)
}

func TestContextWithNilTest(t *testing.T) {
	ctx := TestContext(nil)
	require.Equal(t, context.Background(), ctx)
}
