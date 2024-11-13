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
// If the cfg parameter is nil, a default configuration is used, which sets 
// the SuccessfulCallResultRecordRatio to 100. The function returns a pointer 
// to the newly created Sampler.
func NewSampler(cfg *SamplerConfig) *Sampler {
	if cfg == nil {
		cfg = &SamplerConfig{SuccessfulCallResultRecordRatio: 100}
	}
	return &Sampler{cfg: cfg}
}

// ShouldRecord determines whether a response should be recorded based on its status and the configured recording ratio. 
// It returns true if the response indicates an error, failure, or timeout, or if it meets the criteria defined by the 
// SuccessfulCallResultRecordRatio configuration. If the response is successful and the random value is less than 
// the configured ratio, it will also return true. Otherwise, it returns false, indicating that the response should 
// not be recorded. The function also updates the provided Stats object to reflect the number of samples recorded 
// or skipped.
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
