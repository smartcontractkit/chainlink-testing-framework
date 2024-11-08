## CTF modules and org dependencies
```mermaid
flowchart LR

	chainlink-testing-framework/examples_wasp --> chainlink-testing-framework/seth
	chainlink-testing-framework/examples_wasp --> chainlink-testing-framework/wasp
	click chainlink-testing-framework/examples_wasp href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/framework
	click chainlink-testing-framework/framework href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/framework/examples --> chainlink-testing-framework/framework
	chainlink-testing-framework/framework/examples --> chainlink-testing-framework/wasp
	click chainlink-testing-framework/framework/examples href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/grafana
	click chainlink-testing-framework/grafana href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/havoc --> chainlink-testing-framework/lib/grafana
	click chainlink-testing-framework/havoc href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/lib --> chainlink-testing-framework/seth
	chainlink-testing-framework/lib --> chainlink-testing-framework/wasp
	click chainlink-testing-framework/lib href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/lib/grafana
	click chainlink-testing-framework/lib/grafana href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/seth --> seth
	click chainlink-testing-framework/seth href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/tools/citool --> chainlink-testing-framework/lib
	click chainlink-testing-framework/tools/citool href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/tools/envresolve --> chainlink-testing-framework/lib
	click chainlink-testing-framework/tools/envresolve href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/tools/gotestloghelper --> chainlink-testing-framework/lib
	click chainlink-testing-framework/tools/gotestloghelper href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/wasp --> chainlink-testing-framework/grafana
	chainlink-testing-framework/wasp --> chainlink-testing-framework/lib/grafana
	click chainlink-testing-framework/wasp href "https://github.com/smartcontractkit/chainlink-testing-framework"
	chainlink-testing-framework/wasp-tests --> chainlink-testing-framework/wasp
	click chainlink-testing-framework/wasp-tests href "https://github.com/smartcontractkit/chainlink-testing-framework"
	seth
	click seth href "https://github.com/smartcontractkit/seth"

	subgraph chainlink-testing-framework-repo[chainlink-testing-framework]
		 chainlink-testing-framework/examples_wasp
		 chainlink-testing-framework/framework
		 chainlink-testing-framework/framework/examples
		 chainlink-testing-framework/grafana
		 chainlink-testing-framework/havoc
		 chainlink-testing-framework/lib
		 chainlink-testing-framework/lib/grafana
		 chainlink-testing-framework/seth
		 chainlink-testing-framework/tools/citool
		 chainlink-testing-framework/tools/envresolve
		 chainlink-testing-framework/tools/gotestloghelper
		 chainlink-testing-framework/wasp
		 chainlink-testing-framework/wasp-tests
	end
	click chainlink-testing-framework-repo href "https://github.com/smartcontractkit/chainlink-testing-framework"

	classDef outline stroke-dasharray:6,fill:none;
	class chainlink-testing-framework-repo outline
```
