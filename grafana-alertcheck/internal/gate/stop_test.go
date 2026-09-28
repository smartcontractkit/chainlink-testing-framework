package gate

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// An empty log is enough: StopRecorder opens it only to probe the flock.
func emptyLog(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "log.jsonl")
	require.NoError(t, os.WriteFile(path, nil, 0o644)) // nolint:gosec // test-only temp file
	return path
}

// Cleanup SIGKILLs a recorder that ignores SIGTERM. The lock holder is a real
// process that ignores SIGTERM (watch_daemon_test.go), so only a real SIGKILL
// releases the flock.
func TestStopRecorderCleanupKillsARecorderThatIgnoresSigterm(t *testing.T) {
	logPath := emptyLog(t)
	pid := startLockHolder(t, logPath)
	writePid(t, logPath+".pid", fmt.Sprintf("%d\n", pid))

	// Real clock, short timeout: the SIGTERM wait must elapse before SIGKILL.
	held, err := StopRecorder(context.Background(), StopConfig{
		Log:     logPath,
		Clock:   SystemClock{},
		Timeout: 200 * time.Millisecond,
		Cleanup: true,
	})
	require.NoError(t, err)
	require.NotNil(t, held, "the log is returned held")
	// A free lock is proof the SIGKILL landed.
	free, err := tryLockExclusive(held)
	require.NoError(t, err)
	require.True(t, free, "the killed recorder still holds the lock")
	require.NoError(t, held.Close())

	require.NoFileExists(t, logPath+".pid", "cleanup removes the pidfile")
}

// With no writer, cleanup sends no signal but still removes the pidfile.
func TestStopRecorderCleanupRemovesPidfileWhenNoWriterHoldsTheLog(t *testing.T) {
	logPath := emptyLog(t)
	writePid(t, logPath+".pid", "12345\n") // lock is free, so no signal is sent

	held, err := StopRecorder(context.Background(), StopConfig{
		Log:     logPath,
		Clock:   newFakeClock(testNow),
		Cleanup: true,
	})
	require.NoError(t, err)
	require.NotNil(t, held)
	require.NoError(t, held.Close())
	require.NoFileExists(t, logPath+".pid")
}

// A second stop while another operation still holds the returned lock must see
// a reader, not a writer; otherwise it signals the stale pid.
func TestStopRecorderDoesNotSignalPastAnotherOperationsLock(t *testing.T) {
	logPath := emptyLog(t)
	first, err := StopRecorder(context.Background(), StopConfig{Log: logPath, Clock: newFakeClock(testNow), Cleanup: true})
	require.NoError(t, err)
	require.NotNil(t, first)
	defer first.Close()

	bystander := exec.Command("sleep", "30")
	require.NoError(t, bystander.Start())
	exited := make(chan struct{})
	go func() { _ = bystander.Wait(); close(exited) }()
	t.Cleanup(func() { _ = bystander.Process.Kill() })
	writePid(t, logPath+".pid", fmt.Sprintf("%d\n", bystander.Process.Pid))

	second, err := StopRecorder(context.Background(), StopConfig{Log: logPath, Clock: newFakeClock(testNow), Cleanup: true})
	require.NoError(t, err)
	require.NotNil(t, second)
	require.NoError(t, second.Close())

	select {
	case <-exited:
		require.Fail(t, "the second stop signalled the stale pid")
	case <-time.After(200 * time.Millisecond):
	}
}

// Cleanup is safe to run twice, or after check already stopped the recorder.
func TestStopRecorderCleanupIsIdempotent(t *testing.T) {
	logPath := emptyLog(t)
	held, err := StopRecorder(context.Background(), StopConfig{
		Log:     logPath,
		Clock:   newFakeClock(testNow),
		Cleanup: true,
	})
	require.NoError(t, err)
	require.NotNil(t, held, "the locked log is returned even when nothing was running")
	require.NoError(t, held.Close())
}

// Cleanup must not read an unreadable pidfile as "nothing to stop" while a
// writer still holds the log: the flock is authoritative.
func TestStopRecorderCleanupRefusesToIgnoreALiveWriterWithNoPidfile(t *testing.T) {
	logPath := emptyLog(t)
	_ = startLockHolder(t, logPath) // deliberately no pidfile

	_, err := StopRecorder(context.Background(), StopConfig{
		Log:     logPath,
		Clock:   newVirtualClock(testNow),
		Cleanup: true,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "unreadable")
}

// Cleanup treats only a missing log as idempotent; a real open error (here
// ENOTDIR) must not be reported as "nothing to stop".
func TestStopRecorderCleanupPropagatesNonMissingOpenErrors(t *testing.T) {
	dir := t.TempDir()
	notADir := filepath.Join(dir, "file")
	require.NoError(t, os.WriteFile(notADir, nil, 0o644)) // nolint:gosec // test-only temp file
	logPath := filepath.Join(notADir, "log.jsonl")

	_, err := StopRecorder(context.Background(), StopConfig{
		Log:     logPath,
		Clock:   newFakeClock(testNow),
		Cleanup: true,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "unreadable pidfile")
}

// Without cleanup an absent pidfile is a hard error, as check requires.
func TestStopRecorderWithoutCleanupRefusesAMissingPidfile(t *testing.T) {
	logPath := emptyLog(t)
	_, err := StopRecorder(context.Background(), StopConfig{
		Log:   logPath,
		Clock: newFakeClock(testNow),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot stop the recorder")
}

// Without cleanup a SIGTERM-ignoring holder is a hard error; the pidfile stays.
func TestStopRecorderWithoutCleanupLeavesThePidfileAndErrors(t *testing.T) {
	logPath := emptyLog(t)
	pid := startLockHolder(t, logPath)
	writePid(t, logPath+".pid", fmt.Sprintf("%d\n", pid))

	_, err := StopRecorder(context.Background(), StopConfig{
		Log:   logPath,
		Clock: newVirtualClock(testNow),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "still holds")
	require.FileExists(t, logPath+".pid")
}
