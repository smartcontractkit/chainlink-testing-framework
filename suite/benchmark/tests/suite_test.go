package benchmark_test

//revive:disable:dot-imports
import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	"github.com/smartcontractkit/chainlink-testing-framework/actions"
)

func Test_Suite(t *testing.T) {
	actions.GinkgoRemoteSuite()
	RunSpecs(t, "Benchmark")
}
