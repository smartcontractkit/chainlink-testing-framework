package gauntlet_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gauntlet "github.com/smartcontractkit/integrations-framework/gauntlet"
)

var _ = Describe("Gauntlet @unit", func() {
	var ls string = "/usr/bin/ls"
	BeforeEach(func() {
		if gauntlet.GetOsVersion() == "macos" {
			ls = "/bin/ls"
		}
	})
	It("should fail to find the executable", func() {
		_, err := gauntlet.NewGauntlet("/path/to/nothing")
		Expect(err).Should(HaveOccurred(), "Successfully found an executable where one does not exist")
	})
	It("should create a new Gauntlet struct", func() {
		g, err := gauntlet.NewGauntlet(ls)
		Expect(err).ShouldNot(HaveOccurred(), "Could not get a new gauntlet struct")
		Expect(g.Network).Should(ContainSubstring("test"), "The network did not contain test")
	})
	It("should return a properly formatted flag", func() {
		g, err := gauntlet.NewGauntlet(ls)
		Expect(err).ShouldNot(HaveOccurred(), "Could not get a new gauntlet struct")
		Expect(g.Flag("flag", "value")).To(Equal("--flag=value"))
	})
	It("should execute a command correctly", func() {
		g, err := gauntlet.NewGauntlet(ls)
		Expect(err).ShouldNot(HaveOccurred(), "Could not get a new gauntlet struct")
		out, err := g.ExecCommand([]string{}, []string{})
		Expect(err).ShouldNot(HaveOccurred(), "Failed to execute a command")
		Expect(out).To(ContainSubstring("unrecognized option"))
	})
	It("should find an expected error in the output", func() {
		g, err := gauntlet.NewGauntlet(ls)
		Expect(err).ShouldNot(HaveOccurred(), "Could not get a new gauntlet struct")
		_, err = g.ExecCommandWithRetries([]string{}, []string{"unrecognized option"}, 1)
		Expect(err).Should(HaveOccurred(), "Failed to find the expected error")
	})
})
