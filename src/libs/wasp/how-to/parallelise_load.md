# WASP - How to Parallelize Load

Parallelizing load can be achieved using the **`Profile`** component, which allows you to combine multiple generators. However, this approach works only if all generators start at the same time. If you need to space out the generators for various reasons, you must use native Go concurrency primitives like `goroutines` and `channels`.

---

### Concept

To parallelize load:
1. Split the load into multiple parts.
2. Run each part in separate `goroutines`, either using a `Profile` or directly with `Generator`.
3. Use `channels` to coordinate the timing and synchronization between the goroutines.

---

### Example Scenario

Suppose you want to execute the following scenario:
1. Gradually ramp up an **RPS load** over 10 seconds from 1 to 10 RPS and hold it for 50 seconds.
2. When RPS reaches 10, gradually ramp up a **VU load** over 16 seconds from 2 to 8 VUs (`VU'`) and hold it for 30 seconds.
3. Once `VU'` reaches 8, introduce another user interaction and ramp up a **VU load** over 14 seconds from 3 to 9 VUs (`VU''`) and hold it for 20 seconds.

Here’s how you can achieve this:

```go
func TestParallelLoad(t *testing.T) {
    labels := map[string]string{
        "branch": "parallel_load",
        "commit": "parallel_load",
    }

    // Define RPS schedule
    rpsSchedule := wasp.Combine(
		// wasp.Steps(from, increase, steps, duration)
        wasp.Steps(1, 1, 9, 10*time.Second), // Start with 1 RPS, increment by 1 RPS in 9 steps over 10 seconds
        // wasp.Plain(count, duration)		
        wasp.Plain(9, 50*time.Second),       // Hold 9 RPS for 50 seconds
    )

    // Define VU' schedule
    vuSchedule := wasp.Combine(
        // wasp.Steps(from, increase, steps, duration)
        wasp.Steps(2, 1, 8, 16*time.Second), // Start with 2 VUs, increment by 1 VU in 8 steps over 16 seconds
        // wasp.Plain(count, duration)
        wasp.Plain(10, 30*time.Second),      // Hold 10 VUs for 30 seconds
    )

    // Define VU'' schedule
    vu2Schedule := wasp.Combine(
        // wasp.Steps(from, increase, steps, duration)
        wasp.Steps(3, 1, 6, 14*time.Second), // Start with 3 VUs, increment by 1 VU in 6 steps over 14 seconds
		// wasp.Plain(count, duration)
        wasp.Plain(9, 20*time.Second),       // Hold 9 VUs for 20 seconds
    )

    // Create generators
    rpsGen, err := wasp.NewGenerator(&wasp.Config{
        LoadType:   wasp.RPS,
        Schedule:   rpsSchedule,
        GenName:    "Kappa",
        Labels:     labels,
        Gun:        NewExampleHTTPGun(srv.URL()),
        LokiConfig: wasp.NewEnvLokiConfig(),
    })
    require.NoError(t, err)

    vuGen, err := wasp.NewGenerator(&wasp.Config{
        LoadType:   wasp.VU,
        Schedule:   vuSchedule,
        GenName:    "Lambda",
        Labels:     labels,
        VU:         NewExampleScenario(srv.URL()),
        LokiConfig: wasp.NewEnvLokiConfig(),
    })
    require.NoError(t, err)

    vu2Gen, err := wasp.NewGenerator(&wasp.Config{
        LoadType: wasp.VU,
        Schedule: vu2Schedule,
        GenName:  "Mu",
        Labels:   labels,
        VU:       NewExampleScenario(srv.URL()), // Use the same VirtualUser implementation for simplicity
        LokiConfig: wasp.NewEnvLokiConfig(),
    })
    require.NoError(t, err)

    wg := sync.WaitGroup{}

    // Run RPS load in a separate goroutine
    go func() {
        wg.Add(1)
        rpsGen.Run(true)
        wg.Done()
    }()

    // Wait for RPS load to stabilize
    time.Sleep(10 * time.Second)

    // Run VU' load in a separate goroutine
    go func() {
        wg.Add(1)
        vuGen.Run(true)
        wg.Done()
    }()

    // Wait for VU' load to stabilize
    time.Sleep(16 * time.Second)

    // Run VU'' load
    vu2Gen.Run(true)

    // Wait for all goroutines to complete
    wg.Wait()

    // Check for load generation errors
    require.Equal(t, 0, len(rpsGen.Errors()), "RPS generator errors")
    require.Equal(t, 0, len(vuGen.Errors()), "VU generator errors")
    require.Equal(t, 0, len(vu2Gen.Errors()), "VU'' generator errors")
}
```

---

### Key Points

- **Parallel Execution:** Run each generator in its own `goroutine`.
- **Synchronization:** Use `time.Sleep` or channels to coordinate timing between load segments.
- **Error Handling:** Always check for errors in load generation, though they may not necessarily indicate a problem (e.g., aggressive load tests).

---

### Full Example

For a complete example, refer to [this file](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/wasp/examples/profiles/node_background_load_test.go).