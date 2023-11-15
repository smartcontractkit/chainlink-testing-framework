package logging

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

func TestGetTestLogger(t *testing.T) {
	l := GetTestLogger(t)
	l.Info().Msg("test")
	require.NotNil(t, l)
}

func TestGetTestContainersGoTestLogger(t *testing.T) {
	l := GetTestContainersGoTestLogger(t)
	require.NotNil(t, l.(CustomT).L)
}
