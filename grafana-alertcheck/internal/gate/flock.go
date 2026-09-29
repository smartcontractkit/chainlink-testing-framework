package gate

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

// lockExclusive takes a non-blocking exclusive lock on f. Non-blocking is the
// point: a second writer must fail immediately with an error the operator
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

// tryLockExclusive is lockExclusive read as a question: held is false when
// another process holds the lock; err is non-nil only for a real failure.
func tryLockExclusive(f *os.File) (held bool, err error) {
	switch err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); {
	case err == nil:
		return true, nil
	case isLockContention(err):
		return false, nil
	default:
		return false, fmt.Errorf("flock %s: %w", f.Name(), err)
	}
}

// tryLockShared is the reader side of the same question: held is false when an
// EXCLUSIVE holder exists. The recorder holds the log exclusively; readers and
// cleanup hold it shared, so shared contention always means the recorder, never
// another reader.
func tryLockShared(f *os.File) (held bool, err error) {
	switch err := syscall.Flock(int(f.Fd()), syscall.LOCK_SH|syscall.LOCK_NB); {
	case err == nil:
		return true, nil
	case isLockContention(err):
		return false, nil
	default:
		return false, fmt.Errorf("flock %s: %w", f.Name(), err)
	}
}
