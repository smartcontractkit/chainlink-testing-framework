package wasp

import (
	"context"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
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

// NewWSMockVU create a ws mock virtual user
func NewWSMockVU(cfg *WSMockVUConfig) *WSMockVU {
	return &WSMockVU{
		VUControl: NewVUControl(),
		cfg:       cfg,
		Data:      make([]string, 0),
	}
}

func (m *WSMockVU) Clone(_ *Generator) VirtualUser {
	return &WSMockVU{
		VUControl: NewVUControl(),
		cfg:       m.cfg,
		Data:      make([]string, 0),
	}
}

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

func (m *WSMockVU) Teardown(_ *Generator) error {
	return m.conn.Close(websocket.StatusInternalError, "")
}

// Call create a virtual user firing read requests against mock ws server
func (m *WSMockVU) Call(l *Generator) {
	startedAt := time.Now()
	v := map[string]string{}
	err := wsjson.Read(context.Background(), m.conn, &v)
	if err != nil {
		l.Log.Error().Err(err).Msg("failed read ws msg from vu")
	}
	l.ResponsesChan <- &Response{StartedAt: &startedAt, Data: v}
}
