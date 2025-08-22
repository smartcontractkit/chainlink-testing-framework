package fake

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

const (
	DefaultFakeServicePort = 9111
)

type Input struct {
	Image string  `toml:"image"`
	Port  int     `toml:"port" validate:"required"`
	Out   *Output `toml:"out"`
}

type Output struct {
	UseCache      bool   `toml:"use_cache"`
	BaseURLHost   string `toml:"base_url_host"`
	BaseURLDocker string `toml:"base_url_docker"`
}

var (
	Service     *gin.Engine
	server      *http.Server
	validMethod = regexp.MustCompile("GET|POST|PATCH|PUT|DELETE")
)

// NewFakeDataProvider creates new fake data provider
func NewFakeDataProvider(in *Input) (*Output, error) {
	if server != nil {
		return nil, fmt.Errorf("fake service is already running, call TerminateService first")
	}

	Service = gin.Default()
	Service.Use(recordMiddleware())

	server = &http.Server{
		Addr:    fmt.Sprintf(":%d", in.Port),
		Handler: Service,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Log error but don't panic - server might be intentionally shut down
			fmt.Printf("Fake service error: %v\n", err)
		}
	}()

	out := &Output{
		BaseURLHost:   fmt.Sprintf("http://localhost:%d", in.Port),
		BaseURLDocker: fmt.Sprintf("%s:%d", framework.HostDockerInternal(), in.Port),
	}
	in.Out = out
	return out, nil
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
func Func(method, path string, f func(ctx *gin.Context)) error {
	if err := validate(method, path); err != nil {
		return err
	}
	Service.Handle(method, path, f)
	return nil
}

// JSON fakes for method, path, response and status code
func JSON(method, path string, response map[string]any, statusCode int) error {
	if err := validate(method, path); err != nil {
		return err
	}
	Service.Handle(method, path, func(c *gin.Context) {
		c.JSON(statusCode, response)
	})
	return nil
}

// TerminateService synchronously shuts down the fake service
func TerminateService(ctx context.Context) error {
	if server == nil {
		return nil
	}
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown fake service gracefully: %w", err)
	}
	Service = nil
	server = nil
	return nil
}

// IsServiceRunning returns true if the fake service is currently running
func IsServiceRunning() bool {
	return server != nil && Service != nil
}

// HostDockerInternal returns host.docker.internal that works both locally and in GHA
func HostDockerInternal() string {
	if os.Getenv("CI") == "true" {
		return "http://172.17.0.1"
	}
	return "http://host.docker.internal"
}
