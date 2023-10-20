ARG BASE_IMAGE
ARG IMAGE_VERSION=latest
FROM ${BASE_IMAGE}:${IMAGE_VERSION}
COPY . testdir/
WORKDIR /go/testdir
RUN ./scripts/buildTests
ENTRYPOINT ["./scripts/entrypoint"]
