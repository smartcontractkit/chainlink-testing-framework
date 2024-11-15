# WASP - Testing User Journeys

Let's look at a more complex scenario, where user needs to authenticate first to be able to perform some action. 
Also, we will introduce a slightly more complex load profile:
* 1 user for the first 30 seconds
* then 2 users for next 30 seconds
* and finally 3 users for the last 30 seconds

Since it's a "user journey", we will use a `VirtualUser` implementation to represent a user.

Again, let's start with defining the `VirtualUser` struct:
```go
type VirtualUser struct {
	*wasp.VUControl
	target    string
	Data      []string
	rateLimit int
	rl        ratelimit.Limiter
	client    *resty.Client
}
```

We have added a rate limiter to the struct, which will be used to limit the number of requests per second, as otherwise we could easily overload the server.

> [!WARNING]
> Contrary to a `Gun` the `VirtualUser` doesn't come with any kind of RPS control and you need to implement it yourself.

For brevity, let's skip the implementation of `Clone()` and `Teardown()` functions, as they are similar to the previous examples. Since we will be using HTTP, no `Setup()` is also necessary.

Now, let's implement a request that will authenticate the user:
```go
// requestOne represents user login
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

And an action that requires authentication, let's say a balance check:
```go
// represents authenticated user action
func (m *VirtualUser) requestTwo(l *wasp.Generator) {
	var result map[string]interface{}
	r, err := m.client.R().
		SetResult(&result).
		Get(m.target)
	if err != nil {
		l.Responses.Err(r, GroupUser, err)
		return
	}
	l.Responses.OK(r, GroupUser)
}
```

Now, let's use them:
```go
func (m *VirtualUser) Call(l *wasp.Generator) {
	m.rl.Take() // rate limit
	m.requestOne(l)
	m.requestTwo(l)
}
```

Now, the test itself. Pay attention to how we define the 3 phases of the load profile under `Schedule`:
```go
func TestScenario(t *testing.T) {
	srv := wasp.NewHTTPMockServer(nil)
	srv.Run()

	_, err := wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			T: t,
			LoadType: wasp.VU,
			VU:       NewExampleScenario(srv.URL()),
			Schedule: wasp.Combine(
				wasp.Plain(1, 30*time.Second),
				wasp.Plain(2, 30*time.Second),
				wasp.Plain(3, 30*time.Second),
			),
			LokiConfig: wasp.NewEnvLokiConfig(),
		})).Run(true)
	require.NoError(t, err)
}
```

And that's it! We have a test that simulates a user journey with authentication and an action that requires authentication. And it changes the load during execution!
As always, you can find the full example code [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/scenario).

But, what if you wanted to combine multiple load generators or even a `Gun` with a `VirtualUser`. Could you to that? You'll find out in the [next chapter](./profile_test.md).