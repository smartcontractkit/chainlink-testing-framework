package loadgen

import (
	"time"

	"github.com/gin-gonic/gin"
)

type HTTPMockServer struct {
	srv   *gin.Engine
	Sleep time.Duration
}

func (s *HTTPMockServer) Run() {
	go func() {
		//nolint
		s.srv.Run()
	}()
}

func NewHTTPMockServer(sleep time.Duration) *HTTPMockServer {
	srv := gin.New()
	gin.SetMode(gin.ReleaseMode)
	srv.GET("/", func(c *gin.Context) {
		time.Sleep(sleep)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return &HTTPMockServer{srv: srv}
}
