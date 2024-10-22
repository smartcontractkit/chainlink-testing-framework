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

// NewSampler creates new Sampler
func NewSampler(cfg *SamplerConfig) *Sampler {
	if cfg == nil {
		cfg = &SamplerConfig{SuccessfulCallResultRecordRatio: 100}
	}
	return &Sampler{cfg: cfg}
}

// ShouldRecord return true if we should save CallResult
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
