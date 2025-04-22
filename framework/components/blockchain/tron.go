package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Accounts struct {
	HDPath      string   `json:"hdPath"`
	Mnemonic    string   `json:"mnemonic"`
	PrivateKeys []string `json:"privateKeys"`
	More        []string `json:"more"`
}

var TRONAccounts = Accounts{
	HDPath:   "m/44'/195'/0'/0/",
	Mnemonic: "resemble birth wool happy sun burger fatal trumpet globe purity health ritual",
	PrivateKeys: []string{
		"932a39242805a1b1095638027f26af9664d1d5bf8ab3b7527ee75e7efb2946dd",
		"1c17c9c049d36cde7e5ea99df6c86e0474b04f0e258ab619a1e674f397a17152",
		"458130a239671674746582184711a6f8d633355df1a491b9f3b323576134c2e9",
		"2676fd1427968e07feaa9aff967d4ba7607c5497c499968c098d0517cd75cfbb",
		"d26b24a691ff2b03ee6ab65bf164def216f73574996b9ca6299c43a9a63767ac",
		"55df6adf3d081944dbe4688205d94f236fb4427ac44f3a286a96d47db0860667",
		"8a9a60ddd722a40753c2a38edd6b6fa38e806d681c9b08a520ba4912e62b6458",
		"75eb182fb623acf5e53d9885c4e8578f2530533a96c753481cc4277ecc6022de",
		"6c4b22b1d9d68ef7a8ecd151cd4ffdd4ecc2a7b3a3f8a9f9f9bbdbcef6671f10",
		"e578d66453cb41b6c923b9caa91c375a0545eeb171ccafc60b46fa834ce5c200",
	},
	// should not be empty, otherwise TRE will panic
	More: []string{},
}

const (
	DefaultTronPort = "9090"
)

func defaultTron(in *Input) {
	if in.Image == "" {
		in.Image = "tronbox/tre"
	}
	if in.Port == "" {
		in.Port = DefaultTronPort
	}
}

func newTron(in *Input) (*Output, error) {
	defaultTron(in)
	ctx := context.Background()

	containerName := framework.DefaultTCName("blockchain-node")
	bindPort := fmt.Sprintf("%s/tcp", in.Port)

	accounts, err := os.CreateTemp("", "accounts.json")
	if err != nil {
		return nil, err
	}
	accountsData, err := json.Marshal(TRONAccounts)
	if err != nil {
		return nil, err
	}

	_, err = accounts.WriteString(string(accountsData))
	if err != nil {
		return nil, err
	}

	req := testcontainers.ContainerRequest{
		AlwaysPullImage: in.PullImage,
		Image:           in.Image,
		Name:            containerName,
		ExposedPorts:    []string{bindPort},
		Networks:        []string{framework.DefaultNetworkName},
		NetworkAliases: map[string][]string{
			framework.DefaultNetworkName: {containerName},
		},
		Labels: framework.DefaultTCLabels(),
		HostConfigModifier: func(h *container.HostConfig) {
			h.PortBindings = framework.MapTheSamePort(bindPort)
			framework.ResourceLimitsFunc(h, in.ContainerResources)
		},
		WaitingFor: wait.ForLog("Mnemonic").WithPollInterval(200 * time.Millisecond).WithStartupTimeout(1 * time.Minute),
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      accounts.Name(),
				ContainerFilePath: "/config/accounts.json",
				FileMode:          0644,
			},
		},
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := c.Host(ctx)
	if err != nil {
		return nil, err
	}

	return &Output{
		UseCache:      true,
		ChainID:       in.ChainID,
		Type:          in.Type,
		Family:        FamilyTron,
		ContainerName: containerName,
		Nodes: []*Node{
			{
				ExternalHTTPUrl: fmt.Sprintf("http://%s:%s", host, in.Port),
				InternalHTTPUrl: fmt.Sprintf("http://%s:%s", containerName, in.Port),
			},
		},
	}, nil
}
