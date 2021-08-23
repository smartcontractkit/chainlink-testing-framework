package config

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/tools"
)

// Test config files
const (
	specifiedConfig     string = "%s/config/test_configs/specified_config"
	badConfig           string = "%s/config/test_configs/bad_config"
	noPrivateKeysConfig string = "%s/config/test_configs/no_private_keys"
)

var _ = Describe("Config unit tests @unit", func() {

	Describe("Verify order of importance for environment variables and config files", func() {
		It("should load the default config file", func() {
			conf, err := NewConfig(LocalConfig, tools.ProjectRoot)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.KeepEnvironments).Should(Equal("Never"))
		})

		It("should load a specified file", func() {
			conf, err := NewConfig(LocalConfig, fmt.Sprintf(specifiedConfig, tools.ProjectRoot))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.KeepEnvironments).Should(Equal("Always"))
		})

		It("should fail to load a bad file", func() {
			_, err := NewConfig(LocalConfig, fmt.Sprintf(badConfig, tools.ProjectRoot))
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("line 1: cannot unmarshal"))
		})

		It("should overwrite default values with ENV variables", func() {
			os.Setenv("KEEP_ENVIRONMENTS", "OnFail")
			conf, err := NewConfig(LocalConfig, tools.ProjectRoot)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.KeepEnvironments).Should(Equal("OnFail"))
		})

		It("should overwrite specified file values with ENV variables", func() {
			os.Setenv("KEEP_ENVIRONMENTS", "OnFail")
			conf, err := NewConfig(LocalConfig, fmt.Sprintf(specifiedConfig, tools.ProjectRoot))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.KeepEnvironments).Should(Equal("OnFail"))

		})

		It("should load the config if we specify a Secret Config type", func() {
			conf, err := NewConfig(SecretConfig, tools.ProjectRoot)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(conf.KeepEnvironments).Should(Equal("Never"))
		})

		AfterEach(func() {
			os.Unsetenv("KEEP_ENVIRONMENTS")
		})
	})

	It("should get the network config with a valid name", func() {
		conf, err := NewConfig(LocalConfig, fmt.Sprintf(specifiedConfig, tools.ProjectRoot))
		Expect(err).ShouldNot(HaveOccurred())
		network, err := conf.GetNetworkConfig("test_this_geth")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(network.Name).Should(Equal("Tester Ted"))
	})

	It("should not get the network config with an invalid name", func() {
		conf, err := NewConfig(LocalConfig, fmt.Sprintf(specifiedConfig, tools.ProjectRoot))
		Expect(err).ShouldNot(HaveOccurred())
		_, err = conf.GetNetworkConfig("bad_name")
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(ContainSubstring("no supported network"))
	})

	It("should fetch private keys when they exist", func() {
		conf, err := NewConfig(LocalConfig, fmt.Sprintf(specifiedConfig, tools.ProjectRoot))
		Expect(err).ShouldNot(HaveOccurred())
		network, err := conf.GetNetworkConfig("test_this_geth")
		Expect(err).ShouldNot(HaveOccurred())
		privateKeys := NewPrivateKeyStore(LocalConfig, network)
		keys, err := privateKeys.Fetch()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(len(keys)).Should(Equal(1))
		Expect(keys[0]).Should((Equal("abcdefg")))
	})

	It("should not fetch private keys when they do not exist", func() {
		conf, err := NewConfig(LocalConfig, fmt.Sprintf(noPrivateKeysConfig, tools.ProjectRoot))
		Expect(err).ShouldNot(HaveOccurred())
		network, err := conf.GetNetworkConfig("test_this_geth")
		Expect(err).ShouldNot(HaveOccurred())
		privateKeys := NewPrivateKeyStore(LocalConfig, network)
		_, err = privateKeys.Fetch()
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(ContainSubstring("no keys found"))
	})

})
