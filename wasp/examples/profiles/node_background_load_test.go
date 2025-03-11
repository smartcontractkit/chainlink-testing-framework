package main

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

func TestBackgroundLoadSimple(t *testing.T) {
	srv := wasp.NewHTTPMockServer(nil)
	srv.Run()

	labels := map[string]string{
		"branch": "background_load_simple",
		"commit": "background_load_simple",
	}

	zetaSchedule := wasp.Combine(
		wasp.Steps(1, 1, 4, 20*time.Second),
		wasp.Plain(5, 50*time.Second),
		wasp.Steps(5, -1, 4, 20*time.Second))

	rpsProfile := wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			T:          t,
			LoadType:   wasp.RPS,
			GenName:    "Zeta",
			Schedule:   zetaSchedule,
			Gun:        NewExampleHTTPGun(srv.URL()),
			Labels:     labels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		}))

	// start load generation without waiting for it to finish
	_, err := rpsProfile.Run(false)

	// Wait for the first generator enter the steady state
	time.Sleep(20 * time.Second)

	etaSchedule := wasp.Combine(
		wasp.Steps(1, 1, 10, 20*time.Second),
		wasp.Plain(10, 20*time.Second),
		wasp.Steps(10, -1, 10, 10*time.Second))

	iotaSchedule := wasp.Combine(
		wasp.Steps(1, 1, 10, 20*time.Second),
		wasp.Plain(10, 20*time.Second),
		wasp.Steps(10, -1, 10, 10*time.Second))

	vuProfile, err := wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			T:          t,
			LoadType:   wasp.VU,
			GenName:    "Eta",
			Schedule:   etaSchedule,
			VU:         NewExampleScenario(srv.URL()),
			Labels:     labels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).
		Add(wasp.NewGenerator(&wasp.Config{
			T:          t,
			LoadType:   wasp.VU,
			GenName:    "Iota",
			Schedule:   iotaSchedule,
			VU:         NewExampleScenario(srv.URL()),
			Labels:     labels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).
		Run(true)
	// check if VU Profile did not return an error (e.g. due to invalid configuration or alerts triggered)
	require.NoError(t, err)

	// wait until RPS Profile finishes
	rpsProfile.Wait()
	// check if RPS Profile did not return an error (e.g. due to invalid configuration or alerts triggered)
	require.NoError(t, err)

	require.Equal(t, 0, len(rpsProfile.Generators[0].Errors()), "RPS generator had errors errors")
	require.Equal(t, 0, len(vuProfile.Generators[0].Errors()), "first VU generator had errors")
	require.Equal(t, 0, len(vuProfile.Generators[1].Errors()), "second VU generator had errors")
}

func TestBackgroundLoadGoRoutines(t *testing.T) {
	srv := wasp.NewHTTPMockServer(nil)
	srv.Run()

	labels := map[string]string{
		"branch": "background_load_goroutines",
		"commit": "background_load_goroutines",
	}

	zetaSchedule := wasp.Combine(
		wasp.Steps(1, 1, 4, 20*time.Second),
		wasp.Plain(5, 50*time.Second),
		wasp.Steps(5, -1, 4, 20*time.Second))

	channel := make(chan error)

	rpsProfile := wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			T:          t,
			LoadType:   wasp.RPS,
			GenName:    "Zeta",
			Schedule:   zetaSchedule,
			Gun:        NewExampleHTTPGun(srv.URL()),
			Labels:     labels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		}))

	go func() {
		_, err := rpsProfile.Run(true)
		channel <- err
	}()

	// Wait for the first generator enter the steady state
	time.Sleep(20 * time.Second)

	etaSchedule := wasp.Combine(
		wasp.Steps(1, 1, 10, 20*time.Second),
		wasp.Plain(10, 20*time.Second),
		wasp.Steps(10, -1, 10, 10*time.Second))

	iotaSchedule := wasp.Combine(
		wasp.Steps(1, 1, 10, 20*time.Second),
		wasp.Plain(10, 20*time.Second),
		wasp.Steps(10, -1, 10, 10*time.Second))

	vuProfile, err := wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			T:          t,
			LoadType:   wasp.VU,
			GenName:    "Eta",
			Schedule:   etaSchedule,
			VU:         NewExampleScenario(srv.URL()),
			Labels:     labels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).
		Add(wasp.NewGenerator(&wasp.Config{
			T:          t,
			LoadType:   wasp.VU,
			GenName:    "Iota",
			Schedule:   iotaSchedule,
			VU:         NewExampleScenario(srv.URL()),
			Labels:     labels,
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).
		Run(true)
	// check if VU Profile did not return an error (e.g. due to invalid configuration or alerts triggered)
	require.NoError(t, err)

	err = <-channel
	// check if RPS Profile did not return an error (e.g. due to invalid configuration or alerts triggered)
	require.NoError(t, err)

	require.Equal(t, 0, len(rpsProfile.Generators[0].Errors()), "RPS generator had errors errors")
	require.Equal(t, 0, len(vuProfile.Generators[0].Errors()), "first VU generator had errors")
	require.Equal(t, 0, len(vuProfile.Generators[1].Errors()), "second VU generator had errors")
}

func TestParallelLoad(t *testing.T) {
	// start mock http server
	srv := wasp.NewHTTPMockServer(nil)
	srv.Run()

	labels := map[string]string{
		"branch": "parallel_load",
		"commit": "parallel_load",
	}

	// define RPS schedule
	rpsSchedule := wasp.Combine(
		wasp.Steps(1, 1, 9, 10*time.Second), // start with 1 RPS, increment by 1 RPS in 9 steps during 10 seconds
		wasp.Plain(9, 50*time.Second))       // hold 100 RPS for 50 seconds

	// define VU schedule
	vuSchedule := wasp.Combine(
		wasp.Steps(2, 1, 8, 16*time.Second), // start with 2 VUs, increment by 1 VU in 8 steps during 16 seconds
		wasp.Plain(10, 30*time.Second))      // hold 10 VUs for 30 seconds

	// define VU'' schedule
	vu2Schedule := wasp.Combine(
		wasp.Steps(3, 1, 6, 14*time.Second), // start with 3 VUs, increment by 1 VU in 6 steps during 14 seconds
		wasp.Plain(9, 20*time.Second))       // hold 9 VUs for 20 seconds

	rpsGen, err := wasp.NewGenerator(&wasp.Config{
		LoadType:   wasp.RPS,
		Schedule:   rpsSchedule,
		GenName:    "Kappa",
		Labels:     labels,
		Gun:        NewExampleHTTPGun(srv.URL()),
		LokiConfig: wasp.NewEnvLokiConfig(),
	})

	require.NoError(t, err)

	vuGen, err := wasp.NewGenerator(&wasp.Config{
		LoadType:   wasp.VU,
		Schedule:   vuSchedule,
		GenName:    "Lambda",
		Labels:     labels,
		VU:         NewExampleScenario(srv.URL()),
		LokiConfig: wasp.NewEnvLokiConfig(),
	})
	require.NoError(t, err)

	vu2Gen, err := wasp.NewGenerator(&wasp.Config{
		LoadType: wasp.VU,
		Schedule: vu2Schedule,
		GenName:  "Mu",
		Labels:   labels,
		// here both VUs use the same VirtualUser implementation to keep it simple, but in a real-world use case they would different ones
		VU:         NewExampleScenario(srv.URL()),
		LokiConfig: wasp.NewEnvLokiConfig(),
	})
	require.NoError(t, err)

	wg := sync.WaitGroup{}

	// run RPS load in a separate goroutine
	go func() {
		wg.Add(1)
		rpsGen.Run(true)
		wg.Done()
	}()

	// wait for RPS load to reach the steady state
	time.Sleep(10 * time.Second)

	// run VU load in a separate goroutine
	go func() {
		wg.Add(1)
		vuGen.Run(true)
		wg.Done()
	}()

	// wait for VU' load to reach the steady state
	time.Sleep(16 * time.Second)
	vu2Gen.Run(true)

	// wait for RPS and VU' to finish
	wg.Wait()

	// check for load generation errors, although keep in mind that an error during load generation
	// is not necessarily a problem, it might be a result of a load test being too aggressive
	// the correct way to check for errors is to look at the dashboard and alerts
	require.Equal(t, 0, len(rpsGen.Errors()), "RPS generator errors")
	require.Equal(t, 0, len(vuGen.Errors()), "VU generator errors")
	require.Equal(t, 0, len(vu2Gen.Errors()), "VU' generator errors")
}
