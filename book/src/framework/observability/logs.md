# Logs

We are using `Loki` for logging, check [localhost:3000](http://localhost:3000/explore?panes=%7B%22vYC%22:%7B%22datasource%22:%22P8E80F9AEF21F6940%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22expr%22:%22%7Bjob%3D%5C%22ctf%5C%22,%20container%3D~%5C%22node0%5C%22%7D%22,%22queryType%22:%22range%22,%22datasource%22:%7B%22type%22:%22loki%22,%22uid%22:%22P8E80F9AEF21F6940%22%7D,%22editorMode%22:%22code%22%7D%5D,%22range%22:%7B%22from%22:%22now-1h%22,%22to%22:%22now%22%7D%7D%7D&schemaVersion=1&orgId=1)

Queries:
- Particular node logs
```json
{job="ctf", container=~"node0"}
```
- All nodes logs
```json
{job="ctf", container=~"node.*"}
```
- Filter by log level
```json
{job="ctf", container=~"node.*"} |= "WARN|INFO|DEBUG"
```
