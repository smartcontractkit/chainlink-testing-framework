package gauntlet_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gauntlet "github.com/smartcontractkit/integrations-framework/gauntlet"
)

var _ = Describe("Gauntlet @unit", func() {
	It("should create a new Gauntlet struct", func() {
		g, err := gauntlet.NewGauntlet()
		Expect(err).ShouldNot(HaveOccurred(), "Could not get a new gauntlet struct")
		Expect(g.Network).Should(ContainSubstring("test"), "The network did not contain test")
	})
	It("should return a properly formatted flag", func() {
		g, err := gauntlet.NewGauntlet()
		Expect(err).ShouldNot(HaveOccurred(), "Could not get a new gauntlet struct")
		Expect(g.Flag("flag", "value")).To(Equal("--flag=value"))
	})
	It("should execute a command correctly", func() {
		g, err := gauntlet.NewGauntlet()
		Expect(err).ShouldNot(HaveOccurred(), "Could not get a new gauntlet struct")
		out, err := g.ExecCommand([]string{}, []string{})
		Expect(err).Should(HaveOccurred(), "The command should technically always fail because we don't have access to a gauntlet executable, if it passed without error then we have an issue")
		Expect(out).To(ContainSubstring("yarn"), "Did not contain expected output")
	})
	It("should find an expected error in the output", func() {
		g, err := gauntlet.NewGauntlet()
		Expect(err).ShouldNot(HaveOccurred(), "Could not get a new gauntlet struct")
		_, err = g.ExecCommandWithRetries([]string{}, []string{"yarn"}, 1)
		Expect(err).Should(HaveOccurred(), "Failed to find the expected error")
	})
})
