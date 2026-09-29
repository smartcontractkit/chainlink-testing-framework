package gate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return b
}

func appendFile(t *testing.T, path string, b []byte) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644) // nolint:gosec // test-only temp file
	require.NoError(t, err)
	_, err = f.Write(b)
	require.NoError(t, err)
	require.NoError(t, f.Close())
}

// A torn write waits in the buffer for its newline; it is never reported, and
// never repeated once completed.
func TestLogTailerReadsOnlyCompleteRecords(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")

	header := mustJSON(t, headerRecord{Type: RecordHeader, Header: Header{
		SchemaVersion: LogSchemaVersion, StartedAt: testNow,
	}})
	poll1 := mustJSON(t, pollRecord{Type: RecordPoll, Poll: Poll{RuleUID: checkUID, GrafanaNow: testNow, Found: true}})
	poll2 := mustJSON(t, pollRecord{Type: RecordPoll, Poll: Poll{RuleUID: checkUID, GrafanaNow: testNow.Add(time.Second), Found: true}})

	// Header and one poll complete; the second poll is torn mid-write.
	split := len(poll2) / 2
	appendFile(t, path, []byte(string(header)+"\n"))
	appendFile(t, path, []byte(string(poll1)+"\n"))
	appendFile(t, path, poll2[:split])

	tail, err := newLogTailer(path)
	require.NoError(t, err)
	defer tail.Close()

	polls, sentinel, err := tail.read()
	require.NoError(t, err)
	require.Len(t, polls, 1, "the torn record must not be reported")
	require.Nil(t, sentinel)

	// A second read with nothing new must not repeat the first poll.
	polls, _, err = tail.read()
	require.NoError(t, err)
	require.Empty(t, polls)

	// Finish the torn line; now it is reportable, exactly once.
	appendFile(t, path, append(poll2[split:], '\n'))
	polls, sentinel, err = tail.read()
	require.NoError(t, err)
	require.Len(t, polls, 1)
	require.True(t, polls[0].GrafanaNow.Equal(testNow.Add(time.Second)))
	require.Nil(t, sentinel)

	// The sentinel is reported when it lands.
	stopped := mustJSON(t, stoppedRecord{Type: RecordStopped, At: testNow.Add(2 * time.Second)})
	appendFile(t, path, append(stopped, '\n'))
	_, sentinel, err = tail.read()
	require.NoError(t, err)
	require.NotNil(t, sentinel)
	require.True(t, sentinel.Equal(testNow.Add(2*time.Second)))
}

// A complete but unparseable line is corruption: fail closed, never skip.
func TestLogTailerRejectsAnUnparseableCompleteLine(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.jsonl")
	appendFile(t, path, []byte("not json\n"))

	tail, err := newLogTailer(path)
	require.NoError(t, err)
	defer tail.Close()

	_, _, err = tail.read()
	require.Error(t, err)
	require.Contains(t, err.Error(), "unparseable")
}
