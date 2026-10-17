FROM golang:1.26.6 AS builder

WORKDIR /src

COPY . .

# The module's go.mod replaces the parent framework module with ../../ (a
# monorepo-dev convenience) which cannot resolve inside this build's
# module-dir context, and the published framework does not compile against
# this module's moby pins. The router binary no longer imports framework
# (only the deploy-side library does, and it is not compiled into the
# image), so drop the replace and fetch the published framework module from
# the proxy for module-graph completeness only.
RUN sed -i '/^replace /d' go.mod && go mod tidy && go mod download

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG CTF_LOG_LEVEL=info
ENV CTF_LOG_LEVEL=${CTF_LOG_LEVEL}

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags="-s -w" -o /out/chip-router ./cmd/chip-router

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/chip-router /chip-router

EXPOSE 50050 50051

ENTRYPOINT ["/chip-router"]
