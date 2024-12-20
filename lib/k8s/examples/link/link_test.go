package link

import (
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/blockchain"
	ctf_config "github.com/smartcontractkit/chainlink-testing-framework/lib/config"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/ptr"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink-testing-framework/seth/contracts/bind/link_token_interface"
)

func TestSimpleDONWithLinkContract(t *testing.T) {
	tomlConfig := `[Feature]
FeedsManager = true
LogPoller = true
UICSAKeys = true

[Database]
MaxIdleConns = 20
MaxOpenConns = 40
MigrateOnStartup = true

[Log]
Level = "debug"
JSONConsole = true

[Log.File]
MaxSize = "0b"

[WebServer]
AllowOrigins = "*"
HTTPWriteTimeout = "3m0s"
HTTPPort = 6688
SecureCookies = false
SessionTimeout = "999h0m0s"

[WebServer.RateLimit]
Authenticated = 2000
Unauthenticated = 1000

[WebServer.TLS]
HTTPSPort = 0

[OCR]
Enabled = true

[P2P]

[P2P.V2]
ListenAddresses = ["0.0.0.0:6690"]

[[EVM]]
ChainID = "1337"
AutoCreateKey = true
FinalityDepth = 1
MinContractPayment = "0"

[EVM.GasEstimator]
PriceMax = "200 gwei"
LimitDefault = 6000000
FeeCapDefault = "200 gwei"

[[EVM.Nodes]]
Name = "Simulated Geth-0"
WSURL = "ws://geth:8546"
HTTPURL = "http://geth:8544"`

	chainlinkImageCfg := &ctf_config.ChainlinkImageConfig{
		Image:   ptr.Ptr("public.ecr.aws/chainlink/chainlink"),
		Version: ptr.Ptr("2.19.0"),
	}

	var overrideFn = func(_ interface{}, target interface{}) {
		ctf_config.MustConfigOverrideChainlinkVersion(chainlinkImageCfg, target)
	}

	cd := chainlink.NewWithOverride(0, map[string]any{
		"replicas": 6,
		"toml":     tomlConfig,
		"db": map[string]any{
			"stateful": true, // stateful DB by default for soak tests
		},
		"prometheus": true,
	}, chainlinkImageCfg, overrideFn)

	productName := "data-feedsv1.0"
	nsLabels, err := environment.GetRequiredChainLinkNamespaceLabels(productName, "soak")
	if err != nil {
		t.Fatal("Error creating required chain.link labels for namespace", err)
	}

	workloadPodLabels, err := environment.GetRequiredChainLinkWorkloadAndPodLabels(productName, "soak")
	if err != nil {
		t.Fatal("Error creating required chain.link labels for workload and pod", err)
	}

	baseEnvironmentConfig := &environment.Config{
		TTL:                time.Hour * 720, // 30 days,
		NamespacePrefix:    "bartek-ocr",
		Test:               t,
		PreventPodEviction: true,
		Labels:             nsLabels,
		WorkloadLabels:     workloadPodLabels,
		PodLabels:          workloadPodLabels,
	}

	nodeNetwork := blockchain.SimulatedEVMNetwork

	ethProps := &ethereum.Props{
		NetworkName: nodeNetwork.Name,
		Simulated:   nodeNetwork.Simulated,
		WsURLs:      nodeNetwork.URLs,
		HttpURLs:    nodeNetwork.HTTPURLs,
	}

	testEnv := environment.New(baseEnvironmentConfig).
		AddHelm(ethereum.New(ethProps)).
		AddHelm(cd)

	err = testEnv.Run()
	if err != nil {
		t.Fatal("Error running environment: ", err)
	}

	log.Info().Bool("Remote runner?", testEnv.WillUseRemoteRunner()).Msg("Started environment")

	if testEnv.WillUseRemoteRunner() {
		log.Info().Msg("Exiting as test will use remote runner")
		return
	}

	// if test is running inside K8s, nothing to do, default network urls are correct
	if !testEnv.Cfg.InsideK8s {
		// Test is running locally, set forwarded URL of Geth blockchain node
		wsURLs := testEnv.URLs[blockchain.SimulatedEVMNetwork.Name]
		httpURLs := testEnv.URLs[blockchain.SimulatedEVMNetwork.Name+"_http"]
		if len(wsURLs) == 0 || len(httpURLs) == 0 {
			t.Fatal("Forwarded Geth URLs should not be empty")
		}
		nodeNetwork.URLs = wsURLs
		nodeNetwork.HTTPURLs = httpURLs
	}

	sethClient, err := seth.NewClientBuilder().
		WithRpcUrl(nodeNetwork.URLs[0]).
		WithPrivateKeys([]string{nodeNetwork.PrivateKeys[0]}).
		Build()
	if err != nil {
		t.Fatal("Error creating Seth client", err)
	}

	for i := 0; i < 5; i++ {
		log.Info().
			Msgf("Deploying LinkToken contract, instance %d/%d", i+1, 5)

		linkTokenAbi, err := link_token_interface.LinkTokenMetaData.GetAbi()
		if err != nil {
			t.Fatal("Error getting LinkToken ABI", err)
		}
		linkDeploymentData, err := sethClient.DeployContract(sethClient.NewTXOpts(), "LinkToken", *linkTokenAbi, common.FromHex(link_token_interface.LinkTokenMetaData.Bin))
		if err != nil {
			t.Fatal("Error deploying LinkToken contract", err)
		}
		linkToken, err := link_token_interface.NewLinkToken(linkDeploymentData.Address, sethClient.Client)
		if err != nil {
			t.Fatal("Error creating LinkToken contract instance", err)
		}

		totalSupply, err := linkToken.TotalSupply(sethClient.NewCallOpts())
		if err != nil {
			t.Fatal("Error getting total supply of LinkToken", err)
		}

		if totalSupply.Cmp(big.NewInt(0)) <= 0 {
			t.Fatal("Total supply of LinkToken should be greater than 0")
		}

		time.Sleep(15 * time.Second)
	}

	// here you could proceed with your test logic
	// for example, by deploying other contracts, funding accounts, etc.
	// and maybe generating some load on the system using WASP?
}
