package fake

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

var (
	Service     *gin.Engine
	validMethod = regexp.MustCompile("GET|POST|PATCH|PUT|DELETE")
)

type Input struct {
	Port int     `toml:"port" validate:"required"`
	Out  *Output `toml:"out"`
}

type Output struct {
	BaseURLHost   string `toml:"base_url_host"`
	BaseURLDocker string `toml:"base_url_docker"`
}

// validate validates method and path, does not allow to override mock
func validate(method, path string) error {
	if Service == nil {
		return fmt.Errorf("mock service is not initialized, please set up NewFakeDataProvider in your tests")
	}
	if match := validMethod.Match([]byte(method)); !match {
		return fmt.Errorf("provide GET, POST, PATCH, PUT or DELETE in fake.JSON() method")
	}
	if _, ok := R.Data[RecordKey(method, path)]; ok {
		return fmt.Errorf("fake with method %s and path %s already exists", method, path)
	}
	R.Data[RecordKey(method, path)] = make([]*Record, 0)
	return nil
}

// Func fakes method and path with a custom func
func Func(method string, path string, f func(ctx *gin.Context)) error {
	if err := validate(method, path); err != nil {
		return err
	}
	Service.Handle(method, path, f)
	return nil
}

// JSON fakes for method, path, response and status code
func JSON(method string, path string, response map[string]any, statusCode int) error {
	if err := validate(method, path); err != nil {
		return err
	}
	Service.Handle(method, path, func(c *gin.Context) {
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
		BaseURLHost:   fmt.Sprintf("http://localhost:%d", in.Port),
		BaseURLDocker: fmt.Sprintf("%s:%d", framework.HostDockerInternal(), in.Port),
	}
	in.Out = out
	return out, nil
}
