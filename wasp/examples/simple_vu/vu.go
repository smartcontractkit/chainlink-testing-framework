package main

import (
	"context"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

type WSVirtualUser struct {
	*wasp.VUControl
	target string
	conn   *websocket.Conn
	Data   []string
}

func NewExampleWSVirtualUser(target string) *WSVirtualUser {
	return &WSVirtualUser{
		VUControl: wasp.NewVUControl(),
		target:    target,
		Data:      make([]string, 0),
	}
}

func (m *WSVirtualUser) Clone(_ *wasp.Generator) wasp.VirtualUser {
	return &WSVirtualUser{
		VUControl: wasp.NewVUControl(),
		target:    m.target,
		Data:      make([]string, 0),
	}
}

func (m *WSVirtualUser) Setup(l *wasp.Generator) error {
	var err error
	m.conn, _, err = websocket.Dial(context.Background(), m.target, &websocket.DialOptions{})
	if err != nil {
		l.Log.Error().Err(err).Msg("failed to connect from vu")
		//nolint
		_ = m.conn.Close(websocket.StatusInternalError, "")
		return err
	}
	return nil
}

func (m *WSVirtualUser) Teardown(_ *wasp.Generator) error {
	return nil
}

func (m *WSVirtualUser) Call(l *wasp.Generator) {
	startedAt := time.Now()
	v := map[string]string{}
	err := wsjson.Read(context.Background(), m.conn, &v)
	if err != nil {
		l.Log.Error().Err(err).Msg("failed read ws msg from vu")
		l.ResponsesChan <- &wasp.Response{StartedAt: &startedAt, Error: err.Error(), Failed: true}
		return
	}
	l.ResponsesChan <- &wasp.Response{StartedAt: &startedAt, Data: v}
}
