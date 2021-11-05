package config_test

import (
	"fmt"
	"github.com/smartcontractkit/integrations-framework/client"
	"github.com/smartcontractkit/integrations-framework/utils"
	"os"

	"github.com/smartcontractkit/integrations-framework/config"
	"github.com/smartcontractkit/integrations-framework/environment"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Test config files
const (
	specifiedConfig string = "%s/config/test_configs/specified_config"
	badConfig       string = "%s/config/test_configs/bad_config"
	fetchConfig     string = "%s/config/test_configs/fetch_config"
)

var _ = Describe("Config unit tests @unit", func() {

	Describe("Verify order of importance for environment variables and config files", func() {
		It("should load the default config file", func() {
			conf, err := config.NewConfig(utils.ProjectRoot)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.KeepEnvironments).Should(Equal("Never"), "We either changed the default value in the config or it did not load correctly")
		})

		It("should load a specified file", func() {
			conf, err := config.NewConfig(fmt.Sprintf(specifiedConfig, utils.ProjectRoot))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.KeepEnvironments).Should(Equal("Always"), "We did not load the correct config file")
		})

		It("should fail to load a bad file", func() {
			_, err := config.NewConfig(fmt.Sprintf(badConfig, utils.ProjectRoot))
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("line 1: cannot unmarshal"))
		})

		It("should overwrite default values with ENV variables", func() {
			err := os.Setenv("KEEP_ENVIRONMENTS", "OnFail")
			Expect(err).ShouldNot(HaveOccurred())
			conf, err := config.NewConfig(utils.ProjectRoot)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.KeepEnvironments).Should(Equal("OnFail"), "The environment variable should have been used to change the config value")
		})

		It("should overwrite specified file values with ENV variables", func() {
			err := os.Setenv("KEEP_ENVIRONMENTS", "OnFail")
			Expect(err).ShouldNot(HaveOccurred())
			conf, err := config.NewConfig(fmt.Sprintf(specifiedConfig, utils.ProjectRoot))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.KeepEnvironments).Should(Equal("OnFail"), "The environment variable should have been used to change the config value")

		})

		It("should load the config if we specify a Secret Config type", func() {
			conf, err := config.NewConfig(utils.ProjectRoot)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.KeepEnvironments).Should(Equal("Never"), "We either changed the default value in the config or it did not load correctly")
		})

		AfterEach(func() {
			err := os.Unsetenv("KEEP_ENVIRONMENTS")
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	It("should get the network config with a valid name", func() {
		conf, err := config.NewConfig(fmt.Sprintf(specifiedConfig, utils.ProjectRoot))
		Expect(err).ShouldNot(HaveOccurred())
		network, err := conf.GetNetworkConfig("test_this_geth")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(network.Name).Should(Equal("Tester Ted"), "The network config was not loaded correctly")
	})

	It("should not get the network config with an invalid name", func() {
		conf, err := config.NewConfig(fmt.Sprintf(specifiedConfig, utils.ProjectRoot))
		Expect(err).ShouldNot(HaveOccurred())
		_, err = conf.GetNetworkConfig("bad_name")
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(ContainSubstring("no supported network"))
	})

	It("should fetch LocalStore raw keys", func() {
		conf, err := config.NewConfig(fmt.Sprintf(fetchConfig, utils.ProjectRoot))
		Expect(err).ShouldNot(HaveOccurred())

		bcNetwork, err := client.DefaultNetworkFromConfig(conf)
		Expect(err).ShouldNot(HaveOccurred())

		bcNetwork.Config().PrivateKeyStore, err = environment.NewPrivateKeyStoreFromEnv(&environment.K8sEnvironment{}, bcNetwork.Config())
		Expect(err).ShouldNot(HaveOccurred())

		privateKeys, err := bcNetwork.Config().PrivateKeyStore.Fetch()
		Expect(err).ShouldNot(HaveOccurred())

		Expect(privateKeys).Should(ContainElements(
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
			"5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a"))
	})
})
