# WASP - First Test

## Requirements
* Go installed
* Loki instance running (read how to start a local instance [here](./local_loki_grafana_stack))

Let's start by creating a very simple test that will send 5 HTTP requests per second for 60 seconds.

We will use a `Gun` that's designed for stateless protocols like HTTP and situations where we execute a single operation.

All that a `Gun` needs to do is implement this single-method interface:
```go
type Gun interface {
	Call(l *Generator) *Response
}
```

So here we go! First let's define a struct that will hold our `Gun` implementation:
```go
type ExampleGun struct {
	target string
	client *resty.Client
	Data   []string
}
```
Our gun will send a GET request to the target URL. if the request is successful, we will return `*wasp.Response` with the response data.
If it fails or responds with HTTP code different than `200`, we will return `*wasp.Response` with the response and error.

Let's implement the `Gun` interface:
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
> By default, WASP stores both all successful and failed responses, but you can use a `Sampler` to store only some of the successful ones. You can read more about it [here](./using_sampler.md).

Now that we have a gun, let's write the test:
```go
func TestGun(t *testing.T) {
	// start mock http server
	srv := wasp.NewHTTPMockServer(nil)
	srv.Run()

    // some parts omitted for brevity
	
	// create generator
	gen, err := wasp.NewGenerator(&wasp.Config{
		LoadType: wasp.RPS,
		// just use plain line profile - 5 RPS for 60s
		Schedule:   wasp.Plain(5, 60*time.Second),
		Gun:        NewExampleHTTPGun(srv.URL()),
		Labels:     labels,
		LokiConfig: wasp.NewEnvLokiConfig(),
	})
	if err != nil {
		panic(err)
	}
	// run the generator and wait until it finish
	gen.Run(true)
}
```

> [!NOTE]
> `wasp.NewEnvLokiConfig()` configures Loki via environments variables. You can read about them [here](./configuration.md).

And that's it! You have just created your first test using WASP. You can now run it using `go test -v -run TestGun`.

You can find the full example [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/simple_rps).


But, what if you want to test a more complex scenario?

First, we will look at [stateful protocol test](./stateful_test.md), like `WebSocket`. Then on a hypothetical [user journey test](./user_journey_test.md), where multiple requests need to be executed.