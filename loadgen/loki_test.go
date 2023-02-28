package loadgen

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/client"

	"github.com/stretchr/testify/require"
)

type ResponseSample struct {
	Data string
}

type StatsSample struct {
	CallTimeout      float64
	CurrentInstances float64
	CurrentRPS       float64
	RunFailed        bool
	RunStopped       bool
	Success          float64
	Failed           float64
}

type LokiSamplesAssertions struct {
	ResponsesSamples []ResponseSample
	StatsSamples     []StatsSample
}

func assertSamples(t *testing.T, samples []client.PromtailSendResult, a LokiSamplesAssertions) {
	var cd CallResult
	for i, s := range samples[0:2] {
		t.Logf("Entry: %s", s.Entry)
		err := json.Unmarshal([]byte(s.Entry), &cd)
		require.NoError(t, err)
		require.NotEmpty(t, cd.Duration)
		require.Equal(t, cd.Data, a.ResponsesSamples[i].Data)
	}
	// marshal to map because atomic can't be marshalled
	var ls map[string]interface{}
	for i, s := range samples[2:4] {
		t.Logf("Stats: %s", s.Entry)
		err := json.Unmarshal([]byte(s.Entry), &ls)
		require.NoError(t, err)
		require.Equal(t, ls["callTimeout"], a.StatsSamples[i].CallTimeout)
		require.Equal(t, ls["current_instances"], a.StatsSamples[i].CurrentInstances)
		require.Equal(t, ls["current_rps"], a.StatsSamples[i].CurrentRPS)
		require.Equal(t, ls["run_failed"], a.StatsSamples[i].RunFailed)
		require.Equal(t, ls["run_stopped"], a.StatsSamples[i].RunStopped)
		require.Equal(t, ls["success"], a.StatsSamples[i].Success)
		require.Equal(t, ls["failed"], a.StatsSamples[i].Failed)
	}
}

func TestLokiSamples(t *testing.T) {
	defaultLabels := map[string]string{
		"cluster":    "test_cluster",
		"namespace":  "test_namespace",
		"app":        "test_app",
		"test_group": "test_group",
		"test_id":    "test_id",
	}

	type test struct {
		name       string
		genCfg     *LoadGeneratorConfig
		assertions LokiSamplesAssertions
	}

	tests := []test{
		{
			name: "successful RPS run should contain at least 2 response samples without errors and 2 stats samples",
			genCfg: &LoadGeneratorConfig{
				T: t,
				// empty URL is a special case for mocked client
				LokiConfig: client.NewDefaultLokiConfig("", ""),
				Labels:     defaultLabels,
				Duration:   55 * time.Millisecond,
				Schedule: &LoadSchedule{
					Type:      RPSScheduleType,
					StartFrom: 1,
				},
				Gun: NewMockGun(&MockGunConfig{
					CallSleep: 50 * time.Millisecond,
				}),
			},
			assertions: LokiSamplesAssertions{
				ResponsesSamples: []ResponseSample{
					{
						Data: "successCallData",
					},
					{
						Data: "successCallData",
					},
				},
				StatsSamples: []StatsSample{
					{
						CallTimeout:      0,
						CurrentInstances: 1,
						CurrentRPS:       1,
						RunFailed:        false,
						RunStopped:       false,
						Success:          2,
						Failed:           0,
					},
					{
						CallTimeout:      0,
						CurrentInstances: 1,
						CurrentRPS:       1,
						RunFailed:        false,
						RunStopped:       false,
						Success:          2,
						Failed:           0,
					},
				},
			}},
	}

	for _, tc := range tests {
		gen, err := NewLoadGenerator(tc.genCfg)
		require.NoError(t, err)
		gen.Run()
		gen.Wait()
		assertSamples(t, gen.loki.AllHandleResults(), tc.assertions)
	}

	t.Parallel()
}
