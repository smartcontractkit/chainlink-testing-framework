package chaos_test

import (
	"github.com/smartcontractkit/integrations-framework/utils"
	"testing"

	. "github.com/onsi/ginkgo"
)

func Test_Suite(t *testing.T) {
	utils.GinkgoSuite()
	RunSpecs(t, "Chaos")
}
