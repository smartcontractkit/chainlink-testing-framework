## How to run the same environment deployment inside k8s

<div class="warning">

Managing k8s is challenging, so we've decided to separate `k8s` deployments here - [CRIB](https://github.com/smartcontractkit/crib)

This documentation is outdated, and we are using it only internally to run our soak tests. For `v2` tests please check [this example](../crib.md) and read [CRIB docs](https://github.com/smartcontractkit/crib)
</div>


You can build a `Dockerfile` to run exactly the same environment interactions inside k8s in case you need to run long-running tests
Base image is [here](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/lib/k8s/Dockerfile.base)

```Dockerfile
FROM <account number>.dkr.ecr.us-west-2.amazonaws.com/test-base-image:latest
COPY . .
RUN env GOOS=linux GOARCH=amd64 go build -o test ./examples/remote-test-runner/env.go
RUN chmod +x ./test
ENTRYPOINT ["./test"]
```

Build and upload it using the "latest" tag for the test-base-image

```bash
build_test_image tag=someTag
```

or if you want to specify a test-base-image tag

```bash
build_test_image tag=someTag base_tag=latest
```

Then run it

```bash
# all environment variables with a prefix TEST_ would be provided for k8s job
export TEST_ENV_VAR=myTestVarForAJob
# your image to run as a k8s job
ACCOUNT=$(aws sts get-caller-identity | jq -r .Account)
export ENV_JOB_IMAGE="${ACCOUNT}.dkr.ecr.us-west-2.amazonaws.com/core-integration-tests:v1.1"
export DETACH_RUNNER=true # if you want the test job to run in the background after it has started
export CHAINLINK_ENV_USER=yourUser # user to run the tests
export CHAINLINK_USER_TEAM=yourTeam # team to run the tests for
# your example test file to run inside k8s
# if ENV_JOB_IMAGE is present it will create a job, wait until it finished and get logs
go run examples/remote-test-runner/env.go
```
