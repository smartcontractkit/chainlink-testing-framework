# WASP - First Test (RPS Test)

## Requirements
* Go installed
* Loki

Let's start by creating a simple test that sends 5 HTTP requests per second for 60 seconds.

We will use a `Gun`, which is designed for stateless protocols like HTTP or measuring throughput.

A `Gun` only needs to implement this single-method interface:

```go
type Gun interface {
	Call(l *Generator) *Response
}
```

### Defining the Gun

First, let's define a struct that will hold our `Gun` implementation:

```go
type ExampleGun struct {
	target string
	client *resty.Client
	Data   []string
}
```

Our `Gun` will send a `GET` request to the target URL. If the request is successful, we return a `*wasp.Response` containing the response data. If it fails or responds with an HTTP status other than `200`, we also return a `*wasp.Response` with the response and an error.

Here’s the implementation of the `Gun` interface:

```go
func (m *ExampleGun) Call(l *wasp.Generator) *wasp.Response {
	var result map[string]interface{}
	r, err := m.client.R().
		SetResult(&result).
		Get(m.target)
	if err != nil {
		return &wasp.Response{Data: result, Error: err.Error()}
	}
	if r.Status() != "200 OK" {
		return &wasp.Response{Data: result, Error: "not 200"}
	}
	return &wasp.Response{Data: result}
}
```

> [!NOTE]  
> By default, WASP stores all successful and failed responses, but you can use a `Sampler` to store only some of the successful ones. You can read more about it [here](./components/sampler.md).

### Writing the Test

Now that we have a `Gun`, let’s write the test:

```go
func TestGun(t *testing.T) {
	// start mock HTTP server
	srv := wasp.NewHTTPMockServer(nil)
	srv.Run()

	// define labels to differentiate one run from another
	labels := map[string]string{
		// check variables in dashboard/dashboard.go
		"go_test_name": "TestGun",
		"gen_name":     "test_gun",
		"branch":       "my-awesome-branch",
		"commit":       "f3729fa",
	}

	// create generator
	gen, err := wasp.NewGenerator(&wasp.Config{
		LoadType: wasp.RPS,
		// plain line profile - 5 RPS for 60s
		Schedule:   wasp.Plain(5, 60*time.Second),
		Gun:        NewExampleHTTPGun(srv.URL()),
		Labels:     labels,
		LokiConfig: wasp.NewEnvLokiConfig(),
	})
	if err != nil {
		panic(err)
	}
	// run the generator and wait until it finishes
	gen.Run(true)
}
```

> [!NOTE]  
> We used the `LoadType` of `wasp.RPS` since this is the only type of load that a `Gun` can handle.  
> You can read more about load types [here](./how-to/chose_rps_vu.md).

> [!NOTE]  
> You can learn more about different labels and their functions [here](./how-to/use_labels.md).

> [!NOTE]  
> `wasp.NewEnvLokiConfig()` configures Loki using environment variables. You can read about them [here](./configuration.md).

### Conclusion

And that's it! You’ve just created your first test using WASP. You can now run it using:

```bash
go test -v -run TestGun
```

You can find the full example [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/simple_rps).

---

### What’s Next?

What if you want to test a more complex scenario?

First, we’ll look at a [stateful protocol test](./stateful_test.md), like `WebSocket`. Then, we’ll explore a hypothetical [user journey test](./user_journey_test.md), where multiple requests need to be executed.