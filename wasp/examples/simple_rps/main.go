package main

import (
	"testing"
	"time"

	"github.com/smartcontractkit/wasp"
)

func TestGun(t *testing.T) {
	// start mock http server
	srv := wasp.NewHTTPMockServer(nil)
	srv.Run()

	// define labels for differentiate one run from another
	labels := map[string]string{
		// check variables in dashboard/dashboard.go
		"go_test_name": "generator_healthcheck",
		"gen_name":     "generator_healthcheck",
		"branch":       "generator_healthcheck",
		"commit":       "generator_healthcheck",
	}

	// create generator
	gen, err := wasp.NewGenerator(&wasp.Config{
		LoadType: wasp.RPS,
		// just use plain line profile - 5 RPS for 10s
		Schedule:   wasp.Plain(5, 60*time.Second),
		Gun:        NewExampleHTTPGun(srv.URL()),
		Labels:     labels,
		LokiConfig: wasp.NewEnvLokiConfig(),
	})
	if err != nil {
		panic(err)
	}
	// run the generator and wait until it finish
	gen.Run(true)
}
