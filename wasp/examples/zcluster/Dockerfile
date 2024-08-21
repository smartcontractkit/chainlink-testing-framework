# Example Dockerfile for k8s run
# Builds all the tests in some directory that must have go.mod
# All tests are built as separate binaries with name "module.test"
FROM golang:1.21 as build
ARG TESTS_ROOT

WORKDIR /go/src
COPY . .

RUN echo $(pwd)
RUN ls -lah
WORKDIR /go/src/${TESTS_ROOT}
RUN echo $(pwd)
RUN ls -lah
RUN cd /go/src/${TESTS_ROOT} && CGO_ENABLED=0 go test -c ./...

FROM debian
ARG TESTS_ROOT

COPY --from=build /go/src/${TESTS_ROOT} .
RUN apt-get update && apt-get install -y ca-certificates
ENTRYPOINT /bin/bash