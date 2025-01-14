package jd_test

import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/stretchr/testify/require"
	"os"
	"sync"
	"testing"
)

// here we only test that we can boot up JD
// client examples are under "examples" dir
// since JD is private this env var should be set locally and in CI
// TODO: add ComponentDocker prefix to turn this on when we'll have access to ECRs
func TestJD(t *testing.T) {
	err := framework.DefaultNetwork(t, &sync.Once{})
	require.NoError(t, err)
	_, err = jd.NewJD(&jd.Input{
		Image: os.Getenv("CTF_JD_IMAGE"),
	})
	require.NoError(t, err)
}
