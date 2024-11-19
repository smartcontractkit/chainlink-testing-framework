# WASP - Stateful Protocol Testing

WASP allows you to test stateful protocols like `WebSocket`. These protocols are more complex than stateless protocols, which is reflected in the slightly higher complexity of the interface to be implemented: the `VirtualUser` interface.

```go
type VirtualUser interface {
	Call(l *Generator)
	Clone(l *Generator) VirtualUser
	Setup(l *Generator) error
	Teardown(l *Generator) error
	Stop(l *Generator)
	StopChan() chan struct{}
}
```

### Defining the Virtual User

As before, let's start by defining a struct that will hold our `VirtualUser` implementation:

```go
type WSVirtualUser struct {
	target string
	*wasp.VUControl	
	conn   *websocket.Conn
	Data   []string
}
```

### Implementing the Clone Method

We will begin by implementing the `Clone()` function, used by WASP to create new instances of the `VirtualUser`:

```go
func (m *WSVirtualUser) Clone(_ *wasp.Generator) wasp.VirtualUser {
	return &WSVirtualUser{
		VUControl: wasp.NewVUControl(),
		target:    m.target,
		Data:      make([]string, 0),
	}
}
```

### Implementing the Setup Method

Next, we implement the `Setup()` function, which establishes a connection to the WebSocket server:

```go
func (m *WSVirtualUser) Setup(l *wasp.Generator) error {
	var err error
	m.conn, _, err = websocket.Dial(context.Background(), m.target, &websocket.DialOptions{})
	if err != nil {
		l.Log.Error().Err(err).Msg("failed to connect from vu")
		_ = m.conn.Close(websocket.StatusInternalError, "")
		return err
	}
	return nil
}
```

We will omit the `Teardown()` function for brevity, but it should be used to close the connection to the WebSocket server.  
Additionally, we do not need to implement the `Stop()` or `StopChan()` functions because they are already implemented in the `VUControl` struct.

### Implementing the Call Method

Now, we implement the `Call()` function, which is used to receive messages from the WebSocket server:

```go
func (m *WSVirtualUser) Call(l *wasp.Generator) {
	startedAt := time.Now()
	v := map[string]string{}
	err := wsjson.Read(context.Background(), m.conn, &v)
	if err != nil {
		l.Log.Error().Err(err).Msg("failed read ws msg from vu")
		l.ResponsesChan <- &wasp.Response{StartedAt: &startedAt, Error: err.Error(), Failed: true}
		return
	}
	l.ResponsesChan <- &wasp.Response{StartedAt: &startedAt, Data: v}
}
```

As you can see, instead of returning a single response directly from `Call()`, we send it to the `ResponsesChan` channel. 
This is done, so that we can send each response independently to Loki instead of waiting for the whole call to finish.

### Writing the Test

Now, let's write the test:

```go
func TestVirtualUser(t *testing.T) {
	// start mock WebSocket server
	s := httptest.NewServer(wasp.MockWSServer{
		Sleep: 50 * time.Millisecond,
	})
	defer s.Close()
	time.Sleep(1 * time.Second)

	// some parts omitted for brevity

	// create generator
	gen, err := wasp.NewGenerator(&wasp.Config{
		LoadType: wasp.VU,
		// plain line profile - 5 VUs for 60s
		Schedule:   wasp.Plain(5, 60*time.Second),
		VU:         NewExampleVirtualUser(url),
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

### Conclusion

That wasn’t so difficult, was it? You can now test your WebSocket server with WASP. You can find a full example [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/simple_vu).

---

### What’s Next?

Now, let’s explore how to test [a more complex scenario](./user_journey_test.md) where a `VirtualUser` needs to perform various operations.