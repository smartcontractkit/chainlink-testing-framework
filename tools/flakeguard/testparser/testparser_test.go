package testparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	defaultTestRunCount  = 5
	flakyTestPackagePath = "./example_test_package"
	debugDir             = "debug_outputs"
)

type expectedTestResult struct {
	allSuccesses  bool
	someSuccesses bool
	allFailures   bool
	someFailures  bool
	allSkips      bool
	testPanic     bool
	packagePanic  bool
	race          bool
	maximumRuns   int

	exactRuns       *int
	minimumRuns     *int
	exactPassRate   *float64
	minimumPassRate *float64
	maximumPassRate *float64

	seen bool
}

func TestAttributePanicToTest(t *testing.T) {
	t.Parallel()

	// Test cases: each test case contains a slice of output strings.
	testCases := []struct {
		name             string
		expectedTestName string
		expectedTimeout  bool
		outputs          []string
	}{
		{
			name:             "properly attributed panic",
			expectedTestName: "TestPanic",
			expectedTimeout:  false,
			outputs: []string{
				"panic: This test intentionally panics [recovered]",
				"\tpanic: This test intentionally panics",
				"goroutine 25 [running]:",
				"testing.tRunner.func1.2({0x1008cde80, 0x1008f7d90})",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1632 +0x1bc",
				"testing.tRunner.func1()",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1635 +0x334",
				"panic({0x1008cde80?, 0x1008f7d90?})",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/runtime/panic.go:785 +0x124",
				"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package.TestPanic(0x140000b6ea0?)",
			},
		},
		{
			name:             "improperly attributed panic",
			expectedTestName: "TestPanic",
			expectedTimeout:  false,
			outputs: []string{
				"panic: This test intentionally panics [recovered]",
				"TestPanic(0x140000b6ea0?)",
				"goroutine 25 [running]:",
				"testing.tRunner.func1.2({0x1008cde80, 0x1008f7d90})",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1632 +0x1bc",
				"testing.tRunner.func1()",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1635 +0x334",
				"panic({0x1008cde80?, 0x1008f7d90?})",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/runtime/panic.go:785 +0x124",
				"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package.TestPanic(0x140000b6ea0?)",
			},
		},
		{
			name:             "timeout panic",
			expectedTestName: "TestTimedOut",
			expectedTimeout:  true,
			outputs: []string{
				"panic: test timed out after 10m0s",
				"running tests",
				"TestTimedOut (10m0s)",
				"goroutine 397631 [running]:",
				"testing.(*M).startAlarm.func1()",
				"\t/opt/hostedtoolcache/go/1.23.3/x64/src/testing/testing.go:2373 +0x385",
				"created by time.goFunc",
				"/opt/hostedtoolcache/go/1.23.3/x64/src/time/sleep.go:215 +0x2d",
			},
		},
		{
			name:             "subtest panic",
			expectedTestName: "TestSubTestsSomePanic",
			expectedTimeout:  false,
			outputs: []string{
				"panic: This subtest always panics [recovered]",
				"panic: This subtest always panics",
				"goroutine 23 [running]:",
				"testing.tRunner.func1.2({0x100489e80, 0x1004b3e30})",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1632 +0x1bc",
				"testing.tRunner.func1()",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1635 +0x334",
				"panic({0x100489e80?, 0x1004b3e30?})",
				"\t/opt/homebrew/Cellar/go/1.23.2/libexec/src/runtime/panic.go:785 +0x124",
				"github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package.TestSubTestsSomePanic.func2(0x140000c81a0?)",
			},
		},
		{
			name:             "memory_test panic extraction",
			expectedTestName: "TestJobClientJobAPI",
			expectedTimeout:  false,
			outputs: []string{
				"panic: freeport: cannot allocate port block [recovered]",
				"\tpanic: freeport: cannot allocate port block",
				"goroutine 321 [running]:",
				"testing.tRunner.func1.2({0x5e0dd80, 0x72ebb40})",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1734 +0x21c",
				"testing.tRunner.func1()",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1737 +0x35e",
				"panic({0x5e0dd80?, 0x72ebb40?})",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/runtime/panic.go:787 +0x132",
				"github.com/hashicorp/consul/sdk/freeport.alloc()",
				"\t/home/runner/go/pkg/mod/github.com/hashicorp/consul/sdk@v0.16.1/freeport/freeport.go:274 +0xad",
				"github.com/hashicorp/consul/sdk/freeport.initialize()",
				"\t/home/runner/go/pkg/mod/github.com/hashicorp/consul/sdk@v0.16.1/freeport/freeport.go:124 +0x2d7",
				"sync.(*Once).doSlow(0xc0018eb600?, 0xc000da4a98?)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/sync/once.go:78 +0xab",
				"sync.(*Once).Do(...)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/sync/once.go:69",
				"github.com/hashicorp/consul/sdk/freeport.Take(0x1)",
				"\t/home/runner/go/pkg/mod/github.com/hashicorp/consul/sdk@v0.16.1/freeport/freeport.go:303 +0xe5",
				"github.com/hashicorp/consul/sdk/freeport.GetN({0x7337708, 0xc000683dc0}, 0x1)",
				"\t/home/runner/go/pkg/mod/github.com/hashicorp/consul/sdk@v0.16.1/freeport/freeport.go:427 +0x48",
				"github.com/smartcontractkit/chainlink/deployment/environment/memory_test.TestJobClientJobAPI(0xc000683dc0)",
				"\t/home/runner/work/chainlink/chainlink/deployment/environment/memory/job_service_client_test.go:116 +0xc6",
				"testing.tRunner(0xc000683dc0, 0x6d6c838)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1792 +0xf4",
				"created by testing.(*T).Run in goroutine 1",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1851 +0x413",
			},
		},
		{
			name:             "changeset_test panic extraction",
			expectedTestName: "TestDeployBalanceReader",
			expectedTimeout:  false,
			outputs: []string{
				"panic: freeport: cannot allocate port block [recovered]",
				"\tpanic: freeport: cannot allocate port block",
				"goroutine 378 [running]:",
				"testing.tRunner.func1.2({0x6063f40, 0x76367f0})",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1734 +0x21c",
				"testing.tRunner.func1()",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1737 +0x35e",
				"panic({0x6063f40?, 0x76367f0?})",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/runtime/panic.go:787 +0x132",
				"github.com/hashicorp/consul/sdk/freeport.alloc()",
				"\t/home/runner/go/pkg/mod/github.com/hashicorp/consul/sdk@v0.16.1/freeport/freeport.go:274 +0xad",
				"github.com/hashicorp/consul/sdk/freeport.initialize()",
				"\t/home/runner/go/pkg/mod/github.com/hashicorp/consul/sdk@v0.16.1/freeport/freeport.go:124 +0x2d7",
				"sync.(*Once).doSlow(0xa94f820?, 0xa8000a?)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/sync/once.go:78 +0xab",
				"sync.(*Once).Do(...)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/sync/once.go:69",
				"github.com/hashicorp/consul/sdk/freeport.Take(0x1)",
				"\t/home/runner/go/pkg/mod/github.com/hashicorp/consul/sdk@v0.16.1/freeport/freeport.go:303 +0xe5",
				"github.com/hashicorp/consul/sdk/freeport.GetN({0x7684150, 0xc000583c00}, 0x1)",
				"\t/home/runner/go/pkg/mod/github.com/hashicorp/consul/sdk@v0.16.1/freeport/freeport.go:427 +0x48",
				"github.com/smartcontractkit/chainlink/deployment/environment/memory.NewNodes(0xc000583c00, 0xff, 0xc001583d10, 0xc005aa0030, 0x1, 0x0, {0x0, {0x0, 0x0, 0x0, ...}, ...}, ...)",
				"\t/home/runner/work/chainlink/chainlink/deployment/environment/memory/environment.go:177 +0xa5",
				"github.com/smartcontractkit/chainlink/deployment/environment/memory.NewMemoryEnvironment(_, {_, _}, _, {0x2, 0x0, 0x0, 0x1, 0x0, {0x0, ...}})",
				"\t/home/runner/work/chainlink/chainlink/deployment/environment/memory/environment.go:223 +0x10c",
				"github.com/smartcontractkit/chainlink/deployment/keystone/changeset_test.TestDeployBalanceReader(0xc000583c00)",
				"\t/home/runner/work/chainlink/chainlink/deployment/keystone/changeset/deploy_balance_reader_test.go:23 +0xf5",
				"testing.tRunner(0xc000583c00, 0x70843d0)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1792 +0xf4",
				"created by testing.(*T).Run in goroutine 1",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/testing/testing.go:1851 +0x413",
				"    logger.go:146: 03:14:04.485880684\tINFO\tDeployed KeystoneForwarder 1.0.0 chain selector 909606746561742123 addr 0x72B66019aCEdc35F7F6e58DF94De95f3cBCC5971\t{\"version\": \"(devel)@unset\"}",
				"    logger.go:146: 03:14:04.486035865\tINFO\tdeploying forwarder\t{\"version\": \"(devel)@unset\", \"chainSelector\": 5548718428018410741}",
				"    logger.go:146: 2025-03-08T03:14:04.490Z\tINFO\tchangeset/jd_register_nodes.go:91\tregistered node\t{\"version\": \"unset@unset\", \"name\": \"node1\", \"id\": \"node:{id:\\\"895776f5ba0cc11c570a47b5cc3dbb8771da9262cfb545cd5d48251796af7f\\\"  public_key:\\\"895776f5ba0cc11c570a47b5cc3dbb8771da9262cfb545cd5d48251796af7f\\\"  is_enabled:true  is_connected:true  labels:{key:\\\"product\\\"  value:\\\"test-product\\\"}  labels:{key:\\\"environment\\\"  value:\\\"test-env\\\"}  labels:{key:\\\"nodeType\\\"  value:\\\"bootstrap\\\"}  labels:{key:\\\"don-0-don1\\\"}\"}",
			},
		},
		{
			name:             "empty",
			expectedTestName: "",
			expectedTimeout:  false,
			outputs:          []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testName, timeout, err := attributePanicToTest(tc.outputs)
			assert.Equal(t, tc.expectedTimeout, timeout, "timeout flag mismatch")
			if tc.expectedTestName == "" {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedTestName, testName, "test name mismatch")
			}
		})
	}
}

func TestFailToAttributePanicToTest(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		outputs []string
	}{
		{
			name: "no test name in panic",
			outputs: []string{
				"panic: reflect: Elem of invalid type bool",
				"goroutine 104182 [running]:",
				"reflect.elem(0xc0569d9998?)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/reflect/type.go:733 +0x9a",
				"reflect.(*rtype).Elem(0xa4dd940?)",
				"\t/opt/hostedtoolcache/go/1.24.0/x64/src/reflect/type.go:737 +0x15",
				"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainreader.setPollingFilterOverrides(0x0, {0xc052040510, 0x1, 0xc?})",
				"\t/home/runner/go/pkg/mod/github.com/smartcontractkit/chainlink-solana@v1.1.2-0.20250319030827-8e2f4d76eb79/pkg/solana/chainreader/chain_reader.go:942 +0x492",
				"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainreader.(*ContractReaderService).addEventRead(_, _, {_, _}, {_, _}, {{0xc0544c4270, 0x9}, {0xc0544c4280, 0xc}, ...}, ...)",
				"\t/home/runner/go/pkg/mod/github.com/smartcontractkit/chainlink-solana@v1.1.2-0.20250319030827-8e2f4d76eb79/pkg/solana/chainreader/chain_reader.go:605 +0x13d",
				"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainreader.(*ContractReaderService).initNamespace(0xc054472540, 0xc01c37d440?)",
				"\t/home/runner/go/pkg/mod/github.com/smartcontractkit/chainlink-solana@v1.1.2-0.20250319030827-8e2f4d76eb79/pkg/solana/chainreader/chain_reader.go:443 +0x28b",
				"github.com/smartcontractkit/chainlink-solana/pkg/solana/chainreader.NewContractReaderService({0x7fcf8b532040?, 0xc015b223e0?}, {0xc6ac960, 0xc05464e470}, {0xc0544384e0?, {0xc01c37d440?, 0xc054163b84?, 0xc054163b80?}}, {0x7fcf8071c7a0, 0xc0157928c0})",
				"\t/home/runner/go/pkg/mod/github.com/smartcontractkit/chainlink-solana@v1.1.2-0.20250319030827-8e2f4d76eb79/pkg/solana/chainreader/chain_reader.go:97 +0x287",
				"github.com/smartcontractkit/chainlink-solana/pkg/solana.(*Relayer).NewContractReader(0xc015b2e150, {0x4d0102030cb384f5?, 0xb938300b5ca1aa13?}, {0xc05469c000, 0x1eedf, 0x20000})",
				"\t/home/runner/go/pkg/mod/github.com/smartcontractkit/chainlink-solana@v1.1.2-0.20250319030827-8e2f4d76eb79/pkg/solana/relay.go:160 +0x205",
				"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/oraclecreator.(*pluginOracleCreator).createReadersAndWriters(_, {_, _}, {_, _}, _, {0x3, {0x0, 0xa, 0x93, ...}, ...}, ...)",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/ccip/oraclecreator/plugin.go:446 +0x338",
				"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/oraclecreator.(*pluginOracleCreator).Create(0xc033a69ad0, {0xc6f5a10, 0xc02e4f9a40}, 0x3, {0x3, {0x0, 0xa, 0x93, 0x8f, 0x67, ...}, ...})",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/ccip/oraclecreator/plugin.go:215 +0xc0c",
				"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/launcher.createDON({0xc6f5a10, 0xc02e4f9a40}, {0x7fcf8b533ad0, 0xc015b97340}, {0xb6, 0x5e, 0x31, 0xd0, 0x35, 0xef, ...}, ...)",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/ccip/launcher/launcher.go:367 +0x451",
				"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/launcher.(*launcher).processAdded(0xc015723080, {0xc6f5a10, 0xc02e4f9a40}, 0xc053de2ff0)",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/ccip/launcher/launcher.go:254 +0x239",
				"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/launcher.(*launcher).processDiff(0xc015723080, {0xc6f5a10, 0xc02e4f9a40}, {0xc053de2ff0?, 0xc053de3020?, 0xc053de3050?})",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/ccip/launcher/launcher.go:192 +0x68",
				"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/launcher.(*launcher).tick(0xc015723080, {0xc6f5a10, 0xc02e4f9a40})",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/ccip/launcher/launcher.go:178 +0x20b",
				"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/launcher.(*launcher).monitor(0xc015723080)",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/ccip/launcher/launcher.go:152 +0x112",
				"created by github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/launcher.(*launcher).Start.func1 in goroutine 1335",
				"\t/home/runner/work/chainlink/chainlink/core/capabilities/ccip/launcher/launcher.go:134 +0xa5",
				"FAIL\tgithub.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana\t184.801s",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testName, timeout, err := attributePanicToTest(tc.outputs)
			require.Error(t, err)
			assert.Empty(t, testName, "test name should be empty")
			assert.False(t, timeout, "timeout flag should be false")
		})
	}
}

func TestAttributeRaceToTest(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		packageName      string
		expectedTestName string
		raceEntries      []entry
	}{
		{
			name:             "properly attributed race",
			expectedTestName: "TestRace",
			packageName:      "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package",
			raceEntries:      properlyAttributedRaceEntries,
		},
		{
			name:             "improperly attributed race",
			expectedTestName: "TestRace",
			packageName:      "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package",
			raceEntries:      improperlyAttributedRaceEntries,
		},
		{
			name:        "empty",
			packageName: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/example_test_package",
			raceEntries: []entry{
				{},
			},
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			testName, err := attributeRaceToTest(tc.packageName, tc.raceEntries)
			if tc.expectedTestName == "" {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedTestName, testName, "test race not attributed correctly")
			}
		})
	}
}

var (
	improperlyAttributedRaceEntries = []entry{
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "WARNING: DATA RACE\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Read at 0x00c000292028 by goroutine 13:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0x94\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Previous write at 0x00c000292028 by goroutine 12:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0xa4\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Goroutine 13 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Goroutine 12 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "WARNING: DATA RACE\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Write at 0x00c000292028 by goroutine 13:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0xa4\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Previous write at 0x00c000292028 by goroutine 14:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0xa4\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Goroutine 13 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Goroutine 14 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "WARNING: DATA RACE\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Read at 0x00c000292028 by goroutine 19:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:68 +0xb8\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Previous write at 0x00c000292028 by goroutine 13:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0xa4\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Goroutine 19 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "Goroutine 13 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestFlaky", Output: "    testing.go:1399: race detected during execution of test\n"},
	}
	properlyAttributedRaceEntries = []entry{
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "WARNING: DATA RACE\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Read at 0x00c000292028 by goroutine 13:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0x94\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Previous write at 0x00c000292028 by goroutine 12:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0xa4\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Goroutine 13 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Goroutine 12 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "WARNING: DATA RACE\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Write at 0x00c000292028 by goroutine 13:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0xa4\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Previous write at 0x00c000292028 by goroutine 14:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.func1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:67 +0xa4\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x44\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Goroutine 13 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "Goroutine 14 (running) created at:\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package.TestRace()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /Users/adamhamrick/Projects/chainlink-testing-framework/tools/flakeguard/runner/example_test_package/example_tests_test.go:74 +0x158\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.tRunner()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1690 +0x184\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "  testing.(*T).Run.gowrap1()\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "      /opt/homebrew/Cellar/go/1.23.2/libexec/src/testing/testing.go:1743 +0x40\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "==================\n"},
		{Action: "output", Package: "github.com/smartcontractkit/chainlink-testing-framework/tools/flakeguard/runner/example_test_package", Test: "TestRace", Output: "    testing.go:1399: race detected during execution of test\n"},
	}
)
