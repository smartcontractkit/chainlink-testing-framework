# WASP - Testing User Journeys

Let's explore a more complex scenario where a user needs to authenticate first before performing an action. Additionally, we will introduce a slightly more advanced load profile:
* 1 user for the first 30 seconds
* 2 users for the next 30 seconds
* 3 users for the final 30 seconds

Since this is a "user journey," we will use a `VirtualUser` implementation to represent a user.

### Defining the Virtual User

First, let's define the `VirtualUser` struct:

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

Here, we’ve added a rate limiter to the struct, which will help limit the number of requests per second. This is useful to prevent overloading the server, especially in this test scenario where the requests are simple and fast. Without a limiter, even a small number of Virtual Users (VUs) could result in a very high RPS.

> [!NOTE]  
> Since `VirtualUser` does not inherently limit RPS (it depends on the server's processing speed), you should implement a rate-limiting mechanism if needed.

For brevity, we'll skip the implementation of the `Clone()` and `Teardown()` functions, as they are similar to previous examples. Additionally, no `Setup()` is required because we are using HTTP.

---

### Implementing Requests

#### User Authentication
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

#### Authenticated Action (e.g., Balance Check)
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

#### Combining the Requests
```go
func (m *VirtualUser) Call(l *wasp.Generator) {
	m.rl.Take() // apply rate limiting
	m.requestOne(l)
	m.requestTwo(l)
}
```

---

### Writing the Test

Now, let’s write the test. Pay attention to how the three phases of the load profile are defined under `Schedule`:

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

#### Load Profile
We generate load in three phases:
1. 1 user for the first 30 seconds
2. 2 users for the next 30 seconds
3. 3 users for the final 30 seconds

---

### Conclusion

And that's it! We’ve created a test that simulates a user journey with authentication and an action requiring authentication, while varying the load during execution. You can find the full example code [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/scenario).

But what if you wanted to combine multiple load generators—or even mix a `Gun` with a `VirtualUser`? Could you do that? Find out in the [next chapter](./profile_test.md).