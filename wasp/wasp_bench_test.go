package wasp

import (
	"github.com/stretchr/testify/require"
	"math"
	"os"
	"testing"
	"time"
)

var (
	NoLimitSchedule = Plain(math.MaxInt64, 20*time.Minute)
)

func BenchmarkPacedCall(b *testing.B) {
	_ = os.Setenv("WASP_LOG_LEVEL", "warn")
	gen, err := NewGenerator(&Config{
		LoadType:          RPS,
		StatsPollInterval: 1 * time.Second,
		Schedule:          NoLimitSchedule,
		Gun:               NewMockGun(&MockGunConfig{}),
	})
	require.NoError(b, err)
	gen.setupSchedule()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gen.pacedCall()
	}
}
