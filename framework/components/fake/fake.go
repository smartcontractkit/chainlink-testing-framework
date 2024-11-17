package fake

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
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

func JSON(path string, response map[string]any, statusCode int) error {
	if Service == nil {
		return fmt.Errorf("mock service is not initialized, please set up NewFakeDataProvider in your tests")
	}
	Service.Any(path, func(c *gin.Context) {
		c.JSON(statusCode, response)
	})
	return nil
}

// NewFakeDataProvider creates new fake data provider
func NewFakeDataProvider(in *Input) (*Output, error) {
	Service = gin.Default()
	Service.Use(recordMiddleware())
	go func() {
		_ = Service.Run(fmt.Sprintf(":%d", in.Port))
	}()
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
