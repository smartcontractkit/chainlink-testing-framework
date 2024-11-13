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

// Run starts the HTTP mock server in a separate goroutine. 
// It invokes the server's Run method, allowing it to handle incoming requests 
// asynchronously. This function does not return any value and is intended 
// to be called to initiate the server's operation.
func (s *HTTPMockServer) Run() {
	go func() {
		//nolint
		_ = s.srv.Run()
	}()
}

// URL returns the base URL of the HTTP mock server as a string. 
// This URL can be used to send requests to the mock server for testing purposes. 
// The returned URL is fixed and points to "http://localhost:8080/1".
func (s *HTTPMockServer) URL() string {
	return "http://localhost:8080/1"
}

// NewHTTPMockServer creates a new instance of HTTPMockServer with the provided configuration. 
// If the configuration is nil, it initializes the server with default settings, including 
// predefined latencies and HTTP response codes for two mock API endpoints. 
// The returned HTTPMockServer can be used to simulate API responses for testing purposes.
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
