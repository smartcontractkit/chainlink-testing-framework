---
layout: default
title: Tear Down
nav_order: 99
parent: Usage
---

## Environments

By default, the `TearDown()` method deletes the environment that was launched by the `SuiteSetup`. Sometimes that's not
desired though, like when debugging failing tests. For that, there's a handy ENV variable, `KEEP_ENVIRONMENTS`.

```sh
KEEP_ENVIRONMENTS = Never # Options: Always, OnFail, Never
```

## Logs

`TearDown()` also checks if the test has failed. If so, it builds a `logs/` directory, and dumps the logs and contents
of each piece of the environment that was launched.
