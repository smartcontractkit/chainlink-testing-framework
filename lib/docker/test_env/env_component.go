package test_env

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logging"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/osutil"
)

const (
	BaseCMD = `docker run -i --rm -v /var/run/docker.sock:/var/run/docker.sock --network %s gaiaadm/pumba --log-level=info`
)

type EnvComponent struct {
	ContainerName      string               `json:"containerName"`
	ContainerImage     string               `json:"containerImage"`
	ContainerVersion   string               `json:"containerVersion"`
	ContainerEnvs      map[string]string    `json:"containerEnvs"`
	WasRecreated       bool                 `json:"wasRecreated"`
	Networks           []string             `json:"networks"`
	Container          tc.Container         `json:"-"`
	LogStream          *logstream.LogStream `json:"-"`
	PostStartsHooks    []tc.ContainerHook   `json:"-"`
	PostStopsHooks     []tc.ContainerHook   `json:"-"`
	PreTerminatesHooks []tc.ContainerHook   `json:"-"`
	LogLevel           string               `json:"-"`
	StartupTimeout     time.Duration        `json:"-"`
	l                  zerolog.Logger
	t                  *testing.T
}

type EnvComponentOption = func(c *EnvComponent)

func WithContainerName(name string) EnvComponentOption {
	return func(c *EnvComponent) {
		if name != "" {
			c.ContainerName = name
		}
	}
}

func WithStartupTimeout(timeout time.Duration) EnvComponentOption {
	return func(c *EnvComponent) {
		if timeout != 0 {
			c.StartupTimeout = timeout
		}
	}
}

func WithContainerImageWithVersion(imageWithVersion string) EnvComponentOption {
	return func(c *EnvComponent) {
		split := strings.Split(imageWithVersion, ":")
		if len(split) == 2 {
			c.ContainerImage = split[0]
			c.ContainerVersion = split[1]
		}
	}
}

func WithLogLevel(logLevel string) EnvComponentOption {
	return func(c *EnvComponent) {
		if logLevel != "" {
			c.LogLevel = logLevel
		}
	}
}

func WithLogStream(ls *logstream.LogStream) EnvComponentOption {
	return func(c *EnvComponent) {
		c.LogStream = ls
	}
}

func WithTestInstance(t *testing.T) EnvComponentOption {
	return func(c *EnvComponent) {
		c.t = t
		c.l = logging.GetTestLogger(t)
	}
}

func WithLogger(l zerolog.Logger) EnvComponentOption {
	return func(c *EnvComponent) {
		c.l = l
	}
}

func WithPostStartsHooks(hooks ...tc.ContainerHook) EnvComponentOption {
	return func(c *EnvComponent) {
		c.PostStartsHooks = hooks
	}
}

func WithPostStopsHooks(hooks ...tc.ContainerHook) EnvComponentOption {
	return func(c *EnvComponent) {
		c.PostStopsHooks = hooks
	}
}

func WithPreTerminatesHooks(hooks ...tc.ContainerHook) EnvComponentOption {
	return func(c *EnvComponent) {
		c.PreTerminatesHooks = hooks
	}
}

func (ec *EnvComponent) SetDefaultHooks() {
	ec.PostStartsHooks = []tc.ContainerHook{
		func(ctx context.Context, c tc.Container) error {
			if ec.LogStream != nil {
				return ec.LogStream.ConnectContainer(ctx, c, "")
			}
			return nil
		},
	}
	ec.PostStopsHooks = []tc.ContainerHook{
		func(ctx context.Context, c tc.Container) error {
			if ec.LogStream != nil {
				return ec.LogStream.DisconnectContainer(c)
			}
			return nil
		},
	}
}

func (ec *EnvComponent) GetImageWithVersion() string {
	return fmt.Sprintf("%s:%s", ec.ContainerImage, ec.ContainerVersion)
}

// ChaosPause pauses the container for the specified duration
func (ec EnvComponent) ChaosPause(
	l zerolog.Logger,
	duration time.Duration,
) error {
	withNet := fmt.Sprintf(BaseCMD, ec.Networks[0])
	return osutil.ExecCmd(l, fmt.Sprintf(`%s pause --duration=%s %s`, withNet, duration.String(), ec.ContainerName))
}

// ChaosNetworkDelay delays the container's network traffic for the specified duration
func (ec EnvComponent) ChaosNetworkDelay(
	l zerolog.Logger,
	duration time.Duration,
	delay time.Duration,
	targetInterfaceName string,
	targetIPs []string,
	targetIngressPorts []string,
	targetEgressPorts []string,
) error {
	var sb strings.Builder
	withNet := fmt.Sprintf(BaseCMD, ec.Networks[0])
	sb.Write([]byte(fmt.Sprintf(`%s netem --tc-image=gaiadocker/iproute2 --duration=%s`, withNet, duration.String())))
	writeTargetNetworkParams(&sb, targetInterfaceName, targetIPs, targetIngressPorts, targetEgressPorts)
	sb.Write([]byte(fmt.Sprintf(` delay --time=%d %s`, delay, ec.ContainerName)))
	return osutil.ExecCmd(l, sb.String())
}

// ChaosNetworkLoss causes the container to lose the specified percentage of network traffic for the specified duration
func (ec EnvComponent) ChaosNetworkLoss(
	l zerolog.Logger,
	duration time.Duration,
	lossPercentage int,
	targetInterfaceName string,
	targetIPs []string,
	targetIngressPorts []string,
	targetEgressPorts []string,
) error {
	var sb strings.Builder
	withNet := fmt.Sprintf(BaseCMD, ec.Networks[0])
	sb.Write([]byte(fmt.Sprintf(`%s netem --tc-image=gaiadocker/iproute2 --duration=%s`, withNet, duration.String())))
	writeTargetNetworkParams(&sb, targetInterfaceName, targetIPs, targetIngressPorts, targetEgressPorts)
	sb.Write([]byte(fmt.Sprintf(` loss --percent %d %s`, lossPercentage, ec.ContainerName)))
	return osutil.ExecCmd(l, sb.String())
}

// writeTargetNetworkParams writes the target network parameters to the provided strings.Builder
func writeTargetNetworkParams(sb *strings.Builder, targetInterfaceName string, targetIPs []string, targetIngressPorts []string, targetEgressPorts []string) {
	if targetInterfaceName == "" {
		targetInterfaceName = "eth0"
	}
	for _, ip := range targetIPs {
		sb.Write([]byte(fmt.Sprintf(` -t %s`, ip)))
	}
	sb.Write([]byte(fmt.Sprintf(" --interface %s", targetInterfaceName)))
	for _, p := range targetIngressPorts {
		sb.Write([]byte(fmt.Sprintf(` --ingress-port %s`, p)))
	}
	for _, p := range targetEgressPorts {
		sb.Write([]byte(fmt.Sprintf(` --egress-port %s`, p)))
	}
}
