package jd_test

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/postgres"
	"github.com/stretchr/testify/require"
	"os"
	"sync"
	"testing"
)

// here we only test that we can boot up JD
// client examples are under "examples" dir
// since JD is private this env var should be set locally and in CI
func TestJD(t *testing.T) {
	err := framework.DefaultNetwork(&sync.Once{})
	require.NoError(t, err)
	pgOut, err := postgres.NewPostgreSQL(&postgres.Input{
		Image: "postgres:12.0",
	})
	require.NoError(t, err)
	out, err := jd.NewJD(&jd.Input{
		DBURL: pgOut.JDDockerInternalURL,
		Image: os.Getenv("CTF_JD_IMAGE"),
	})
	require.NoError(t, err)
	spew.Dump(out)
}
