package generators

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

type CLNodeGun struct {
	Mode   string
	client *clclient.ChainlinkClient
	Data   []string
}

func NewCLNodeGun(c *clclient.ChainlinkClient, mode string) *CLNodeGun {
	return &CLNodeGun{
		Mode:   mode,
		client: c,
		Data:   make([]string, 0),
	}
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *CLNodeGun) Call(l *wasp.Generator) *wasp.Response {
	switch m.Mode {
	case "bridges":
		return m.bridges()
	default:
		panic("unknown generator mode")
	}
	return nil
}

func (m *CLNodeGun) bridges() *wasp.Response {
	b, rr, err := m.client.ReadBridges()
	if b == nil {
		return &wasp.Response{Error: "bridges response is nil"}
	}
	if err != nil {
		return &wasp.Response{Error: err.Error(), Failed: true}
	}
	if rr.Status() != "200 OK" {
		return &wasp.Response{Error: "not 200", Failed: true}
	}
	return &wasp.Response{Data: b}
}
