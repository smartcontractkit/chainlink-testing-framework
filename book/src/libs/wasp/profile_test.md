# WASP - Using Profiles

In this section, we’ll explore the most complex scenario: using a `VirtualUser` to represent a user and a `Gun` to generate background load. We’ll also introduce a new load segment type that varies over time.

To bind together different `Generators` (such as a `Gun` or a `VirtualUser`), we’ll use a `Profile`. Think of it as the highest-level abstraction that fully describes the load profile.

### `Gun` and `VirtualUser`

We’ll skip defining both the `Gun` and the `VirtualUser`, as they are nearly identical to previous examples. However, for the `VirtualUser`, we’ll use a handy wrapper for sending `Response` back to the channel. See if you can spot it:

```go
// represents user login
func (m *VirtualUser) requestOne(l *wasp.Generator) {
	var result map[string]interface{}
	r, err := m.client.R().
		SetResult(&result).
		Get(m.target)
	if err != nil {
		l.Responses.Err(r, GroupAuth, err)
		return
	}
	l.Responses.OK(r, GroupAuth)
}
```

> [!NOTE]
> You might have noticed `GroupAuth` constant passed to `OK()` and `Err()` methods. This is used to group responses in the dashboard. 
> You can read more about it [here](./how-to/use_labels.md).

---

### Background Load Schedule for the `Gun`

For the RPS `Gun`, we’ll define a schedule where:
1. During the first 10 seconds, the RPS increases by 2.5 every second.
2. For the next 40 seconds, it remains steady at 5 RPS.

```go
epsilonSchedule := wasp.Combine(
    // wasp.Steps(from, increase, steps, duration)	
    wasp.Steps(1, 1, 4, 10*time.Second), // Start at 1 RPS, increment by 1 RPS in 4 steps over 10 seconds (1 increment every 2.5 seconds)
    // wasp.Plain(count, duration)
    wasp.Plain(5, 40*time.Second))       // Hold 5 RPS for 40 seconds
```

---

### Virtual User Schedule

For the `VirtualUser`, the schedule will:
1. Start with 1 user for the first 10 seconds.
2. Add 1 user every 3 seconds for 30 seconds.
3. Gradually reduce from 10 users to 0 over 10 seconds.

```go
thetaSchedule := wasp.Combine(
    // wasp.Plain(count, duration)
    wasp.Plain(1, 10*time.Second),         // 1 user for the first 10 seconds
    // wasp.Steps(from, increase, steps, duration)
    wasp.Steps(1, 1, 9, 30*time.Second),  // Increment by 1 user every ~3 seconds over 30 seconds
    // wasp.Steps(from, increase, steps, duration)
    wasp.Steps(10, -1, 10, 10*time.Second)) // Decrement by 1 user every second over 10 seconds
```

---

### Defining the Profile

We’ll now define our `Profile` to combine both the `Gun` and the `VirtualUser`:

```go
_, err := wasp.NewProfile().
    Add(wasp.NewGenerator(&wasp.Config{
        T:          t,
        LoadType:   wasp.VU,
        GenName:    "Theta",
        Schedule:   thetaSchedule,
        VU:         NewExampleScenario(srv.URL()),
        LokiConfig: wasp.NewEnvLokiConfig(),
    })).
    Add(wasp.NewGenerator(&wasp.Config{
        T:          t,
        LoadType:   wasp.RPS,
        GenName:    "Epsilon",
        Schedule:   epsilonSchedule,
        Gun:        NewExampleHTTPGun(srv.URL()),
        LokiConfig: wasp.NewEnvLokiConfig(),
    })).
    Run(true)
```

---

### Conclusion

And that’s it! You’ve created a complex load profile that simulates a growing number of users alongside a background load that varies over time. Notice the `.Run(true)` method, which blocks until all the `Profile`'s generators have finished.

You can find the full example [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/profiles).

---

> [!NOTE]  
> The `error` returned by the `.Run(true)` function of a `Profile` might indicate the following:
> * load generator did not start and instead returned an error
> * or **if** `Profile` was created with `WithGrafana` option that enabled `CheckDashboardAlertsAfterRun` that at some alerts fired.
> To check for errors during load generation, you need to call the `Errors()` method on each `Generator` within the `Profile`.

> [!NOTE]  
> Currently, it’s not possible to have a "waiting" schedule that doesn’t generate any load. 
> To implement such logic, start one `Profile` in a goroutine and use `time.Sleep` before starting the second `Profile` or start use a combination `.Run(false)` and `.Wait()` methods.
> Both examples of can be found [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/profiles/node_background_load_test.go).

---

### What’s Next?

Now that you know how to generate load, it’s time to learn how to monitor and assert on it. Let’s move on to the next section: [Testing Alerts](./testing_alerts.md).