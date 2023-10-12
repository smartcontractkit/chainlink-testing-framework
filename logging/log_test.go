package logging

import (
	"testing"
	"time"
)

func TestLogAfterTestEnd(t *testing.T) {
	l := GetTestLogger(t)
	go func() {
		for i := 0; i < 5000; i++ {
			l.Info().Msg("test")
		}
	}()
	time.Sleep(1 * time.Millisecond)
}
