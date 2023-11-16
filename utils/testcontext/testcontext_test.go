package testcontext

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContextWithTest(t *testing.T) {
	ctx := Get(t)
	require.NotEqual(t, context.Background(), ctx)
}

func TestContextWithNilTest(t *testing.T) {
	ctx := Get(nil)
	require.Equal(t, context.Background(), ctx)
}
