# `config` Package

Fairly self explanatory.

## `charts.json`

A convenient helper file that you can use to squash chart values, namely the chainlink image and version. This is usually more convenient than having to set the `CHARTS` environment variable to the raw JSON.

```sh
CHARTS='../../config/charts.json' make test_smoke
```
