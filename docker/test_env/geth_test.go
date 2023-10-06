package test_env

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/docker"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"
)

const count = 500
const count2 = 1

func TestGeth2(t *testing.T) {
	t.Parallel()
	testHelper(t)
}

func testHelper(t *testing.T) {
	l := logging.GetTestLogger(t)
	network, err := docker.CreateNetwork(l)
	require.NoError(t, err)
	for i := 0; i < count; i++ {
		g := NewGeth([]string{network.Name}).
			WithTestLogger(t)
		_, _, err = g.StartContainer()
		require.NoError(t, err)
		fmt.Printf("Finished %d\n", i)

		time.Sleep(5 * time.Second)
		ns := blockchain.SimulatedEVMNetwork
		ns.URLs = []string{g.ExternalWsUrl}
		_, evmError := blockchain.ConnectEVMClient(ns, l)

		if evmError != nil {
			d, err1 := g.Container.Logs(context.Background())
			require.NoError(t, err1)
			defer d.Close()
			buf := new(bytes.Buffer)
			_, err1 = io.Copy(buf, d)
			require.NoError(t, err1)
			fmt.Println(buf.String())
			out, err := exec.Command("docker", "inspect", g.ContainerName).Output()
			require.NoError(t, err)
			fmt.Println(string(out))
			require.NoError(t, evmError, "Couldn't connect to the evm client")
		}

		err = g.Container.Terminate(context.Background())
		require.NoError(t, err)
	}
}

func TestGeth1(t *testing.T) {
	t.Parallel()
	testHelper2(t)
}

func testHelper2(t *testing.T) {
	l := logging.GetTestLogger(t)
	network, err := docker.CreateNetwork(l)
	require.NoError(t, err)
	for i := 0; i < count2; i++ {
		g := NewGeth([]string{network.Name}).
			WithTestLogger(t)
		_, _, err = g.StartContainer()
		require.NoError(t, err)
		fmt.Printf("Finished %d\n", i)

		time.Sleep(5 * time.Second)
		ns := blockchain.SimulatedEVMNetwork
		ns.URLs = []string{g.ExternalWsUrl}
		_, evmError := blockchain.ConnectEVMClient(ns, l)

		// if evmError != nil {
		d, err1 := g.Container.Logs(context.Background())
		require.NoError(t, err1)
		defer d.Close()
		buf := new(bytes.Buffer)
		_, err1 = io.Copy(buf, d)
		require.NoError(t, err1)
		fmt.Println(buf.String())
		out, err := exec.Command("docker", "inspect", g.ContainerName).Output()
		require.NoError(t, err)
		fmt.Println(string(out))
		require.NoError(t, evmError, "Couldn't connect to the evm client")
		// }

		err = g.Container.Terminate(context.Background())
		require.NoError(t, err)
	}
}
