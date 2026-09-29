package gate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// tailReadChunk bounds one read of the growing log.
const tailReadChunk = 64 * 1024

// logTailer incrementally reads a recorder's log while the recorder may still
// be appending to it, for the fail-fast guard only.
//
// It consumes complete, newline-terminated records only — a torn write stays
// buffered until its newline arrives — and never feeds the final
// classification. Once the run stops, check still reads the log once through
// the strict whole-file ReadLog.
type logTailer struct {
	f      *os.File
	offset int64  // file position up to which the file has been read
	buf    []byte // bytes read but not yet forming a complete line
}

func newLogTailer(path string) (*logTailer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("tail log %s: %w", path, err)
	}
	return &logTailer{f: f}, nil
}

func (t *logTailer) Close() error { return t.f.Close() }

// read returns polls appended since the previous call, plus the sentinel when
// the recorder has finished.
func (t *logTailer) read() ([]Poll, *time.Time, error) {
	if _, err := t.f.Seek(t.offset, io.SeekStart); err != nil {
		return nil, nil, fmt.Errorf("tail log: seek: %w", err)
	}
	chunk := make([]byte, tailReadChunk)
	for {
		n, err := t.f.Read(chunk)
		if n > 0 {
			t.buf = append(t.buf, chunk[:n]...)
			t.offset += int64(n)
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, nil, fmt.Errorf("tail log: read: %w", err)
		}
		if n == 0 {
			break
		}
	}

	var (
		polls    []Poll
		sentinel *time.Time
	)
	for {
		i := bytes.IndexByte(t.buf, '\n')
		if i < 0 {
			break
		}
		line := t.buf[:i]
		t.buf = t.buf[i+1:]
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var probe struct {
			Type RecordType `json:"type"`
		}
		if err := json.Unmarshal(line, &probe); err != nil {
			return nil, nil, fmt.Errorf("tail log: unparseable complete record: %w", err)
		}
		switch probe.Type {
		case RecordHeader:
			// Already read authoritatively by ReadLogHeader; nothing to do.
		case RecordPoll:
			var rec pollRecord
			if err := json.Unmarshal(line, &rec); err != nil {
				return nil, nil, fmt.Errorf("tail log: unparseable poll: %w", err)
			}
			polls = append(polls, rec.Poll)
		case RecordStopped:
			var rec stoppedRecord
			if err := json.Unmarshal(line, &rec); err != nil {
				return nil, nil, fmt.Errorf("tail log: unparseable sentinel: %w", err)
			}
			at := rec.At
			sentinel = &at
		default:
			return nil, nil, fmt.Errorf("tail log: unknown record type %q", probe.Type)
		}
	}
	return polls, sentinel, nil
}
