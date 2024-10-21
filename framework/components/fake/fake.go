package fake

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"net/http"
	"os"
)

var (
	Service *gin.Engine
)

type Input struct {
	Port int     `toml:"port" validate:"required"`
	Out  *Output `toml:"out"`
}

type Output struct {
	BaseURLHost   string `toml:"base_url_host"`
	BaseURLDocker string `toml:"base_url_docker"`
}

func FakeJSON(path string, response gin.H, statusCode int) error {
	if Service == nil {
		return fmt.Errorf("mock service is not initialized, please set up NewFakeDataProvider in your tests")
	}
	Service.Any(path, func(c *gin.Context) {
		c.JSON(statusCode, response)
	})
	return nil
}

func runMocks(in *Input) {
	router := gin.Default()
	router.GET("/mock1", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "This is a GET request response from the mock service.",
		})
	})
	router.POST("/mock2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "This is a POST request response from the mock service.",
		})
	})
	_ = router.Run(fmt.Sprintf(":%d", in.Port))
}

func NewFakeDataProvider(in *Input) (*Output, error) {
	go runMocks(in)
	out := &Output{
		BaseURLHost: fmt.Sprintf("http://localhost:%d", in.Port),
	}
	if os.Getenv(framework.EnvVarCI) == "true" {
		out.BaseURLDocker = fmt.Sprintf("http://172.17.0.1:%d", in.Port)
	} else {
		out.BaseURLDocker = fmt.Sprintf("http://host.docker.internal:%d", in.Port)
	}
	in.Out = out
	return out, nil
}
