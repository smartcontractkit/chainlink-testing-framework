//go:build unix

package gate

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

// readyFD is where ExtraFiles[0] lands in the child: exec.Cmd starts extra
// descriptors at 3, after stdin, stdout and stderr.
const readyFD = 3

// detachedChild is a started recorder: the process, the read end of its
// readiness pipe, and the size the daemon log had before it wrote anything.
type detachedChild struct {
	cmd   *exec.Cmd
	ready *os.File
	// logOffset is where this run's output starts in the daemon log, so a
	// failure quotes this child and never a previous run's.
	logOffset int64
}

// spawnChild re-execs this binary as the detached recorder (§4.4). A trailing
// `&` is NOT sufficient: the child would keep the parent's session and process
// group, so it would still take the terminal's signals and, on a runner, die
// with the step that started it. Setsid gives it a new session AND a new
// process group, which is what makes it survive to the end of the window.
//
// There is deliberately no Windows counterpart — runners are Linux and
// goreleaser builds linux+darwin only (P12) — so the package does not build
// there at all rather than silently recording in the foreground.
func spawnChild(cfg WatchConfig) (detachedChild, error) {
	exe, err := os.Executable()
	if err != nil {
		return detachedChild{}, fmt.Errorf("find own executable: %w", err)
	}

	// The child is detached, so its output has nowhere to go but a file, and
	// that file is the only place it can ever explain a failure. Opened
	// O_APPEND, never O_TRUNC: --daemon-log is an operator-supplied path and
	// this process does not get to destroy what is already in it. The offset
	// below is what keeps a shared path from misattributing a failure.
	logFile, err := os.OpenFile(cfg.DaemonLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return detachedChild{}, fmt.Errorf("open daemon log %s: %w", cfg.DaemonLog, err)
	}
	// The child inherits the descriptor at Start; this process does not need
	// its own copy afterwards.
	defer logFile.Close()

	var logOffset int64
	if info, err := logFile.Stat(); err == nil {
		logOffset = info.Size()
	}

	// The readiness pipe (§ReadyFDFlag): the child gets the write end as
	// descriptor 3 and reports on it once it holds the log and is polling.
	readyRead, readyWrite, err := os.Pipe()
	if err != nil {
		return detachedChild{}, fmt.Errorf("open readiness pipe: %w", err)
	}

	cmd := exec.Command(exe, childArgs(cfg)...)
	cmd.Stdin = nil // /dev/null
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.ExtraFiles = []*os.File{readyWrite} // descriptor 3 in the child
	// The environment is how the connection details reach the child (§20.2):
	// the token must never appear in argv, where it would land in the process
	// table and in CI logs.
	cmd.Env = os.Environ()
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	if err := cmd.Start(); err != nil {
		readyRead.Close()
		readyWrite.Close()
		return detachedChild{}, fmt.Errorf("start recorder %s: %w", exe, err)
	}
	// Drop the parent's copy of the write end at once: with only the child
	// holding it, a child that dies before signalling closes the pipe and the
	// parent reads EOF instead of waiting out the whole timeout.
	readyWrite.Close()

	return detachedChild{cmd: cmd, ready: readyRead, logOffset: logOffset}, nil
}

// childArgs builds the child's command line. The rule set, every cadence and
// the recording's identity all come from the header the parent already wrote,
// and the connection details come from the environment — so what is left here
// is the log path plus the two run facts the header does not carry: the
// optional hard stop and the concurrency limit.
//
// Notably absent: --pidfile (the parent writes it, so check can find the pid
// the instant Watch returns), --alerts, --folder, --poll-interval, and
// anything derived from them. The CLI's `watch` FlagSet needs one flag of its
// own for this path — ReadyFDFlag — and dispatches to RunDaemonChild when it
// sees DaemonChildFlag.
func childArgs(cfg WatchConfig) []string {
	args := []string{"watch", DaemonChildFlag, "--out", cfg.Out, ReadyFDFlag, strconv.Itoa(readyFD)}
	if !cfg.Until.IsZero() {
		args = append(args, "--until", cfg.Until.Format(time.RFC3339))
	}
	if cfg.Concurrency > 0 {
		args = append(args, "--concurrency", strconv.Itoa(cfg.Concurrency))
	}
	return args
}
