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

// NewWSMockVU creates a new instance of WSMockVU with the provided configuration.
// It initializes the VUControl using NewVUControl and sets up an empty data slice.
func NewWSMockVU(cfg *WSMockVUConfig) *WSMockVU {
	return &WSMockVU{
		VUControl: NewVUControl(),
		cfg:       cfg,
		Data:      make([]string, 0),
	}
}

// Clone creates a new instance of WSMockVU with a fresh VUControl and the same configuration as the original.
// It initializes the Data field as an empty string slice and returns the new VirtualUser instance.
func (m *WSMockVU) Clone(_ *Generator) VirtualUser {
	return &WSMockVU{
		VUControl: NewVUControl(),
		cfg:       m.cfg,
		Data:      make([]string, 0),
	}
}

// Setup establishes a WebSocket connection to the target URL specified in the configuration.
// It returns an error if the connection attempt fails, logging the error and closing the connection.
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

// Teardown closes the WebSocket connection of the WSMockVU instance with an internal error status.
// It is used to clean up resources associated with the virtual user after execution.
// The function returns any error encountered during the closure of the connection.
func (m *WSMockVU) Teardown(_ *Generator) error {
	return m.conn.Close(websocket.StatusInternalError, "")
}

// Call reads a WebSocket message from the virtual user connection and sends the response data to the Generator's ResponsesChan. 
// It logs an error if reading the message fails. The response includes the time the call started and the message data.
func (m *WSMockVU) Call(l *Generator) {
	startedAt := time.Now()
	v := map[string]string{}
	err := wsjson.Read(context.Background(), m.conn, &v)
	if err != nil {
		l.Log.Error().Err(err).Msg("failed read ws msg from vu")
	}
	l.ResponsesChan <- &Response{StartedAt: &startedAt, Data: v}
}
