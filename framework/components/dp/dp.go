package dp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"net/http"
)

var (
	MockService *gin.Engine
)

type Input struct {
	Port int     `toml:"port" validate:"required"`
	Out  *Output `toml:"out"`
}

type Output struct {
	Urls []string `toml:"data_provider_urls"`
}

func Mock(path string, response gin.H, statusCode int) error {
	if MockService == nil {
		return fmt.Errorf("mock service is not initialized, please set up NewMockedDataProvider in your tests")
	}
	MockService.Any(path, func(c *gin.Context) {
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

func NewMockedDataProvider(in *Input) (*Output, error) {
	if in.Out != nil && framework.NoCache() {
		return in.Out, nil
	}
	go runMocks(in)
	out := &Output{
		Urls: []string{
			"http://localhost:8080/mock1",
			"http://localhost:8080/mock2",
		},
	}
	in.Out = out
	return out, nil
}
