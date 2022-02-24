package soak_test

//revive:disable:dot-imports
import (
	"testing"

	"github.com/smartcontractkit/integrations-framework/utils"

	. "github.com/onsi/ginkgo/v2"
)

func Test_Suite(t *testing.T) {
	utils.GinkgoSuite(utils.ProjectRoot)
	RunSpecs(t, "Soak")
}
