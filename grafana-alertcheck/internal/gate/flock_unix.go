//go:build unix

package gate

import (
	"fmt"
	"os"
	"syscall"
)

// lockExclusive takes a non-blocking exclusive lock on f. Non-blocking is the
// point (§8): a second writer must fail immediately with an error the operator
// sees, not queue behind the first and start appending to a log somebody else
// already finished.
//
// There is deliberately no Windows implementation — runners are Linux and
// goreleaser builds linux+darwin only (P6, P12) — so the package does not
// build there at all rather than silently skipping the lock.
func lockExclusive(f *os.File) error {
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return fmt.Errorf("flock: %w", err)
	}
	return nil
}
