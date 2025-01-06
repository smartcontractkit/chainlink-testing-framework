package wasp

import (
	"context"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

/* This is a naive WS mock server to check the tool performance */

type MockWSServer struct {
	// Logf controls where logs are sent.
	Logf  func(f string, v ...interface{})
	Sleep time.Duration
}

// ServeHTTP upgrades the HTTP connection to a WebSocket and continuously sends predefined responses.
// It is used to mock WebSocket server behavior for testing purposes.
func (s MockWSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		s.Logf("%v", err)
		return
	}
	// nolint
	defer c.Close(websocket.StatusInternalError, "")
	for {
		//nolint
		err = constantAnswer(s.Sleep, c)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			s.Logf("failed to answer with %v: %v", r.RemoteAddr, err)
			return
		}
	}
}

// constantAnswer sends a predefined "epico!" answer to the websocket connection.
// It is used to consistently respond to clients with a fixed message.
func constantAnswer(sleep time.Duration, c *websocket.Conn) error {
	time.Sleep(sleep)
	return wsjson.Write(context.Background(), c, map[string]string{
		"answer": "epico!",
	})
}
