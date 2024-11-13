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

// NewWSMockVU initializes a new instance of WSMockVU with the provided configuration. 
// It sets up the VUControl component and prepares an empty slice to hold data. 
// The returned WSMockVU instance is ready for use in a WebSocket mock environment.
func NewWSMockVU(cfg *WSMockVUConfig) *WSMockVU {
	return &WSMockVU{
		VUControl: NewVUControl(),
		cfg:       cfg,
		Data:      make([]string, 0),
	}
}

// Clone creates a new instance of WSMockVU by copying the configuration from the original instance. 
// It initializes a new VUControl and sets up an empty slice for Data. 
// This function is typically used to generate multiple virtual users with the same configuration in a load testing scenario. 
// The returned VirtualUser can be utilized independently of the original instance.
func (m *WSMockVU) Clone(_ *Generator) VirtualUser {
	return &WSMockVU{
		VUControl: NewVUControl(),
		cfg:       m.cfg,
		Data:      make([]string, 0),
	}
}

// Setup establishes a WebSocket connection to the target URL specified in the configuration. 
// It returns an error if the connection fails, logging the error details for troubleshooting. 
// If the connection is successful, it initializes the connection for further use.
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

// Teardown closes the WebSocket connection associated with the WSMockVU instance. 
// It returns an error if the connection cannot be closed successfully. 
// This function is typically called during the teardown phase of a virtual user's lifecycle, 
// allowing for proper cleanup of resources.
func (m *WSMockVU) Teardown(_ *Generator) error {
	return m.conn.Close(websocket.StatusInternalError, "")
}

// Call reads a WebSocket message from the virtual user connection and sends the result to the provided generator's response channel. 
// It logs an error if reading the message fails. The function captures the time when the call started and includes this timestamp 
// along with the received data in the response sent to the channel.
func (m *WSMockVU) Call(l *Generator) {
	startedAt := time.Now()
	v := map[string]string{}
	err := wsjson.Read(context.Background(), m.conn, &v)
	if err != nil {
		l.Log.Error().Err(err).Msg("failed read ws msg from vu")
	}
	l.ResponsesChan <- &Response{StartedAt: &startedAt, Data: v}
}
