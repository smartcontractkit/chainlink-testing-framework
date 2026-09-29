package gate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// recorderStopTimeout bounds the wait for the recorder's exit after SIGTERM.
// recorderStopPoll is how often the wait re-checks the lock.
const (
	recorderStopTimeout = 30 * time.Second
	recorderStopPoll    = 100 * time.Millisecond
)

// StopConfig locates the recorder and selects the stop semantics.
type StopConfig struct {
	// Log is the JSONL path whose flock proves a writer exists right now.
	// PidFile defaults to <Log>.pid.
	Log     string
	PidFile string

	Clock Clock
	Notes io.Writer

	// Timeout bounds each wait (SIGTERM, then SIGKILL when Cleanup is set).
	// Zero means recorderStopTimeout; a test overrides it to drive the SIGKILL
	// path quickly.
	Timeout time.Duration

	// Cleanup is the `stop` command's semantics rather than check's: SIGKILL a
	// recorder that ignores SIGTERM, and remove the pidfile once no writer
	// holds the log. check leaves it false — a writer that will not exit means
	// the log cannot be trusted, which is a could-not-check, never a silent
	// kill.
	Cleanup bool
}

func (c StopConfig) withDefaults() StopConfig {
	if c.Clock == nil {
		c.Clock = SystemClock{}
	}
	if c.Notes == nil {
		c.Notes = io.Discard
	}
	if c.PidFile == "" && c.Log != "" {
		c.PidFile = c.Log + ".pid"
	}
	if c.Timeout <= 0 {
		c.Timeout = recorderStopTimeout
	}
	return c
}

// StopRecorder signals the recorder holding cfg.Log and waits for it to go.
// It returns the log under a shared flock — the caller must keep it open across
// any read and close it when done — even when no writer was found.
//
// The recorder holds the log exclusively; readers and cleanup hold it shared.
// Shared is refused only by an exclusive holder, so contention always means the
// recorder: a pid can be reused, the lock cannot.
func StopRecorder(ctx context.Context, cfg StopConfig) (*os.File, error) {
	cfg = cfg.withDefaults()
	if cfg.Log == "" {
		return nil, fmt.Errorf("stop recorder: no log path")
	}
	if cfg.PidFile == "" {
		return nil, fmt.Errorf("stop recorder: no pidfile path")
	}

	pid, pidErr := ReadPidFile(cfg.PidFile)
	if pidErr == nil {
		return cfg.stopNamed(ctx, pid)
	}
	if !cfg.Cleanup {
		return nil, fmt.Errorf("cannot stop the recorder: %w; a pidfile is written only once a recorder reports that it is running, so an unreadable one means the recording never started", pidErr)
	}
	return cfg.stopUnnamed(pidErr)
}

// stopNamed stops the recorder the pidfile names.
func (c StopConfig) stopNamed(ctx context.Context, pid int) (*os.File, error) {
	log, err := os.Open(c.Log)
	if err != nil {
		return nil, fmt.Errorf("cannot check whether a recorder is still running: %w; pidfile %s names a recorder that may still hold the log", err, c.PidFile)
	}

	// held is true when we TAKE the lock, i.e. when no writer holds it.
	held, err := tryLockShared(log)
	if err != nil {
		log.Close()
		return nil, err
	}
	if held {
		// No writer, so no signal — the pid may belong to somebody else by now.
		fmt.Fprintf(c.Notes, "note: no writer holds %s; the recorder has already finished\n", c.Log)
		return c.finish(log)
	}
	return c.signalAndWait(ctx, pid, log)
}

// stopUnnamed is cleanup with an unreadable pidfile: the log's lock still
// decides whether a writer exists, but a live one cannot be named, so it can
// only be reported, never signalled.
func (c StopConfig) stopUnnamed(pidErr error) (*os.File, error) {
	log, err := os.Open(c.Log)
	if err != nil {
		// Nothing was recorded at all: no readable pidfile and no log. This is
		// the only open failure cleanup may read as "nothing to stop".
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(c.Notes, "note: no recorder to stop: %s is absent and its pidfile %s is unreadable\n", c.Log, c.PidFile)
			return nil, nil
		}
		return nil, fmt.Errorf("cannot check whether a recorder is still running: %w; its unreadable pidfile %s: %w", err, c.PidFile, pidErr)
	}

	held, err := tryLockShared(log)
	if err != nil {
		log.Close()
		return nil, err
	}
	if held {
		fmt.Fprintf(c.Notes, "note: no writer holds %s; nothing to stop\n", c.Log)
		return c.finish(log)
	}
	log.Close()
	return nil, fmt.Errorf("a writer holds %s but pidfile %s is unreadable: %w", c.Log, c.PidFile, pidErr)
}

// finish removes the pidfile while the lock is still held — so a concurrent
// watch cannot take the lock and write a pidfile this then deletes — and
// returns the locked log.
func (c StopConfig) finish(log *os.File) (*os.File, error) {
	if err := c.removePidFile(); err != nil {
		log.Close()
		return nil, err
	}
	return log, nil
}

// errStillHeld is waitForRelease's timeout result.
var errStillHeld = errors.New("recorder still holds the log")

// waitForRelease polls the lock until it is released, the timeout passes, or
// ctx is cancelled.
func (c StopConfig) waitForRelease(ctx context.Context, log *os.File, timeout time.Duration) error {
	deadline := c.Clock.Now().Add(timeout)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.Clock.After(recorderStopPoll):
		}
		held, err := tryLockShared(log)
		if err != nil {
			return err
		}
		if held {
			return nil
		}
		if !c.Clock.Now().Before(deadline) {
			return errStillHeld
		}
	}
}

// signalAndWait sends SIGTERM to pid and waits for the log's lock to release.
// The lock, never the pid, is the proof the writer is gone.
func (c StopConfig) signalAndWait(ctx context.Context, pid int, log *os.File) (*os.File, error) {
	fail := func(err error) (*os.File, error) {
		log.Close()
		return nil, err
	}

	gone, err := signalRecorder(pid)
	if err != nil {
		return fail(err)
	}
	if gone {
		// The pid died while the lock was held. A free lock now means the
		// recorder finished in the race window; remaining contention means the
		// pidfile does not name the process holding the log.
		held, err := tryLockShared(log)
		if err != nil {
			return fail(err)
		}
		if held {
			return c.finish(log)
		}
		return fail(fmt.Errorf("a writer holds %s but pidfile %s names pid %d, which does not exist: the pidfile does not name the process that holds the log",
			c.Log, c.PidFile, pid))
	}

	err = c.waitForRelease(ctx, log, c.Timeout)
	if err == nil {
		return c.finish(log)
	}
	if !errors.Is(err, errStillHeld) {
		return fail(err)
	}
	if !c.Cleanup {
		return fail(fmt.Errorf("recorder pid %d still holds %s %s after SIGTERM; refusing to read a log a writer can still append to",
			pid, c.Log, c.Timeout))
	}

	// The recorder ignored SIGTERM; SIGKILL cannot be caught, so the lock drops
	// as soon as the process is gone.
	if err := killRecorder(pid); err != nil {
		return fail(err)
	}
	if err := c.waitForRelease(ctx, log, c.Timeout); err != nil {
		if errors.Is(err, errStillHeld) {
			err = fmt.Errorf("recorder pid %d still holds %s %s after SIGKILL", pid, c.Log, c.Timeout)
		}
		return fail(err)
	}
	return c.finish(log)
}

// removePidFile is a no-op outside cleanup semantics: check must not remove the
// pidfile, because its presence is what a later run reads to learn a recording
// started. An already-absent pidfile is success; any other failure is reported,
// so cleanup cannot claim to have cleaned up when it did not.
func (c StopConfig) removePidFile() error {
	if !c.Cleanup {
		return nil
	}
	if err := os.Remove(c.PidFile); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove pidfile %s: %w", c.PidFile, err)
	}
	return nil
}
