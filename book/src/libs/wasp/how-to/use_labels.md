# WASP - How to Use Labels

The default WASP dashboard utilizes the following labels to differentiate between various generators and requests:
* `go_test_name`
* `gen_name`
* `branch`
* `commit`
* `call_group`

---

## `go_test_name`

This label is used to differentiate between different tests. If an instance of `*testing.T` is passed to the generator configuration, the test name will automatically be added to the `go_test_name` label:

```go
gen, err := wasp.NewGenerator(&wasp.Config{
    LoadType:   wasp.RPS,
    T:          t,                                  // < ----------------------- HERE
    Schedule:   wasp.Plain(5, 60*time.Second),
    Gun:        NewExampleHTTPGun(srv.URL()),
    Labels:     labels,
    LokiConfig: wasp.NewEnvLokiConfig(),
})
```

---

## `gen_name`

This label differentiates between different generators:

```go
gen, err := wasp.NewGenerator(&wasp.Config{
    LoadType:   wasp.RPS,
    T:          t,
    GenName:    "my_super_generator",                // < ----------------------- HERE
    Schedule:   wasp.Plain(5, 60*time.Second),
    Gun:        NewExampleHTTPGun(srv.URL()),
    Labels:     labels,
    LokiConfig: wasp.NewEnvLokiConfig(),
})
```

---

## `branch` and `commit`

These labels are primarily used to distinguish between different branches and commits when running tests in CI.  
Since there is no automated way to fetch this data, it must be passed manually:

```go
branch := os.Getenv("BRANCH")
commit := os.Getenv("COMMIT")

labels := map[string]string{
    "branch":      branch,                       // < ----------------------- HERE
    "commit":      commit,                       // < ----------------------- HERE
}

gen, err := wasp.NewGenerator(&wasp.Config{
    LoadType:   wasp.RPS,
    T:          t,
    GenName:    "my_super_generator",
    Schedule:   wasp.Plain(5, 60*time.Second),
    Gun:        NewExampleHTTPGun(srv.URL()),
    Labels:     labels,                         // < ----------------------- HERE
    LokiConfig: wasp.NewEnvLokiConfig(),
})
```

> [!NOTE]  
> Make sure to set the `BRANCH` and `COMMIT` environment variables in your CI workflow.

---

## `call_group`

This label allows differentiation between different requests within a Virtual User's `Call()` function.

### Example:

```go
// represents user login
func (m *VirtualUser) requestOne(l *wasp.Generator) {
	var result map[string]interface{}
	r, err := m.client.R().
		SetResult(&result).
		Get(m.target)
	if err != nil {
		l.Responses.Err(r, "auth", err) // < ----------------------- HERE
		return
	}
	l.Responses.OK(r, "auth")         // < ----------------------- HERE
}

// represents authenticated user action
func (m *VirtualUser) requestTwo(l *wasp.Generator) {
    var result map[string]interface{}
    r, err := m.client.R().
        SetResult(&result).
        Get(m.target)
    if err != nil {
        l.Responses.Err(r, "user", err) // < ----------------------- HERE
        return
    }
    l.Responses.OK(r, "user")         // < ----------------------- HERE
}
```

In this example, the authentication request and the user action request are differentiated by the `call_group` label.

---

## Custom Labels

> [!NOTE]  
> You can add as many custom labels as needed for use with your dashboard.

```go
labels := map[string]string{
    "custom_label_1": "value1",
    "custom_label_2": "value2",
}
```