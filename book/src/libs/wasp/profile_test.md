# WASP - Using Profiles

Finally, let's look at the most complex scenario, where we will be using a `VirtualUser` to represent a user, and a `Gun` to represent some background load.
Also, we will use a new load segment type, one that varies over time.

In order to bind together different `Generators` (like a `Gun` or a `VirtualUser`) we will use a `Profile`. You can think of it as a highest-level abstraction that fully describes the load profile.

We will skip defining both the `Gun` and the `VirtualUser` as they are identical to the previous examples. Well, almost.
For the `VU` we will use a handy wrapper for sending `Response` back to the channel. See if you can spot it:
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

Let's define our background load schedule for the RPS `Gun`. During first 10 seconds, it will increment the RPS by 2.5 every seconds, and then it will keep it at 5 RPS for 40 seconds: 
```go
epsilonSchedule := wasp.Combine(
    wasp.Steps(1, 1, 4, 10*time.Second), // start with 1 RPS, increment by 1 RPS in 4 steps during 10 seconds (increment of 1 every 2.5 seconds)
    wasp.Plain(5, 40*time.Second)) // hold 5 RPS for 40 seconds
```

And our virtual user, will first do nothing for the first 10 seconds, waiting for the background load to reach the desired level and then add 1 virtual user every 3 seconds for 30 seconds, to finally wind down from 10 to 0 users in 10 seconds:
```go
	thetaSchedule := wasp.Combine(
		wasp.Plain(0, 10*time.Second), // do nothing for 10 seconds
		wasp.Steps(1, 1, 9, 30*time.Second), // start with 1 user, increment by 1 virtual user in 9 steps during 30 seconds (increment of 1 every ~3 seconds)
        wasp.Steps(10, -1, 10, 10*time.Second)) // start with 10 users, decrement by 1 virtual user in 10 steps during 10 seconds (decrement of 1 every second)
```

Now, let's define our `Profile`:
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

And done! You have just created a complex load profile that will simulate a growing number of users and a background load that varies over time.
Notice the handly `.Run(true)` method that will block until all of the `Profile`'s generators have finished.

You can find the full example [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/profiles).

Now that you know how to generate load, it's time to learn how to monitor it and assert on it. Let's move on to the next section: [Testing Alerts](./testing_alerts.md).
```