## How to run the same environment deployment inside k8s

You can build a `Dockerfile` to run exactly the same environment interactions inside k8s in case you need to run long-running tests
Base image is [here](Dockerfile.base)
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
# your example test file to run inside k8s
# if ENV_JOB_IMAGE is present chainlink-env will create a job, wait until it finished and get logs
go run examples/remote-test-runner/env.go
```