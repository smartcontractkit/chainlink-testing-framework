package compatibility

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	"github.com/smartcontractkit/integrations-framework/utils"
)

func Test_Suite(t *testing.T) {
	utils.GinkgoSuite()
	RunSpecs(t, "Compatibility")
}
