## Observability tools
We have some observability tools we use with our harness, you can use them by calling
```
ctf observability up
```
Change your `Loki` config in your `.envrc` you use to run tests
```
export LOKI_TENANT_ID=promtail
export LOKI_URL=http://host.docker.internal:3030/loki/api/v1/push
```
Then check [Loki](http://localhost:3000/explore?panes=%7B%220EE%22:%7B%22datasource%22:%22P8E80F9AEF21F6940%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22expr%22:%22%7Bjob%3D%5C%22ctf%5C%22%7D%22,%22queryType%22:%22range%22,%22datasource%22:%7B%22type%22:%22loki%22,%22uid%22:%22P8E80F9AEF21F6940%22%7D,%22editorMode%22:%22code%22%7D%5D,%22range%22:%7B%22from%22:%22now-5m%22,%22to%22:%22now%22%7D%7D%7D&schemaVersion=1&orgId=1) logs