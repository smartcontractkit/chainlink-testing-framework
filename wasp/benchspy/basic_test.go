package benchspy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

func TestBenchSpy_NewBasicData(t *testing.T) {
	testName := t.Name()
	gen := &wasp.Generator{
		Cfg: &wasp.Config{
			T:        t,
			GenName:  "gen1",
			Schedule: []*wasp.Segment{{From: 1, Duration: time.Hour, StartTime: time.Now().Add(-time.Hour), EndTime: time.Now()}},
		},
	}

	tests := []struct {
		name        string
		commitOrTag string
		generators  []*wasp.Generator
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid basic data",
			commitOrTag: "abc123",
			generators:  []*wasp.Generator{gen},
			wantErr:     false,
		},
		{
			name:        "no generators",
			commitOrTag: "abc123",
			generators:  []*wasp.Generator{},
			wantErr:     true,
			errMsg:      "at least one generator is required",
		},
		{
			name:        "generator without testing.T",
			commitOrTag: "abc123",
			generators:  []*wasp.Generator{{Cfg: &wasp.Config{GenName: "gen1"}}},
			wantErr:     true,
			errMsg:      "generators are not associated with a testing.T instance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bd, err := NewBasicData(tt.commitOrTag, tt.generators...)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.commitOrTag, bd.CommitOrTag)
			assert.Equal(t, testName, bd.TestName)
			assert.Len(t, bd.GeneratorConfigs, len(tt.generators))
		})
	}
}

func TestBenchSpy_MustNewBasicData(t *testing.T) {
	gen := &wasp.Generator{
		Cfg: &wasp.Config{
			T:       t,
			GenName: "gen1",
			Schedule: []*wasp.Segment{
				{From: 1, Duration: time.Hour, StartTime: time.Now().Add(-time.Hour), EndTime: time.Now()},
			},
		},
	}

	t.Run("successful creation", func(t *testing.T) {
		assert.NotPanics(t, func() {
			bd := MustNewBasicData("abc123", gen)
			assert.Equal(t, "abc123", bd.CommitOrTag)
		})
	})

	t.Run("panics on error", func(t *testing.T) {
		assert.Panics(t, func() {
			MustNewBasicData("abc123") // no generators
		})
	})
}

func TestBenchSpy_BasicData_FillStartEndTimes(t *testing.T) {
	now := time.Now()

	t.Run("successful fill", func(t *testing.T) {
		gen := &wasp.Generator{
			Cfg: &wasp.Config{
				T:       t,
				GenName: "gen1",
				Schedule: []*wasp.Segment{
					{StartTime: now, EndTime: now.Add(time.Hour)},
					{StartTime: now.Add(2 * time.Hour), EndTime: now.Add(3 * time.Hour)},
				},
			},
		}

		bd, err := NewBasicData("abc123", gen)
		require.NoError(t, err)

		err = bd.FillStartEndTimes()
		require.NoError(t, err)
		assert.Equal(t, now, bd.TestStart)
		assert.Equal(t, now.Add(3*time.Hour), bd.TestEnd)
	})

	t.Run("error on missing start time", func(t *testing.T) {
		gen := &wasp.Generator{
			Cfg: &wasp.Config{
				T:       t,
				GenName: "gen1",
				Schedule: []*wasp.Segment{
					{EndTime: now.Add(-time.Hour), Type: wasp.SegmentType_Plain, From: 1, Duration: time.Hour},
				},
			},
		}

		bd, err := NewBasicData("abc123", gen)
		require.Nil(t, bd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "start time is missing")
	})

	t.Run("error on missing end time", func(t *testing.T) {
		gen := &wasp.Generator{
			Cfg: &wasp.Config{
				T:       t,
				GenName: "gen1",
				Schedule: []*wasp.Segment{
					{StartTime: now.Add(time.Hour)},
				},
			},
		}

		bd, err := NewBasicData("abc123", gen)
		require.Nil(t, bd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "end time is missing")
	})

	t.Run("error on empty schedule", func(t *testing.T) {
		gen := &wasp.Generator{
			Cfg: &wasp.Config{
				T:       t,
				GenName: "gen1",
			},
		}

		bd, err := NewBasicData("abc123", gen)
		require.Nil(t, bd)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "schedule is empty for generator gen1")
	})
}

func TestBenchSpy_BasicData_FillStartEndTimes_MultipleGenerators(t *testing.T) {
	baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		gens      []*wasp.Generator
		wantStart time.Time
		wantEnd   time.Time
		wantErr   bool
		errMsg    string
	}{
		{
			name: "multiple generators with different schedules",
			gens: []*wasp.Generator{
				{
					Cfg: &wasp.Config{
						T:       t,
						GenName: "gen1",
						Schedule: []*wasp.Segment{
							{StartTime: baseTime, EndTime: baseTime.Add(2 * time.Hour)},
						},
					},
				},
				{
					Cfg: &wasp.Config{
						T:       t,
						GenName: "gen2",
						Schedule: []*wasp.Segment{
							{StartTime: baseTime.Add(time.Hour), EndTime: baseTime.Add(3 * time.Hour)},
						},
					},
				},
			},
			wantStart: baseTime,
			wantEnd:   baseTime.Add(3 * time.Hour),
			wantErr:   false,
		},
		{
			name: "multiple segments in one generator",
			gens: []*wasp.Generator{
				{
					Cfg: &wasp.Config{
						T:       t,
						GenName: "gen1",
						Schedule: []*wasp.Segment{
							{StartTime: baseTime.Add(2 * time.Hour), EndTime: baseTime.Add(3 * time.Hour)},
							{StartTime: baseTime, EndTime: baseTime.Add(time.Hour)},
						},
					},
				},
			},
			wantStart: baseTime,
			wantEnd:   baseTime.Add(3 * time.Hour),
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bd, err := NewBasicData("test-commit", tt.gens...)
			require.NoError(t, err)

			err = bd.FillStartEndTimes()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestBenchSpy_TestBasicData_Validate(t *testing.T) {
	gen := &wasp.Generator{
		Cfg: &wasp.Config{
			T:       t,
			GenName: "gen1",
		},
	}

	tests := []struct {
		name    string
		bd      *BasicData
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid data",
			bd: &BasicData{
				TestStart:        time.Now(),
				TestEnd:          time.Now().Add(time.Hour),
				GeneratorConfigs: map[string]*wasp.Config{"gen1": gen.Cfg},
			},
			wantErr: false,
		},
		{
			name: "missing start time",
			bd: &BasicData{
				TestEnd:          time.Now().Add(time.Hour),
				GeneratorConfigs: map[string]*wasp.Config{"gen1": gen.Cfg},
			},
			wantErr: true,
			errMsg:  "test start time is missing",
		},
		{
			name: "missing generator configs",
			bd: &BasicData{
				TestStart: time.Now(),
				TestEnd:   time.Now().Add(time.Hour),
			},
			wantErr: true,
			errMsg:  "generator configs are missing",
		},
		{
			name: "test start and end time are the same",
			bd: &BasicData{
				TestStart: time.Now(),
				TestEnd:   time.Now(),
			},
			wantErr: true,
			errMsg:  "test duration is less than a second",
		},
		{
			name: "test end time before start time",
			bd: &BasicData{
				TestStart: time.Now().Add(time.Hour),
				TestEnd:   time.Now(),
			},
			wantErr: true,
			errMsg:  "test end time is before test start time",
		},
		{
			name: "test end time are start time < 1s apart",
			bd: &BasicData{
				TestStart: time.Now(),
				TestEnd:   time.Now().Add(time.Second - time.Millisecond),
			},
			wantErr: true,
			errMsg:  "test duration is less than a second",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.bd.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestBenchSpy_BasicData_IsComparable(t *testing.T) {
	baseConfig := &wasp.Config{
		T:        t,
		GenName:  "gen1",
		LoadType: wasp.VU,
		Schedule: []*wasp.Segment{
			{From: 1, Duration: time.Hour},
		},
		CallTimeout: time.Second,
	}

	tests := []struct {
		name    string
		bd1     BasicData
		bd2     BasicData
		wantErr bool
		errMsg  string
	}{
		{
			name: "identical configs",
			bd1: BasicData{GeneratorConfigs: map[string]*wasp.Config{
				"gen1": baseConfig,
			}},
			bd2: BasicData{GeneratorConfigs: map[string]*wasp.Config{
				"gen1": baseConfig,
			}},
			wantErr: false,
		},
		{
			name: "different generator count",
			bd1: BasicData{GeneratorConfigs: map[string]*wasp.Config{
				"gen1": baseConfig,
			}},
			bd2: BasicData{GeneratorConfigs: map[string]*wasp.Config{
				"gen1": baseConfig,
				"gen2": baseConfig,
			}},
			wantErr: true,
			errMsg:  "generator configs count is different",
		},
		{
			name: "other report is missing a generator #1",
			bd1: BasicData{GeneratorConfigs: map[string]*wasp.Config{
				"gen1": baseConfig,
			}},
			bd2: BasicData{GeneratorConfigs: map[string]*wasp.Config{
				"gen2": baseConfig,
			}},
			wantErr: true,
			errMsg:  "generator config gen1 is missing from the other report",
		},
		{
			name: "other report is missing a generator #2",
			bd1: BasicData{GeneratorConfigs: map[string]*wasp.Config{
				"gen1": baseConfig,
				"gen2": baseConfig,
			}},
			bd2: BasicData{GeneratorConfigs: map[string]*wasp.Config{
				"gen2": baseConfig,
				"gen3": baseConfig,
			}},
			wantErr: true,
			errMsg:  "generator config gen1 is missing from the other report",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.bd1.IsComparable(tt.bd2)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestBenchSpy_compareGeneratorConfigs(t *testing.T) {
	baseConfig := &wasp.Config{
		GenName:  "gen",
		LoadType: wasp.VU,
		Schedule: []*wasp.Segment{
			{From: 1, Duration: time.Hour, Type: wasp.SegmentType_Plain},
		},
		CallTimeout: time.Second,
	}

	t.Run("identical configs", func(t *testing.T) {
		err := compareGeneratorConfigs(baseConfig, baseConfig)
		require.NoError(t, err)
	})

	t.Run("different generator names", func(t *testing.T) {
		cfg2 := *baseConfig
		cfg2.GenName = "gen2"
		err := compareGeneratorConfigs(baseConfig, &cfg2)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "generator names are different. Expected gen, got gen")
	})

	t.Run("different segment load", func(t *testing.T) {
		cfg2 := *baseConfig
		cfg2.Schedule = []*wasp.Segment{
			{From: 2, Duration: time.Hour, Type: wasp.SegmentType_Plain},
		}
		err := compareGeneratorConfigs(baseConfig, &cfg2)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "segments at index 0 are different")
	})

	t.Run("different segment duration", func(t *testing.T) {
		cfg2 := *baseConfig
		cfg2.Schedule = []*wasp.Segment{
			{From: 1, Duration: time.Minute, Type: wasp.SegmentType_Plain},
		}
		err := compareGeneratorConfigs(baseConfig, &cfg2)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "segments at index 0 are different")
	})

	t.Run("different segment type", func(t *testing.T) {
		cfg2 := *baseConfig
		cfg2.Schedule = []*wasp.Segment{
			{From: 1, Duration: time.Hour, Type: wasp.SegmentType_Steps},
		}
		err := compareGeneratorConfigs(baseConfig, &cfg2)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "segments at index 0 are different")
	})

	t.Run("different segment count", func(t *testing.T) {
		cfg2 := *baseConfig
		cfg2.Schedule = []*wasp.Segment{
			{From: 2, Duration: time.Hour},
			{From: 3, Duration: time.Hour},
		}
		err := compareGeneratorConfigs(baseConfig, &cfg2)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "schedules are different. Expected 1, got 2")
	})

	t.Run("different load types", func(t *testing.T) {
		cfg2 := *baseConfig
		cfg2.LoadType = wasp.RPS
		err := compareGeneratorConfigs(baseConfig, &cfg2)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "load types are different")
	})

	t.Run("different callTimeout", func(t *testing.T) {
		cfg2 := *baseConfig
		cfg2.CallTimeout = time.Minute
		err := compareGeneratorConfigs(baseConfig, &cfg2)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "call timeouts are different. Expected 1s, got 1m0s")
	})

	t.Run("different rate limit duration", func(t *testing.T) {
		cfg2 := *baseConfig
		cfg2.RateLimitUnitDuration = time.Minute
		err := compareGeneratorConfigs(baseConfig, &cfg2)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "rate limit unit durations are different. Expected 0s, got 1m0s")
	})
}

func TestBenchSpy_mustMarshallSegment(t *testing.T) {
	segment := &wasp.Segment{
		From:     1,
		Duration: time.Hour,
	}

	t.Run("successful marshal", func(t *testing.T) {
		assert.NotPanics(t, func() {
			result := mustMarshallSegment(segment)
			assert.Contains(t, result, `"from": 1`)
		})
	})
}
