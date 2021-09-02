package experiments

import "time"

type DNSChaos struct {
	Base
	Duration time.Duration
	Patterns []string
}

func (e *DNSChaos) SetBase(base Base) {
	e.Base = base
}

func (e *DNSChaos) Resource() string {
	return "dnschaos"
}

func (e *DNSChaos) Filename() string {
	return "dns-chaos.yml"
}
