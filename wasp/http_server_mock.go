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

func (s *HTTPMockServer) Run() {
	go func() {
		//nolint
		_ = s.srv.Run()
	}()
}

func (s *HTTPMockServer) URL() string {
	return "http://localhost:8080/1"
}

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
