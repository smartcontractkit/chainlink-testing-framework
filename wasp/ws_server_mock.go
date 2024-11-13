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

// ServeHTTP handles WebSocket connections by accepting incoming requests and establishing a WebSocket connection. 
// It continuously processes messages until the connection is closed normally or an error occurs. 
// If an error occurs during the connection setup or message processing, it logs the error and terminates the connection. 
// The function ensures that the WebSocket connection is properly closed when finished.
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

// constantAnswer pauses for the specified duration before sending a predefined response 
// over the provided WebSocket connection. It returns an error if the write operation fails. 
// The response sent is a JSON object containing the key "answer" with the value "epico!". 
// If the connection is closed normally, the function will return a status indicating that.
func constantAnswer(sleep time.Duration, c *websocket.Conn) error {
	time.Sleep(sleep)
	return wsjson.Write(context.Background(), c, map[string]string{
		"answer": "epico!",
	})
}
