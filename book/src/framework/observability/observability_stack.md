# Local Observability Stack

You can use a local observability stack, framework is connected to it by default

```bash
ctf obs up
```

To remove it use

```bash
ctf obs down
```

Read more about how to check [logs](logs.md) and [profiles](profiling.md)

## Developing

Change compose files under `framework/cmd/observability` and restart the stack (removing volumes too)
```
just reload-cli && ctf obs r
```

## Local Dashboards (Docker)

You can create a dashboard using [UI](http://localhost:3000) and put them under `$pwd/dashboards` folder then commit, they'll be loaded automatically on start and you can find them [here](http://localhost:3000/dashboards) under `local` directory.

`$pwd` is you current working directory from which you call `ctf obs u`

## Troubleshooting

### `cadvisor` is not working
Make sure your `Advanced` Docker settings look like this
![img_2.png](img_2.png)

