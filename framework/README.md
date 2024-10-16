## Chainlink Testing Framework

This is a common component to stream all your local docker logs to `k8s`

### Streaming logs to Loki
Add labels to all your local docker. Here is an example you can use to test component, do that in component `Golang` code for framework components
```
docker run --name ctf_nginx --label=logging=promtail nginx
```

Add environment variables or secrets in CI
```
export LOKI_STREAM=true
export LOKI_TENANT_ID=
export LOKI_URL=           # should point ot v1/push API of Loki
export LOKI_BASIC_AUTH=
```