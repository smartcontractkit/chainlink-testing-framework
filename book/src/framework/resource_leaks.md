## Resource Leak Detector

We have a simple utility to detect resource leaks in our tests

## CL Nodes Leak Detection

In this example test will fail if any node will consume more than 2 additional cores and allocate 20% more memory at the end of a test.
```go
import (
	"github.com/smartcontractkit/chainlink-testing-framework/framework/leak"
)
```

```go
			l, err := leak.NewCLNodesLeakDetector(leak.NewResourceLeakChecker())
			require.NoError(t, err)
			errs := l.Check(&leak.CLNodesCheck{
				NumNodes:        in.NodeSets[0].Nodes,
				Start:           start,
				End:             time.Now(),
				WarmUpDuration:  10 * time.Minute,
				CPUThreshold:    2000.0,
				MemoryThreshold: 20.0,
			})
			require.NoError(t, errs)
```

## Custom Resource Assertion

You can also use low-level API to verify