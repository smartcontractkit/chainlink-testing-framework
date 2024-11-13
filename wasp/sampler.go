package wasp

import "math/rand"

type SamplerConfig struct {
	SuccessfulCallResultRecordRatio int
}

// Sampler is a CallResult filter that stores a percentage of successful call results
// errored and timed out results are always stored
type Sampler struct {
	cfg *SamplerConfig
}

// NewSampler creates a new Sampler instance using the provided SamplerConfig.
// If the provided configuration is nil, a default configuration with a
// SuccessfulCallResultRecordRatio of 100 is used. It returns a pointer to the
// newly created Sampler.
func NewSampler(cfg *SamplerConfig) *Sampler {
	if cfg == nil {
		cfg = &SamplerConfig{SuccessfulCallResultRecordRatio: 100}
	}
	return &Sampler{cfg: cfg}
}

// ShouldRecord determines whether a response should be recorded based on its error, failure, or timeout status.
// It updates the statistics for samples recorded or skipped. The decision is influenced by the configured
// ratio for recording successful call results, allowing for probabilistic recording of successful responses.
func (m *Sampler) ShouldRecord(cr *Response, s *Stats) bool {
	if cr.Error != "" || cr.Failed || cr.Timeout {
		s.SamplesRecorded.Add(1)
		return true
	}
	if m.cfg.SuccessfulCallResultRecordRatio == 0 {
		s.SamplesSkipped.Add(1)
		return false
	}
	if m.cfg.SuccessfulCallResultRecordRatio == 100 {
		s.SamplesRecorded.Add(1)
		return true
	}
	r := rand.Intn(100)
	if cr.Error == "" && r < m.cfg.SuccessfulCallResultRecordRatio {
		s.SamplesRecorded.Add(1)
		return true
	}
	s.SamplesSkipped.Add(1)
	return false
}
