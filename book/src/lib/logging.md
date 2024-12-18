# Logging

This small library was created to address two issues:
* mixed up logging for parallel tests, when using vanilla loggers
* conformity with logging interface required by `testcontainers-go` (a Docker container library)

It uses `"github.com/rs/zerolog"` for the logger.

## Configuration
There's only one configuration option: the log level. You can set it via `TEST_LOG_LEVEL` environment variable to:
* `trace`
* `debug`
* `info` (default)
* `warn`
* `error`

## How to use
The main way to get a Logger instance is to call `logging.GetTestLogger(*testing.T)`. `testing.T` instance can be `nil`.

When using it together with `testcontainers-go`, which is a library we use to interact with Docker containers you should
use `GetTestContainersGoTestLogger(*testing.T)` instead.

And that's all there is to it :-)