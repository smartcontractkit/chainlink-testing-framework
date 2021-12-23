package compatibility

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/utils"
)

var _ = Describe("Json RPC compatibility @json_rpc", func() {
	var (
		rpcClients []*rpc.Client
	)

	BeforeEach(func() {
		By("Getting RPC client", func() {
			nc, err := config.LoadNetworksConfig(filepath.Join(utils.ProjectRoot, "networks.yaml"))
			Expect(err).ShouldNot(HaveOccurred())
			for _, networkName := range nc.SelectedNetworks {
				networkSettings, ok := nc.NetworkSettings[networkName]
				Equal(ok)
				var urls = networkSettings["urls"]
				for _, url := range urls.([]interface{}) {
					rpcClient, err := rpc.Dial(fmt.Sprintf("%v", url))
					Expect(err).ShouldNot(HaveOccurred())
					rpcClients = append(rpcClients, rpcClient)
				}
			}
		})
	})

	Describe("HTTP requests", func() {
		It("methods", func() {
			for _, rpcClient := range rpcClients {
				var block interface{}
				err := rpcClient.CallContext(context.Background(), &block, "eth_getBlockByNumber", "0x5BAD55")
				Expect(err).ShouldNot(HaveOccurred())
			}
		})

		It("schemas", func() {

		})
	})
})
