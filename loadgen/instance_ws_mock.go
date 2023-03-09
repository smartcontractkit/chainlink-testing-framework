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
	cfg  WSMockConfig
	stop chan struct{}
	Data []string
}

// NewWSMockInstance create a ws mock instanceTemplate
func NewWSMockInstance(cfg WSMockConfig) WSMockInstance {
	return WSMockInstance{
		cfg:  cfg,
		stop: make(chan struct{}, 1),
		Data: make([]string, 0),
	}
}

// Run create an instanceTemplate firing read requests against mock ws server
func (m WSMockInstance) Run(l *Generator) {
	l.ResponsesWaitGroup.Add(1)
	c, _, err := websocket.Dial(context.Background(), m.cfg.TargetURl, &websocket.DialOptions{})
	if err != nil {
		l.Log.Error().Err(err).Msg("failed to connect from instanceTemplate")
		//nolint
		c.Close(websocket.StatusInternalError, "")
	}
	go func() {
		defer l.ResponsesWaitGroup.Done()
		for {
			select {
			case <-l.ResponsesCtx.Done():
				//nolint
				c.Close(websocket.StatusNormalClosure, "")
				return
			case <-m.stop:
				return
			default:
				startedAt := time.Now()
				v := map[string]string{}
				err = wsjson.Read(context.Background(), c, &v)
				if err != nil {
					l.Log.Error().Err(err).Msg("failed read ws msg from instanceTemplate")
				}
				l.ResponsesChan <- CallResult{StartedAt: &startedAt, Data: v}
			}
		}
	}()
}

func (m WSMockInstance) Stop(l *Generator) {
	m.stop <- struct{}{}
}
