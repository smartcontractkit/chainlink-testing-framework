package wasp

import (
	"context"
	"net/http"
	"time"

	"nhooyr.io/websocket/wsjson"

	"nhooyr.io/websocket"
)

/* This is a naive WS mock server to check the tool performance */

type MockWSServer struct {
	// Logf controls where logs are sent.
	Logf  func(f string, v ...interface{})
	Sleep time.Duration
}

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

func constantAnswer(sleep time.Duration, c *websocket.Conn) error {
	time.Sleep(sleep)
	return wsjson.Write(context.Background(), c, map[string]string{
		"answer": "epico!",
	})
}
