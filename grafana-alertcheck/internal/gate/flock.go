package gate

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

// lockExclusive takes a non-blocking exclusive lock on f. Non-blocking is the
// point (§8): a second writer must fail immediately with an error the operator
// sees, not queue behind the first and start appending to a log somebody else
// already finished.
func lockExclusive(f *os.File) error {
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return fmt.Errorf("flock: %w", err)
	}
	return nil
}

// isLockContention reports whether a flock failure means another writer holds
// the lock, as opposed to an unrelated failure.
func isLockContention(err error) bool {
	return errors.Is(err, syscall.EWOULDBLOCK) || errors.Is(err, syscall.EAGAIN)
}
