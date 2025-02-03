# Asserting Container Logs

You can either assert that CL nodes have no errors like that, we check `(CRIT|PANIC|FATAL)` levels by default for all the nodes

```golang
	in, err := framework.Load[Cfg](t)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := framework.SaveAndCheckLogs(t)
		require.NoError(t, err)
	})
```

or customize file assertions

```golang
	in, err := framework.Load[Cfg](t)
	require.NoError(t, err)
	t.Cleanup(func() {
		// save all the logs to default directory "logs/docker-$test_name"
		logs, err := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, err)
		// check that CL nodes has no errors (CRIT|PANIC|FATAL) levels
		err = framework.CheckCLNodeContainerErrors()
		require.NoError(t, err)
		// do custom assertions
		for _, l := range logs {
			matches, err := framework.SearchLogFile(l, " name=HeadReporter version=\\d")
			require.NoError(t, err)
			_ = matches
		}
	})
```

Full [example](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/framework/examples/myproject/smoke_logs_test.go)