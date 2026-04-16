# Asserting Container Logs

Use built-in critical-level assertion (`CRIT|PANIC|FATAL`) for Chainlink node logs:

```golang
in, err := framework.Load[Cfg](t)
require.NoError(t, err)
t.Cleanup(func() {
	err := framework.SaveAndCheckLogs(t)
	require.NoError(t, err)
})
```

For custom checks, assert logs directly from streams with `StreamCTFContainerLogsFanout`.

```golang
re := regexp.MustCompile(`name=HeadReporter version=\d+`)
t.Cleanup(func() {
	err := framework.StreamCTFContainerLogsFanout(
		framework.LogStreamConsumer{
			Name: "custom-regex-assert",
			Consume: func(logStreams map[string]io.ReadCloser) error {
				for name, stream := range logStreams {
					scanner := bufio.NewScanner(stream)
					found := false
					for scanner.Scan() {
						if re.MatchString(scanner.Text()) {
							found = true
							break
						}
					}
					if err := scanner.Err(); err != nil {
						return fmt.Errorf("scan %s: %w", name, err)
					}
					if !found {
						return fmt.Errorf("missing HeadReporter log in %s", name)
					}
				}
				return nil
			},
		},
	)
	require.NoError(t, err)
})
```

Full [example](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/examples/myproject/smoke_logs_test.go)