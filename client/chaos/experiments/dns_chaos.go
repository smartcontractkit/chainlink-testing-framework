package experiments

import "time"

// DNSChaos stuct with objects for DNS chaos testing
type DNSChaos struct {
	Base
	Duration time.Duration
	Patterns []string
}

// SetBase sets the base
func (e *DNSChaos) SetBase(base Base) {
	e.Base = base
}

// Resource returns the resource for dns chaos
func (e *DNSChaos) Resource() string {
	return "dnschaos"
}

// Filename returns the file name of the dns chaos yaml
func (e *DNSChaos) Filename() string {
	return "dns-chaos.yml"
}
