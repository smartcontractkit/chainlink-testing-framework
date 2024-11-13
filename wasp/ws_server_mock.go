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

// ServeHTTP handles HTTP requests by upgrading them to WebSocket connections.
// It accepts a WebSocket connection and continuously sends a constant answer
// until a normal closure is detected or an error occurs. Errors are logged
// using the server's logging function. The connection is closed with an
// internal error status if an error occurs during processing.
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

// constantAnswer sends a predefined JSON message with a constant answer to the given WebSocket connection after a specified sleep duration. 
// It returns any error encountered during the message writing process.
func constantAnswer(sleep time.Duration, c *websocket.Conn) error {
	time.Sleep(sleep)
	return wsjson.Write(context.Background(), c, map[string]string{
		"answer": "epico!",
	})
}
