package resty

import (
	"github.com/go-resty/resty/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/utils/header"
	"os"
)

func NewDefaultResty() (*resty.Client, error) {
	h, err := header.ReadEnvHTTPHeaders(logging.L)
	if err != nil {
		return nil, err
	}
	return resty.New().
		SetDebug(os.Getenv("DEBUG_RESTY") == "true").
		SetHeaders(header.HeaderToMultiValueFormat(h)), nil
}
