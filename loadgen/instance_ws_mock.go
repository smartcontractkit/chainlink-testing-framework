package loadgen

import (
	"context"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// WSMockConfig ws mock config
type WSMockConfig struct {
	TargetURl string
}

// WSMockInstance ws mock mock
type WSMockInstance struct {
	cfg  *WSMockConfig
	Data []string
}

// NewWSMockInstance create a ws mock instance
func NewWSMockInstance(cfg *WSMockConfig) *WSMockInstance {
	return &WSMockInstance{
		cfg:  cfg,
		Data: make([]string, 0),
	}
}

// Run create an instance firing read requests against mock ws server
func (m *WSMockInstance) Run(l *Generator) {
	l.ResponsesWaitGroup.Add(1)
	c, _, err := websocket.Dial(context.Background(), m.cfg.TargetURl, &websocket.DialOptions{})
	if err != nil {
		l.Log.Error().Err(err).Msg("failed to connect from instance")
		//nolint
		c.Close(websocket.StatusInternalError, "")
	}
	go func() {
		for {
			select {
			case <-l.ResponsesCtx.Done():
				l.ResponsesWaitGroup.Done()
				//nolint
				c.Close(websocket.StatusNormalClosure, "")
				return
			default:
				startedAt := time.Now()
				v := map[string]string{}
				err = wsjson.Read(context.Background(), c, &v)
				if err != nil {
					l.Log.Error().Err(err).Msg("failed read ws msg from instance")
				}
				l.ResponsesChan <- CallResult{StartedAt: &startedAt, Data: v}
			}
		}
	}()
}
