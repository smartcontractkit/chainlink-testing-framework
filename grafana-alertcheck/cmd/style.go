package main

import (
	"bytes"
	"io"
	"os"
	"strings"
)

// ANSI SGR codes for the human-facing notes and table footer. The colours are
// applied only when the destination is a terminal (see colorEnabled); a pipe,
// file or CI log gets plain text, so stdout stays reserved for --output json
// and no machine reader ever sees escape sequences.
const (
	ansiReset  = "\x1b[0m"
	ansiRed    = "\x1b[31m"
	ansiGreen  = "\x1b[32m"
	ansiYellow = "\x1b[33m"
	ansiCyan   = "\x1b[36m"
	// Orange has no entry in the base-16 palette; 256-colour 208 is a legible
	// orange used for warnings, distinct from the yellow used for notes.
	ansiOrange = "\x1b[38;5;208m"
)

// colorEnabled reports whether ANSI colour should be written to w. Colour is
// written only when three things hold: NO_COLOR is unset, w is a real *os.File
// (so text/tabwriter buffers, strings.Builder and bytes.Buffer tests all stay
// plain), and that file is a character device (a terminal, not a redirect).
func colorEnabled(w io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// styleLine applies the note vocabulary's colour to one line when enabled. The
// colour wraps the text only; the terminating newline is written uncoloured so
// the terminal's line discipline is never inside the escape sequence.
func styleLine(line string, enabled bool) string {
	if !enabled {
		return line
	}
	content := strings.TrimRight(line, "\n")
	var color string
	switch {
	case strings.HasPrefix(content, "warning:"):
		color = ansiOrange
	case strings.HasPrefix(content, "note:"):
		color = ansiYellow
	case strings.HasPrefix(content, "drain wait:"):
		color = ansiCyan
	}
	if color == "" {
		return line
	}
	return color + content + ansiReset + "\n"
}

// noteStyler wraps the gate package's Notes stream — a presentation seam that
// keeps colour out of the library. It colourises each line by its known prefix
// and separates the collection countdown from the setup phase with a single
// blank line before the first "collecting:" line. The gate keeps emitting plain
// prose; only the CLI lays it out.
type noteStyler struct {
	w             io.Writer
	enabled       bool
	pending       []byte
	sawCollecting bool
}

func newNoteStyler(w io.Writer) *noteStyler {
	return &noteStyler{w: w, enabled: colorEnabled(w)}
}

// startsSection reports whether a line opens a new phase of the stream and so
// deserves a blank line above it. "collecting:" opens the countdown (once —
// later countdown lines follow on from the first), and "drain wait:" opens the
// drain phase. The setup lines (planned run time, warning, min-observed, notes)
// are one contiguous block and are not separated from each other.
func (s *noteStyler) startsSection(line string) bool {
	switch {
	case strings.HasPrefix(line, "warning:"):
		return true
	case strings.HasPrefix(line, "drain wait:"):
		return true
	case strings.HasPrefix(line, "collecting:"):
		if s.sawCollecting {
			return false
		}
		s.sawCollecting = true
		return true
	}
	return false
}

func (s *noteStyler) Write(p []byte) (int, error) {
	n := len(p)
	s.pending = append(s.pending, p...)
	for {
		i := bytes.IndexByte(s.pending, '\n')
		if i < 0 {
			break
		}
		line := string(s.pending[:i+1])
		s.pending = s.pending[i+1:]

		if s.startsSection(line) {
			if _, err := io.WriteString(s.w, "\n"); err != nil {
				return n, err
			}
		}
		if _, err := io.WriteString(s.w, styleLine(line, s.enabled)); err != nil {
			return n, err
		}
	}
	return n, nil
}
