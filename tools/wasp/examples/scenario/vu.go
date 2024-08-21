package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/smartcontractkit/wasp"
	"go.uber.org/ratelimit"
)

const (
	GroupAuth   = "auth"
	GroupUser   = "user"
	GroupCommon = "common"
)

type VirtualUser struct {
	*wasp.VUControl
	target string
	Data   []string
	rl     ratelimit.Limiter
	client *resty.Client
	stop   chan struct{}
}

func NewExampleScenario(target string) *VirtualUser {
	return &VirtualUser{
		VUControl: wasp.NewVUControl(),
		target:    target,
		rl:        ratelimit.New(10),
		client:    resty.New().SetBaseURL(target),
		Data:      make([]string, 0),
	}
}

func (m *VirtualUser) Clone(_ *wasp.Generator) wasp.VirtualUser {
	return &VirtualUser{
		VUControl: wasp.NewVUControl(),
		target:    m.target,
		rl:        ratelimit.New(10),
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
