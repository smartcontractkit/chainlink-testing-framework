# WASP - Stateful protocol testing

WASP allows you to test stateful protocols, like `WebSocket`. They are a bit more complex than stateless protocols, which is reflected in slightly higher complexity of interface to be implemented, namely the `VirtualUser` interface:

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

As before, let's start by defining a struct that will hold our `VirtualUser` implementation:

```go
type WSVirtualUser struct {
	target string
	*wasp.VUControl	
	conn   *websocket.Conn
	Data   []string
}
```

We will start with implementing the `Clone()` function used by WASP for creating new instances of the `VirtualUser`:

```go
func (m *WSVirtualUser) Clone(_ *wasp.Generator) wasp.VirtualUser {
	return &WSVirtualUser{
		VUControl: wasp.NewVUControl(),
		target:    m.target,
		Data:      make([]string, 0),
	}
}
```

Next, we will implement the `Setup()` function, which will be used to establish a connection to the WebSocket server:

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

We will omit `Teardown()` for brevity, but it should be used to close the connection to the WebSocket server.
Also, we don't have to implement neither `Stop()` nor `StopChan()` functions, as they are already implemented in the `VUControl` struct.

Next, we will implement the `Call()` function, which will be used to receive messages from the WebSocket server:
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
As you can see, we are not returning a response directly from `Call()`, but sending it to the `ResponsesChan` channel. 

Now, let's write the test.
```go
func TestVirtualUser(t *testing.T) {
	// start mock ws server
	s := httptest.NewServer(wasp.MockWSServer{
		Sleep: 50 * time.Millisecond,
	})
	defer s.Close()
	time.Sleep(1 * time.Second)
	
	// some parts omitted for brevity

	// create generator
	gen, err := wasp.NewGenerator(&wasp.Config{
		LoadType: wasp.VU,
		// just use plain line profile - 5 VUs for 60s
		Schedule:   wasp.Plain(5, 60*time.Second),
		VU:         NewExampleVirtualUser(url),
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

That wasn't so difficult, was it? You can now test your WebSocket server with WASP. You can find a full example [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/wasp/examples/simple_vu).

Now, let's see how we can test a more complex scenario, where a `VirtualUser` needs to perform various operations.