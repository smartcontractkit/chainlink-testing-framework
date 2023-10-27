ARG BASE_IMAGE
ARG IMAGE_VERSION=latest
FROM ${BASE_IMAGE}:${IMAGE_VERSION}
COPY . testdir/
WORKDIR /go/testdir
RUN ./k8s/scripts/buildTests
ENTRYPOINT ["./k8s/scripts/entrypoint"]
