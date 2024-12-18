# Postgres Connector

Postgres Connect simplifies connecting to a Postgres DB, by either:
* providing it full db config
* providing it a pointer to `enviroment.Environment`

# Arbitrary DB

```go
config := &PostgresConfig{
    Host: "db-host.io",
    Port: 5432,
    User: "admin",
    Password: "oh-so-secret",
    DBName: "database",
    //SSLMode:
}
connector, err := NewPostgresConnector(config)
if err != nil {
    panic(err)
}
```

If no `SSLMode` is supplied it will default to `sslmode=disable`.

# K8s
This code assumes it connects to k8s enviroment created with the CTF, where each Chainlink Node
has an underlaying Postgres DB instance using CTF's default configuration:
```go
var myk8sEnv *environment.Environment

// ... create k8s environment

node0PgClient, node0err := ConnectDB(0, myk8sEnv)
if node0err != nil {
    panic(node0err)
}

node1PgClient, node1err := ConnectDB(1, myk8sEnv)
if node1err != nil {
    panic(node1err)
}
```

