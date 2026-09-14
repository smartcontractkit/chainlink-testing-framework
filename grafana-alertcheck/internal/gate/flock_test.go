package gate

import (
	"errors"
	"fmt"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsLockContention(t *testing.T) {
	contended := []error{syscall.EWOULDBLOCK, syscall.EAGAIN}
	for _, e := range contended {
		require.Truef(t, isLockContention(e), "isLockContention(%v)", e)
		// lockExclusive wraps the raw error via fmt.Errorf("flock: %w", ...).
		require.Truef(t, isLockContention(fmt.Errorf("flock: %w", e)), "isLockContention(wrapped %v)", e)
	}

	notContended := []error{
		syscall.EROFS,
		syscall.ENOTSUP,
		syscall.ENOLCK,
		syscall.EBADF,
		syscall.EIO,
		errors.New("something else"),
	}
	for _, e := range notContended {
		require.Falsef(t, isLockContention(e), "isLockContention(%v)", e)
		require.Falsef(t, isLockContention(fmt.Errorf("flock: %w", e)), "isLockContention(wrapped %v)", e)
	}
}
