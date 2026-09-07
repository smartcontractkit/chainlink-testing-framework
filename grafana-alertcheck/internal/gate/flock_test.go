package gate

import (
	"errors"
	"fmt"
	"syscall"
	"testing"
)

func TestIsLockContention(t *testing.T) {
	contended := []error{syscall.EWOULDBLOCK, syscall.EAGAIN}
	for _, e := range contended {
		if !isLockContention(e) {
			t.Errorf("isLockContention(%v) = false, want true", e)
		}
		// lockExclusive wraps the raw error via fmt.Errorf("flock: %w", ...).
		if !isLockContention(fmt.Errorf("flock: %w", e)) {
			t.Errorf("isLockContention(wrapped %v) = false, want true", e)
		}
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
		if isLockContention(e) {
			t.Errorf("isLockContention(%v) = true, want false (not a contender)", e)
		}
		if isLockContention(fmt.Errorf("flock: %w", e)) {
			t.Errorf("isLockContention(wrapped %v) = true, want false", e)
		}
	}
}
