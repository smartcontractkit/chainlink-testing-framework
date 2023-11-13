package utils

import (
	"context"
	"testing"
)

// TestContext returns a context with the test's deadline, if available.
func TestContext(tb testing.TB) context.Context {
	ctx := context.Background()
	if tb == nil {
		return ctx
	}
	var cancel func()
	switch t := tb.(type) {
	case *testing.T:
		if d, ok := t.Deadline(); ok {
			ctx, cancel = context.WithDeadline(ctx, d)
		}
	}
	if cancel == nil {
		ctx, cancel = context.WithCancel(ctx)
	}
	tb.Cleanup(cancel)
	return ctx
}
