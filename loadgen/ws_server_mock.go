package loadgen

import (
	"context"
	"net/http"
	"time"

	"nhooyr.io/websocket/wsjson"

	"nhooyr.io/websocket"
)

/* This is a naive WS mock server to check the tool performance */

type MockWSServer struct {
	// logf controls where logs are sent.
	logf  func(f string, v ...interface{})
	sleep time.Duration
}

func (s MockWSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		s.logf("%v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")
	for {
		//nolint
		err = constantAnswer(s.sleep, c)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			s.logf("failed to constantAnswer with %v: %v", r.RemoteAddr, err)
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
