package test_env

import (
	"fmt"
	"time"

	"strings"

	"github.com/rs/zerolog"
	tc "github.com/testcontainers/testcontainers-go"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/osutil"
)

const (
	BaseCMD = `docker run -i --rm -v /var/run/docker.sock:/var/run/docker.sock --network %s gaiaadm/pumba --log-level=info`
)

type EnvComponent struct {
	ContainerName      string             `json:"containerName"`
	ContainerImage     string             `json:"containerImage"`
	ContainerVersion   string             `json:"containerVersion"`
	ContainerEnvs      map[string]string  `json:"containerEnvs"`
	WasRecreated       bool               `json:"wasRecreated"`
	Networks           []string           `json:"networks"`
	Container          tc.Container       `json:"-"`
	PostStartsHooks    []tc.ContainerHook `json:"-"`
	PostStopsHooks     []tc.ContainerHook `json:"-"`
	PreTerminatesHooks []tc.ContainerHook `json:"-"`
	LogLevel           string             `json:"-"`
	StartupTimeout     time.Duration      `json:"-"`
}

type EnvComponentOption = func(c *EnvComponent)

// WithContainerName sets the container name for an EnvComponent.
// It allows customization of the container's identity, enhancing clarity
// and organization in containerized environments.
func WithContainerName(name string) EnvComponentOption {
	return func(c *EnvComponent) {
		if name != "" {
			c.ContainerName = name
		}
	}
}

// WithStartupTimeout sets a custom startup timeout for an EnvComponent.
// This option allows users to specify how long to wait for the component to start
// before timing out, enhancing control over component initialization.
func WithStartupTimeout(timeout time.Duration) EnvComponentOption {
	return func(c *EnvComponent) {
		if timeout != 0 {
			c.StartupTimeout = timeout
		}
	}
}

// WithContainerImageWithVersion sets the container image and version for an EnvComponent.
// It splits the provided image string by ':' and assigns the values accordingly.
// This function is useful for configuring specific container images in a deployment.
func WithContainerImageWithVersion(imageWithVersion string) EnvComponentOption {
	return func(c *EnvComponent) {
		split := strings.Split(imageWithVersion, ":")
		if len(split) == 2 {
			c.ContainerImage = split[0]
			c.ContainerVersion = split[1]
		}
	}
}

// WithLogLevel sets the logging level for an environment component.
// It allows customization of log verbosity, enhancing debugging and monitoring capabilities.
func WithLogLevel(logLevel string) EnvComponentOption {
	return func(c *EnvComponent) {
		if logLevel != "" {
			c.LogLevel = logLevel
		}
	}
}

// WithPostStartsHooks sets the PostStarts hooks for an EnvComponent.
// This allows users to define custom actions that should occur after the component starts.
func WithPostStartsHooks(hooks ...tc.ContainerHook) EnvComponentOption {
	return func(c *EnvComponent) {
		c.PostStartsHooks = hooks
	}
}

// WithPostStopsHooks sets the PostStops hooks for an EnvComponent.
// This allows users to define custom actions that should occur after the component stops.
func WithPostStopsHooks(hooks ...tc.ContainerHook) EnvComponentOption {
	return func(c *EnvComponent) {
		c.PostStopsHooks = hooks
	}
}

// WithPreTerminatesHooks sets the pre-termination hooks for an EnvComponent.
// This allows users to define custom behavior that should occur before the component is terminated.
func WithPreTerminatesHooks(hooks ...tc.ContainerHook) EnvComponentOption {
	return func(c *EnvComponent) {
		c.PreTerminatesHooks = hooks
	}
}

// SetDefaultHooks initializes the default hooks for the environment component.
// This function is useful for ensuring that the component has a consistent starting state before further configuration.
func (ec *EnvComponent) SetDefaultHooks() {
	// no default hooks
}

// GetImageWithVersion returns the container image name combined with its version.
// This function is useful for generating a complete image identifier needed for container requests.
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
