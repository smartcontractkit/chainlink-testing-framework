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

// tryLockExclusive is the same call read as a question rather than as a
// demand: held is false when another process holds the lock, and err is
// non-nil only for a failure that is not contention.
//
// check needs that distinction where NewWriter does not. NewWriter is entitled
// to treat any refusal as "another writer has it", because it wants the lock;
// check only wants to know whether a writer EXISTS (§4.4). The lock answers
// that directly, where a pid can only infer it — the kernel releases a flock
// when the holder exits, crash included, and pids get reused.
func tryLockExclusive(f *os.File) (held bool, err error) {
	switch err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); {
	case err == nil:
		return true, nil
	case errors.Is(err, syscall.EWOULDBLOCK):
		return false, nil
	default:
		return false, fmt.Errorf("flock %s: %w", f.Name(), err)
	}
}
