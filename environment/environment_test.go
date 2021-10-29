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
	specifiedConfig     string = "%s/config/test_configs/specified_config"
	noPrivateKeysConfig string = "%s/config/test_configs/no_private_keys_config"
	secretKeysConfig    string = "%s/config/test_configs/secret_keys_config"
)

var _ = Describe("Environment unit tests @unit", func() {

	Describe("NewPrivateKeyStoreFromEnv unit tests", func() {
		It("should fetch private keys when they exist", func() {
			conf, err := config.NewConfig(fmt.Sprintf(specifiedConfig, tools.ProjectRoot))
			Expect(err).ShouldNot(HaveOccurred())

			bcNetwork, err := client.DefaultNetworkFromConfig(conf)
			Expect(err).ShouldNot(HaveOccurred())

			bcNetwork.Config().PrivateKeyStore, err = environment.NewPrivateKeyStoreFromEnv(&environment.K8sEnvironment{}, bcNetwork.Config())
			Expect(err).ShouldNot(HaveOccurred())

			privateKeys, err := bcNetwork.Config().PrivateKeyStore.Fetch()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(len(privateKeys)).Should(Equal(1), "The number of private keys was incorrect")
			Expect(privateKeys[0]).Should(Equal("5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"), "The private key did not get read correctly")
		})

		It("should not fetch private keys when they do not exist", func() {
			conf, err := config.NewConfig(fmt.Sprintf(noPrivateKeysConfig, tools.ProjectRoot))
			Expect(err).ShouldNot(HaveOccurred())

			bcNetwork, err := client.DefaultNetworkFromConfig(conf)
			Expect(err).ShouldNot(HaveOccurred())

			bcNetwork.Config().PrivateKeyStore, err = environment.NewPrivateKeyStoreFromEnv(&environment.K8sEnvironment{}, bcNetwork.Config())
			Expect(err).ShouldNot(HaveOccurred())

			_, err = bcNetwork.Config().PrivateKeyStore.Fetch()
			Expect(err.Error()).Should(ContainSubstring("no keys found"))
		})

		It("should fetch secret private keys", func() {
			Skip("Not ready to be run in github")
			conf, err := config.NewConfig(fmt.Sprintf(secretKeysConfig, tools.ProjectRoot))
			Expect(err).ShouldNot(HaveOccurred())

			bcNetwork, err := client.DefaultNetworkFromConfig(conf)
			Expect(err).ShouldNot(HaveOccurred())

			env, err := environment.NewK8sEnvironment(conf, bcNetwork)
			Expect(err).ShouldNot(HaveOccurred())

			err = env.DeploySpecs(environment.NewChainlinkCluster(1))
			Expect(err).ShouldNot(HaveOccurred())

			bcNetwork.Config().PrivateKeyStore, err = environment.NewPrivateKeyStoreFromEnv(env, bcNetwork.Config())
			Expect(err).ShouldNot(HaveOccurred())

			privateKeys, err := bcNetwork.Config().PrivateKeyStore.Fetch()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(len(privateKeys)).Should(Equal(2), "The number of private keys was incorrect")
			Expect(privateKeys[0]).Should(Equal("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"), "The private key did not get read correctly")
			Expect(privateKeys[1]).Should(Equal("59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"), "The private key did not get read correctly")
		})
	})
})
