package havoc

import (
	"errors"
	"time"
)

// WaitForAllChaosRunning blocks until chaos experiments are running
func WaitForAllChaosRunning(chaosObjects []*Chaos, timeoutDuration time.Duration) error {
	timeout := time.NewTimer(timeoutDuration)
	defer timeout.Stop()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	runningStatus := make(map[*Chaos]bool)
	for _, chaos := range chaosObjects {
		runningStatus[chaos] = false
	}

	for {
		allRunning := true

		select {
		case <-timeout.C:
			return errors.New("timeout reached before all chaos experiments became running")
		case <-ticker.C:
			for chaos, isRunning := range runningStatus {
				if !isRunning { // Only check if not already marked as running
					if chaos.Status == StatusRunning {
						runningStatus[chaos] = true
					} else {
						allRunning = false
					}
				}
			}

			if allRunning {
				return nil // All chaos objects are running, can exit
			}
			// Otherwise, continue the loop
		}
	}
}
