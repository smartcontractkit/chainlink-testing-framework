package gate

import (
	"errors"
	"fmt"
	"syscall"
)

// signalRecorder asks the recorder to stop (§4.4 step 2).
//
// The caller must have established that a writer is alive — by taking the
// log's flock and being refused — before it calls this. Nothing removes the
// pidfile when a recorder exits cleanly, so a pid read without that proof can
// name any same-user process that has since inherited it.
//
// gone reports ESRCH. Given the lock proof, that is a broken contract rather
// than a clean stop, and stopRecorder treats it as one; the value is reported
// instead of raised here because this function knows the errno and not what
// it means.
func signalRecorder(pid int) (gone bool, err error) {
	switch err := syscall.Kill(pid, syscall.SIGTERM); {
	case err == nil:
		return false, nil
	case errors.Is(err, syscall.ESRCH):
		return true, nil
	default:
		return false, fmt.Errorf("signal recorder pid %d: %w", pid, err)
	}
}
