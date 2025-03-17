package main

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/ratelimit"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

const (
	GroupAuth = "auth"
	GroupUser = "user"
)

type VirtualUser struct {
	*wasp.VUControl
	target    string
	rateLimit int
	rl        ratelimit.Limiter
	Data      []string
	client    *resty.Client
}

func NewExampleScenario(target string) *VirtualUser {
	rateLimit := 10
	return &VirtualUser{
		VUControl: wasp.NewVUControl(),
		target:    target,
		rateLimit: rateLimit,
		rl:        ratelimit.New(rateLimit, ratelimit.WithoutSlack),
		client:    resty.New().SetBaseURL(target),
		Data:      make([]string, 0),
	}
}

func (m *VirtualUser) Clone(_ *wasp.Generator) wasp.VirtualUser {
	return &VirtualUser{
		VUControl: wasp.NewVUControl(),
		target:    m.target,
		rateLimit: m.rateLimit,
		rl:        ratelimit.New(m.rateLimit, ratelimit.WithoutSlack),
		client:    resty.New().SetBaseURL(m.target),
		Data:      make([]string, 0),
	}
}

func (m *VirtualUser) Setup(_ *wasp.Generator) error {
	return nil
}

func (m *VirtualUser) Teardown(_ *wasp.Generator) error {
	return nil
}

// represents user login
func (m *VirtualUser) requestOne(l *wasp.Generator) {
	var result map[string]interface{}
	r, err := m.client.R().
		SetResult(&result).
		Get(m.target)
	if err != nil {
		l.Responses.Err(r, GroupAuth, err)
		return
	}
	l.Responses.OK(r, GroupAuth)
}

// represents authenticated user action
func (m *VirtualUser) requestTwo(l *wasp.Generator) {
	var result map[string]interface{}
	r, err := m.client.R().
		SetResult(&result).
		Get(m.target)
	if err != nil {
		l.Responses.Err(r, GroupUser, err)
		return
	}
	l.Responses.OK(r, GroupUser)
}

func (m *VirtualUser) Call(l *wasp.Generator) {
	m.rl.Take()
	m.requestOne(l)
	m.requestTwo(l)
}
