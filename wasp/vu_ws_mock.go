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

// NewWSMockVU creates a new WSMockVU instance using the provided WSMockVUConfig.
// It initializes the VUControl and prepares the Data slice for storing mock data.
func NewWSMockVU(cfg *WSMockVUConfig) *WSMockVU {
	return &WSMockVU{
		VUControl: NewVUControl(),
		cfg:       cfg,
		Data:      make([]string, 0),
	}
}

// Clone creates and returns a new VirtualUser by duplicating the WSMockVU's configuration.
func (m *WSMockVU) Clone(_ *Generator) VirtualUser {
	return &WSMockVU{
		VUControl: NewVUControl(),
		cfg:       m.cfg,
		Data:      make([]string, 0),
	}
}

// Setup initializes the WebSocket connection for the WSMockVU using the provided Generator.
// It attempts to dial the target URL specified in the configuration.
// If the connection fails, it logs the error, closes the connection, and returns the encountered error.
// On successful connection, it returns nil.
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

// Teardown closes the virtual user's connection and performs necessary cleanup.
// It returns an error if the connection fails to close properly.
func (m *WSMockVU) Teardown(_ *Generator) error {
	return m.conn.Close(websocket.StatusInternalError, "")
}

// Call reads a WebSocket message from the WSMockVU's connection and sends a Response containing the start time and data to the provided Generator.
// If the read operation fails, it logs the encountered error.
func (m *WSMockVU) Call(l *Generator) {
	startedAt := time.Now()
	v := map[string]string{}
	err := wsjson.Read(context.Background(), m.conn, &v)
	if err != nil {
		l.Log.Error().Err(err).Msg("failed read ws msg from vu")
	}
	l.ResponsesChan <- &Response{StartedAt: &startedAt, Data: v}
}
