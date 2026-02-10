package pods_test

import (
	"os"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	p "github.com/smartcontractkit/chainlink-testing-framework/pods"
	"github.com/stretchr/testify/require"
)

func defaultNoErr(t *testing.T, err error) { require.NoError(t, err) }

func TestPods(t *testing.T) {
	tests := []struct {
		name             string
		props            *p.Config
		skipCI           bool
		validateManifest func(t *testing.T, err error)
	}{
		{
			name: "test-single-pod",
			props: &p.Config{
				Namespace: p.S("test-single-pod"),
				Pods: []*p.PodConfig{
					{
						Name:  p.S("test-pod-1"),
						Image: p.S("nginx:latest"),
						Ports: []string{"80:80"},
					},
				},
			},
			validateManifest: defaultNoErr,
		},
		{
			name: "test-command",
			props: &p.Config{
				Namespace: p.S("test-command"),
				Pods: []*p.PodConfig{
					{
						Name:        p.S("anvil"),
						Labels:      map[string]string{"chain.link/component": "cl"},
						Annotations: map[string]string{"custom-annotation": "custom"},
						Image:       p.S("ghcr.io/foundry-rs/foundry"),
						Ports:       []string{"8545:8545"},
						Command:     p.S("anvil --host=0.0.0.0 -b=1"),
					},
				},
			},
			validateManifest: defaultNoErr,
		},
		{
			name: "test-instances",
			props: &p.Config{
				Namespace: p.S("test-instances"),
				Pods: []*p.PodConfig{
					{
						Name:     p.S("test-pod-1"),
						Image:    p.S("nginx:latest"),
						Ports:    []string{"80:80"},
						Replicas: p.I(2),
					},
				},
			},
			validateManifest: defaultNoErr,
		},
		{
			name: "test-multiple-pods",
			props: &p.Config{
				Namespace: p.S("test-multiple-pods"),
				Pods: []*p.PodConfig{
					{
						Name:  p.S("test-pod-1"),
						Image: p.S("nginx:latest"),
						Ports: []string{"80:80"},
					},
					{
						Name:  p.S("test-pod-2"),
						Image: p.S("nginx:latest"),
						Ports: []string{"80:80"},
					},
				},
			},
			validateManifest: defaultNoErr,
		},
		{
			name: "test-custom-resources",
			props: &p.Config{
				Namespace: p.S("test-custom-resources"),
				Pods: []*p.PodConfig{
					{
						Name:     p.S("test-pod-1"),
						Image:    p.S("nginx:latest"),
						Ports:    []string{"80:80"},
						Requests: p.Resources("250m", "1Gi"),
						Limits:   p.Resources("500m", "2Gi"),
					},
				},
			},
			validateManifest: defaultNoErr,
		},
		{
			name: "test-invalid-ports",
			props: &p.Config{
				Namespace: p.S("test-invalid-ports"),
				Pods: []*p.PodConfig{
					{
						Name:  p.S("test-pod-1"),
						Image: p.S("nginx:latest"),
						Ports: []string{"80"},
					},
				},
			},
			validateManifest: func(t *testing.T, err error) {
				require.Contains(t, err.Error(), "invalid port mapping")
			},
		},
		{
			name: "test-configmaps",
			props: &p.Config{
				Namespace: p.S("test-configmaps"),
				Pods: []*p.PodConfig{
					{
						Name:  p.S("test-pod-1"),
						Image: p.S("nginx:latest"),
						Ports: []string{"80:80"},
						ConfigMap: map[string]string{
							"config.toml":  `test`,
							"config2.toml": `test`,
						},
						ConfigMapMountPath: map[string]string{
							"config.toml":  "/config.toml",
							"config2.toml": "/config2.toml",
						},
					},
				},
			},
			validateManifest: defaultNoErr,
		},
		{
			name: "test-secrets",
			props: &p.Config{
				Namespace: p.S("test-secrets"),
				Pods: []*p.PodConfig{
					{
						Name:  p.S("test-pod-1"),
						Image: p.S("nginx:latest"),
						Ports: []string{"80:80"},
						Secrets: map[string]string{
							"secret.toml":  `test`,
							"secret2.toml": `test`,
						},
						SecretsMountPath: map[string]string{
							"secret.toml":  "/secret.toml",
							"secret2.toml": "/secret2.toml",
						},
					},
				},
			},
			validateManifest: defaultNoErr,
		},
		{
			name: "test-services",
			props: &p.Config{
				Namespace: p.S("test-services"),
				Pods: []*p.PodConfig{
					{
						Name:  p.S("test-pod-1"),
						Image: p.S("nginx:latest"),
						Ports: []string{"80:80"},
					},
				},
			},
			validateManifest: defaultNoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if os.Getenv("CI") == "true" && tt.skipCI {
				t.Skip("this test can't be run in CI because of GHA limitations")
			}
			manifest, err := p.Run(tt.props)
			tt.validateManifest(t, err)
			if err == nil {
				snaps.MatchSnapshot(t, manifest)
			}
		})
	}
}
