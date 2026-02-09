package main

import (
	"fmt"
	"strconv"
	"time"

	p "github.com/smartcontractkit/pods"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/clnode"
)

func main() {
	if err := NewEnvironment(); err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Minute)

	c, err := clclient.New([]*clnode.Output{
		{Node: &clnode.NodeOut{APIAuthUser: clnode.DefaultAPIUser, APIAuthPassword: clnode.DefaultAPIPassword, ExternalURL: "http://localhost:6688"}},
		{Node: &clnode.NodeOut{APIAuthUser: clnode.DefaultAPIUser, APIAuthPassword: clnode.DefaultAPIPassword, ExternalURL: "http://localhost:6689"}},
	})
	if err != nil {
		panic(err)
	}
	for _, n := range c {
		_, resp, err := n.GetForwarders()
		if err != nil {
			panic(err)
		}
		fmt.Println(resp)
	}
}

func NewEnvironment() error {
	ns := "pods"
	chainlinkImg := "public.ecr.aws/chainlink/chainlink:v2.17.0"
	pgImg := "postgres:15"

	cfg := &p.Config{
		Namespace: p.S(ns),
		Pods: []*p.PodConfig{
			{
				Name:    p.S("anvil"),
				Labels:  map[string]string{"chain.link/component": "anvil"},
				Image:   p.S("ghcr.io/foundry-rs/foundry:stable"),
				Ports:   []string{"8545:8545"},
				Command: p.S("anvil --host=0.0.0.0 -b=1 --mixed-mining"),
			},
		},
	}

	k8sDeployment, err := p.New(cfg)
	if err != nil {
		return err
	}
	err = k8sDeployment.Apply()
	if err != nil {
		return err
	}
	k8sDeployment.ResetPodsConfig()

	for i := 0; i < 2; i++ {
		cfg.Pods = append(cfg.Pods, &p.PodConfig{
			Name:                     p.S(fmt.Sprintf("node-%d", i)),
			Labels:                   map[string]string{"chain.link/component": "cl", "instance": strconv.Itoa(i)},
			ContainerSecurityContext: p.CLUserContainerSecurityCtx(),
			Image:                    p.S(chainlinkImg),
			Ports:                    []string{"6688:6688", "6690:6690"},
			ConfigMap: map[string]*string{
				"config.toml": p.S(`
	[WebServer]
	AllowOrigins = "*"
	HTTPWriteTimeout = "3m0s"
	HTTPPort = 6688
	SecureCookies = false
	SessionTimeout = "999h0m0s"
	[WebServer.TLS]
	HTTPSPort = 0
	[Feature]
	LogPoller = true
	[OCR2]
	Enabled = true
	DatabaseTimeout = '1s'
	[P2P.V2]
	Enabled = true
	ListenAddresses = ['0.0.0.0:5001']
	[[EVM]]
	ChainID = "31337"
	AutoCreateKey = true
	FinalityDepth = 1
	MinContractPayment = "0"
	[[EVM.Nodes]]
	Name = "Anvil"
	WSURL = "ws://anvil-svc:8545"
	HTTPURL = "http://anvil-svc:8545"
	`),
				"apicredentials": p.S(fmt.Sprintf(`%s
	%s`, clnode.DefaultAPIUser, clnode.DefaultAPIPassword)),
			},
			ConfigMapMountPath: map[string]*string{
				"config.toml":    p.S("/config.toml"),
				"apicredentials": p.S("/apicredentials"),
			},
			Secrets: map[string]*string{"secrets.toml": p.S(fmt.Sprintf(`
	[Database]
	URL = 'postgresql://chainlink:thispasswordislongenough@postgres-%d-svc:5432/chainlink?sslmode=disable'
	[Password]
	Keystore = 'thispasswordislongenough'
	`, i))},
			SecretsMountPath: map[string]*string{"secrets.toml": p.S("/secrets.toml")},
			Command:          p.S("chainlink -c /config.toml -s /secrets.toml node start -d -a /apicredentials"),
		})
		cfg.Pods = append(cfg.Pods, p.PostgreSQL(fmt.Sprintf("postgres-%d", i), pgImg, p.ResourcesSmall(), p.ResourcesSmall(), p.S("10Gi")))
	}

	err = k8sDeployment.Apply()
	if err != nil {
		return err
	}

	return p.NewForwarder(k8sDeployment.API).
		Forward([]p.PortForwardConfig{
			{
				ServiceName:   "anvil-svc",
				LocalPort:     8545,
				ContainerPort: 8545,
				Namespace:     ns,
			},
			{
				ServiceName:   "node-0-svc",
				LocalPort:     6688,
				ContainerPort: 6688,
				Namespace:     ns,
			},
			{
				ServiceName:   "node-1-svc",
				LocalPort:     6689,
				ContainerPort: 6688,
				Namespace:     ns,
			},
		},
		)
}
