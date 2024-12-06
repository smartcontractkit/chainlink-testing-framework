package wasp

import (
	"time"

	"github.com/gin-gonic/gin"
)

type HTTPMockServerConfig struct {
	FirstAPILatency   time.Duration
	FirstAPIHTTPCode  int
	SecondAPILatency  time.Duration
	SecondAPIHTTPCode int
}

type HTTPMockServer struct {
	srv   *gin.Engine
	Sleep time.Duration
}

// Run starts the HTTPMockServer in a separate goroutine.
// It enables the server to handle incoming HTTP requests concurrently.
func (s *HTTPMockServer) Run() {
	go func() {
		//nolint
		_ = s.srv.Run()
	}()
}

// URL returns the base URL of the HTTPMockServer.
// Use it to configure clients to send requests to the mock server during testing.
func (s *HTTPMockServer) URL() string {
	return "http://localhost:8080/1"
}

// NewHTTPMockServer initializes an HTTP mock server with configurable latencies and response codes.
// If cfg is nil, default settings are applied.
// Use it to simulate HTTP endpoints for testing purposes.
func NewHTTPMockServer(cfg *HTTPMockServerConfig) *HTTPMockServer {
	if cfg == nil {
		cfg = &HTTPMockServerConfig{
			FirstAPILatency:   50 * time.Millisecond,
			FirstAPIHTTPCode:  200,
			SecondAPILatency:  50 * time.Millisecond,
			SecondAPIHTTPCode: 200,
		}
	}
	srv := gin.New()
	gin.SetMode(gin.ReleaseMode)
	srv.GET("/1", func(c *gin.Context) {
		time.Sleep(cfg.FirstAPILatency)
		c.JSON(cfg.FirstAPIHTTPCode, gin.H{
			"message": "pong",
		})
	})
	srv.GET("/2", func(c *gin.Context) {
		time.Sleep(cfg.SecondAPILatency)
		c.JSON(cfg.SecondAPIHTTPCode, gin.H{
			"message": "pong",
		})
	})
	return &HTTPMockServer{srv: srv}
}
