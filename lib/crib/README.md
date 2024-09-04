### CRIB Connector

This is a simple CRIB connector for OCRv1 CRIB
This code is temporary and may be removed in the future if connection logic will be simplified with [ARC](https://github.com/actions/actions-runner-controller)

## Example

Go to the [CRIB](https://github.com/smartcontractkit/crib) repository and spin up a cluster.

```shell
./scripts/cribbit.sh crib-oh-my-crib
devspace deploy --debug --profile local-dev-simulated-core-ocr1
```

## Run an example test

```shell
export CRIB_NAMESPACE=crib-oh-my-crib
export CRIB_NETWORK=geth # only "geth" is supported for now
export CRIB_NODES=5 # min 5 nodes
#export SETH_LOG_LEVEL=debug # these two can be enabled to debug connection issues
#export RESTY_DEBUG=true
export GAP_URL=https://localhost:8080/primary # only applicable in CI, unset the var to connect locally
go test -v -run TestCRIB
```
