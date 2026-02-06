package pods_test

import (
	"context"
	"os"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	p "github.com/smartcontractkit/pods"
	"github.com/stretchr/testify/require"
)

func Apply() bool { return os.Getenv("APPLY") == "true" }

func defaultNoErr(t *testing.T, _ *p.Pods, err error) { require.NoError(t, err) }

func onePod(t *testing.T, p *p.Pods, err error) {
	require.NoError(t, err)
	pods, err := p.GetPods(context.Background())
	require.NoError(t, err)
	require.Len(t, pods.Items, 1)
}

func TestPods(t *testing.T) {
	tests := []struct {
		name               string
		props              *p.Config
		skipCI             bool
		validateManifest   func(t *testing.T, p *p.Pods, err error)
		validateDeployment func(t *testing.T, p *p.Pods, err error)
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
			validateManifest:   defaultNoErr,
			validateDeployment: onePod,
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
			validateManifest:   defaultNoErr,
			validateDeployment: onePod,
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
			validateDeployment: func(t *testing.T, p *p.Pods, err error) {
				require.NoError(t, err)
				pods, err := p.GetPods(context.Background())
				require.NoError(t, err)
				require.Len(t, pods.Items, 2)
			},
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
			validateDeployment: func(t *testing.T, p *p.Pods, err error) {
				require.NoError(t, err)
				pods, err := p.GetPods(context.Background())
				require.NoError(t, err)
				require.Len(t, pods.Items, 2)
			},
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
			validateManifest:   defaultNoErr,
			validateDeployment: onePod,
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
			validateManifest: func(t *testing.T, p *p.Pods, err error) {
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
						ConfigMap: map[string]*string{
							"config.toml":  p.S(`test`),
							"config2.toml": p.S(`test`),
						},
						ConfigMapMountPath: map[string]*string{
							"config.toml":  p.S("/config.toml"),
							"config2.toml": p.S("/config2.toml"),
						},
					},
				},
			},
			validateManifest: defaultNoErr,
			validateDeployment: func(t *testing.T, p *p.Pods, err error) {
				onePod(t, p, err)
				cms, err := p.GetConfigMaps(context.Background())
				require.NoError(t, err)
				require.Equal(t, 2, len(cms))
				require.Equal(t, `test`, cms["test-pod-1-configmap"]["config.toml"])
				require.Equal(t, `test`, cms["test-pod-1-configmap"]["config2.toml"])
			},
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
						Secrets: map[string]*string{
							"secret.toml":  p.S(`test`),
							"secret2.toml": p.S(`test`),
						},
						SecretsMountPath: map[string]*string{
							"secret.toml":  p.S("/secret.toml"),
							"secret2.toml": p.S("/secret2.toml"),
						},
					},
				},
			},
			validateManifest: defaultNoErr,
			validateDeployment: func(t *testing.T, p *p.Pods, err error) {
				onePod(t, p, err)
				secrets, err := p.GetSecrets(context.Background())
				require.NoError(t, err)
				require.Equal(t, 1, len(secrets))
				require.Equal(t, []byte(`test`), secrets["test-pod-1-secret"]["secret.toml"])
				require.Equal(t, []byte(`test`), secrets["test-pod-1-secret"]["secret2.toml"])
			},
		},
		{
			name:   "test-volumes",
			skipCI: true,
			props: &p.Config{
				Namespace: p.S("test-volumes"),
				Pods:      []*p.PodConfig{p.PostgreSQL("pg-x", "postgres:15", p.ResourcesSmall(), p.ResourcesSmall(), p.S("1Gi"))},
			},
			validateManifest: defaultNoErr,
			validateDeployment: func(t *testing.T, p *p.Pods, err error) {
				onePod(t, p, err)
				volumes, err := p.GetPersistentVolumes(context.Background())
				require.NoError(t, err)
				require.Len(t, volumes, 1)
			},
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
			validateDeployment: func(t *testing.T, p *p.Pods, err error) {
				require.NoError(t, err)
				svcs, err := p.GetServices(context.Background())
				require.NoError(t, err)
				require.Len(t, svcs, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if os.Getenv("CI") == "true" && tt.skipCI {
				t.Skip("this test can't be run in CI because of GHA limitations")
			}
			p := p.New(tt.props)
			err := p.Generate()
			if tt.validateManifest != nil {
				tt.validateManifest(t, p, err)
			}
			if Apply() {
				err := p.CreateNamespace(*tt.props.Namespace)
				require.NoError(t, err)
				t.Cleanup(func() {
					_ = p.RemoveNamespace(*tt.props.Namespace)
				})
				err = p.Apply()
				if tt.validateDeployment != nil {
					tt.validateDeployment(t, p, err)
				}
				return
			}
			if err == nil {
				snaps.MatchSnapshot(t, *p.Manifest())
			}
		})
	}
}
