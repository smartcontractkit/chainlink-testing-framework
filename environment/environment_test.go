package environment_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/environment"
	"github.com/smartcontractkit/integrations-framework/tools"
)

// Test config files
const (
	specifiedConfig string = "%s/config/test_configs/specified_config"
	noPrivateKeysConfig string = "%s/config/test_configs/no_private_keys_config"
)

var _ = Describe("Environment unit tests @unit", func() {

	Describe("NewPrivateKeyStoreFromEnv unit tests", func() {
		FIt("should fetch private keys when they exist", func() {
			conf, err := config.NewConfig(fmt.Sprintf(specifiedConfig, tools.ProjectRoot))
			Expect(err).ShouldNot(HaveOccurred())

			bcNetwork, err := client.NewNetworkFromConfig(conf)
			Expect(err).ShouldNot(HaveOccurred())

			bcNetwork.Config().PrivateKeyStore, err = environment.NewPrivateKeyStoreFromEnv(environment.K8sEnvironment{}, bcNetwork.Config())
			Expect(err).ShouldNot(HaveOccurred())

			privateKeys, err := bcNetwork.Config().PrivateKeyStore.Fetch()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(len(privateKeys)).Should(Equal(1), "The number of private keys was incorrect")
			Expect(privateKeys[0]).Should(Equal("5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"), "The private key did not get read correctly")
		})

		FIt("should not fetch private keys when they do not exist", func() {
			conf, err := config.NewConfig(fmt.Sprintf(noPrivateKeysConfig, tools.ProjectRoot))
			Expect(err).ShouldNot(HaveOccurred())

			bcNetwork, err := client.NewNetworkFromConfig(conf)
			Expect(err).ShouldNot(HaveOccurred())

			bcNetwork.Config().PrivateKeyStore, err = environment.NewPrivateKeyStoreFromEnv(environment.K8sEnvironment{}, bcNetwork.Config())
			Expect(err).ShouldNot(HaveOccurred())

			_, err = bcNetwork.Config().PrivateKeyStore.Fetch()
			Expect(err.Error()).Should(ContainSubstring("no keys found"))
		})
	})
})