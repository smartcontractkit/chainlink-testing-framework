package wasp

import (
	"context"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

// WSMockVUConfig ws mock config
type WSMockVUConfig struct {
	TargetURl string
}

// WSMockVU ws mock virtual user
type WSMockVU struct {
	*VUControl
	cfg  *WSMockVUConfig
	conn *websocket.Conn
	Data []string
}

// NewWSMockVU initializes a WSMockVU with the provided configuration.
// It sets up control mechanisms and data storage, enabling the simulation
// of a WebSocket virtual user for testing scenarios.
func NewWSMockVU(cfg *WSMockVUConfig) *WSMockVU {
	return &WSMockVU{
		VUControl: NewVUControl(),
		cfg:       cfg,
		Data:      make([]string, 0),
	}
}

// Clone creates a new VirtualUser instance based on the current WSMockVU.
// It is used to instantiate additional virtual users for scaling load tests.
func (m *WSMockVU) Clone(_ *Generator) VirtualUser {
	return &WSMockVU{
		VUControl: NewVUControl(),
		cfg:       m.cfg,
		Data:      make([]string, 0),
	}
}

// Setup establishes a WebSocket connection to the configured target URL using the provided Generator.
// It returns an error if the connection cannot be established, allowing callers to handle setup failures.
func (m *WSMockVU) Setup(l *Generator) error {
	var err error
	m.conn, _, err = websocket.Dial(context.Background(), m.cfg.TargetURl, &websocket.DialOptions{})
	if err != nil {
		l.Log.Error().Err(err).Msg("failed to connect from virtual user")
		//nolint
		_ = m.conn.Close(websocket.StatusInternalError, "")
		return err
	}
	return nil
}

// Teardown gracefully closes the WebSocket connection for the VirtualUser.
// It should be called when the user simulation is complete to release resources.
func (m *WSMockVU) Teardown(_ *Generator) error {
	return m.conn.Close(websocket.StatusInternalError, "")
}

// Call reads a WebSocket message from the connection and sends the response with a timestamp to the generator's ResponsesChan.
// It is used by a virtual user to handle incoming WebSocket data during execution.
func (m *WSMockVU) Call(l *Generator) {
	startedAt := time.Now()
	v := map[string]string{}
	err := wsjson.Read(context.Background(), m.conn, &v)
	if err != nil {
		l.Log.Error().Err(err).Msg("failed read ws msg from vu")
	}
	l.ResponsesChan <- &Response{StartedAt: &startedAt, Data: v}
}
