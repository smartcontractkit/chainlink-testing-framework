# Chainlink ecosystem

Currently, only one application has a wrapper that lives in this framework: the Job Distributor.
Chainlink Node wrapper can be found in the [chainlink repository](https://github.com/smartcontractkit/chainlink/blob/develop/integration-tests/docker/test_env/cl_node.go).

## Job Distributor
JD is a component for a centralized creation and management of jobs executed by Chainlink Nodes. It's a single point of entry
that frees you from having to setup each job separately on each node from the DON.

It requires a Postgres DB instance, which can also be started using the CTF:
```go
pg, err := test_env.NewPostgresDb(
    []string{network.Name},
    test_env.WithPostgresDbName("jd-db"),
    test_env.WithPostgresImageVersion("14.1"))
if err != nil {
    panic(err)
}
err = pg.StartContainer()
if err != nil {
    panic(err)
}
```

Then all you need to do, to start a new instance is:
```go
jd := New([]string{network.Name},
    //replace with actual image
    WithImage("localhost:5001/jd"),
    WithVersion("latest"),
    WithDBURL(pg.InternalURL.String()),
)

err = jd.StartContainer()
if err != nil {
    panic(err)
}
```

Once you have JD started you can create a new GRPC connection and start interacting with it
to register DON nodes, create jobs, etc:
```go
conn, cErr := grpc.NewClient(jd.Grpc, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}...)
if cErr != nil {
    panic(cErr)
}

// use the connection
```