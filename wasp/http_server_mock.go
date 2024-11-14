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

// Run starts the HTTPMockServer by launching its internal server in a new goroutine.
// It invokes s.srv.Run() asynchronously and ignores any errors.
func (s *HTTPMockServer) Run() {
	go func() {
		//nolint
		_ = s.srv.Run()
	}()
}

// URL returns the base URL of the HTTP mock server.
func (s *HTTPMockServer) URL() string {
	return "http://localhost:8080/1"
}

// NewHTTPMockServer creates and initializes a new HTTPMockServer.
// If cfg is nil, it sets default latencies and HTTP status codes for two endpoints.
// The server defines two GET endpoints (/1 and /2) that respond with a JSON message after a simulated delay.
// It returns a pointer to the configured HTTPMockServer.
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
