package client_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/hooks"
	"github.com/smartcontractkit/integrations-framework/utils"
)

const (
	fetchConfig string = "%s/config/test_configs/fetch_config"
)

var _ = Describe("Blockchain @unit", func() {
	var conf *config.Config

	BeforeEach(func() {
		var err error
		conf, err = config.NewConfig(fmt.Sprintf(fetchConfig, utils.ProjectRoot))
		Expect(err).ShouldNot(HaveOccurred())
	})

	DescribeTable("create new wallet configurations", func(
		initFunc hooks.NewNetworkHook,
	) {
		networkConfig, err := initFunc(conf)
		Expect(err).ShouldNot(HaveOccurred())

		networkConfig.Config().PrivateKeyStore, err = environment.NewPrivateKeyStoreFromEnv(&environment.K8sEnvironment{}, networkConfig.Config())
		Expect(err).ShouldNot(HaveOccurred())

		wallets, err := networkConfig.Wallets()
		Expect(err).ShouldNot(HaveOccurred())
		rawWallets := wallets.All()
		for index := range rawWallets {
			_, err := wallets.Wallet(index)
			Expect(err).ShouldNot(HaveOccurred())
		}
	},
		Entry("on Ethereum", client.DefaultNetworkFromConfig),
	)
})
