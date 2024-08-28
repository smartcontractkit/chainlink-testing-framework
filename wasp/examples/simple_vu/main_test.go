package main

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"strings"
)

func TestVirtualUser(t *testing.T) {
	// start mock http server
	s := httptest.NewServer(wasp.MockWSServer{
		Sleep: 50 * time.Millisecond,
	})
	defer s.Close()
	time.Sleep(1 * time.Second)

	// define labels for differentiate one run from another
	labels := map[string]string{
		// check variables in dashboard/dashboard.go
		"go_test_name": "generator_healthcheck",
		"branch":       "generator_healthcheck",
		"commit":       "generator_healthcheck",
	}

	url := strings.Replace(s.URL, "http", "ws", -1)
	log.Warn().Interface("URL", url).Send()

	// create generator
	gen, err := wasp.NewGenerator(&wasp.Config{
		LoadType: wasp.VU,
		// just use plain line profile - 5 VUs for 10s
		Schedule:   wasp.Plain(5, 10*time.Second),
		VU:         NewExampleWSVirtualUser(url),
		Labels:     labels,
		LokiConfig: wasp.NewEnvLokiConfig(),
	})
	if err != nil {
		panic(err)
	}
	// run the generator and wait until it finish
	gen.Run(true)
}
