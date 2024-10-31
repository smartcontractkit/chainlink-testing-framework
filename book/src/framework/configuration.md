# Configuration

### Environment variables
|             Name             |                                                                      Description                                                                       |          Possible values | Default |        Required?         |
|:----------------------------:|:------------------------------------------------------------------------------------------------------------------------------------------------------:|-------------------------:|:-------:|:------------------------:|
|         CTF_CONFIGS          | Path(s) to test config files. <br/>Can be more than one, ex.: smoke.toml,smoke_1.toml,smoke_2.toml.<br/>First filepath will hold all the merged values | Any valid TOML file path |         |            ✅             |
|        CTF_LOG_LEVEL         |                                                                   Harness log level                                                                    | `info`, `debug`, `trace` | `info`  |            🚫            |
|       CTF_LOKI_STREAM        |                                                Streams all components logs to `Loki`, see params below                                                 |          `true`, `false` | `false` |            🚫            |
|           LOKI_URL           |                                            URL to `Loki` push api, should be like`${host}/loki/api/v1/push`                                            |                      URL |    -    | If you use `Loki` then ✅ |
|        LOKI_TENANT_ID        |                                                Streams all components logs to `Loki`, see params below                                                 |          `true`, `false` |    -    | If you use `Loki` then ✅ |
| TESTCONTAINERS_RYUK_DISABLED |                                   Testcontainers-Go reaper container, removes all the containers after the test exit                                   |          `true`, `false` | `false` |            🚫            |
|         RESTY_DEBUG          |                                                            Log all Resty client HTTP calls                                                             |          `true`, `false` | `false` |            🚫            |
