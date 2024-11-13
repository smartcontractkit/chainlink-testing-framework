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

// Run starts the HTTPMockServer in a new goroutine. 
// It invokes the Run method of the underlying server without handling any errors. 
// This allows the server to operate concurrently with other processes.
func (s *HTTPMockServer) Run() {
	go func() {
		//nolint
		_ = s.srv.Run()
	}()
}

// URL returns the base URL of the HTTPMockServer as a string.
// It is typically used to construct full URLs for testing purposes.
func (s *HTTPMockServer) URL() string {
	return "http://localhost:8080/1"
}

// NewHTTPMockServer creates and returns a new HTTPMockServer instance configured with the provided HTTPMockServerConfig.
// If cfg is nil, it defaults to a configuration with 50ms latency and HTTP 200 status code for both API endpoints.
// The server has two endpoints, "/1" and "/2", each responding with a JSON message "pong" after the specified latency.
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
