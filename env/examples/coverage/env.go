package main

import (
	"time"

	"github.com/smartcontractkit/chainlink-env/environment"
	goc "github.com/smartcontractkit/chainlink-env/pkg/cdk8s/goc"
	dummy "github.com/smartcontractkit/chainlink-env/pkg/cdk8s/http_dummy"
)

func main() {
	e := environment.New(nil).
		AddChart(goc.New()).
		AddChart(dummy.New())
	if err := e.Run(); err != nil {
		panic(err)
	}
	// run your test logic here
	time.Sleep(1 * time.Minute)
	if err := e.SaveCoverage(); err != nil {
		panic(err)
	}
	// clear the coverage, rerun the tests again if needed
	if err := e.ClearCoverage(); err != nil {
		panic(err)
	}
}
