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
```go
		diff, err := lc.MeasureDelta(&leak.CheckConfig{
			Query:          fmt.Sprintf(`quantile_over_time(0.5, container_memory_rss{name="don-node%d"}[1h]) / 1024 / 1024`, i),
			Start:          mustTime("2026-01-12T21:53:00Z"),
			End:            mustTime("2026-01-13T10:11:00Z"),
			WarmUpDuration: 1 * time.Hour,
		})
		require.NoError(t, err)
```

## Adding New Queries

You can use our test `hog` to debug new metrics and verify its correctness
```bash
cd framework/leak/cmd
just build
```

Run different hogs
```bash
ctf obs up
go test -v -timeout 1h -run TestCyclicHog
```
Then verify your query
```bash
go test -v -run TestVerifyCyclicHog
```
