FROM golang:1.26.6 AS builder

# cmd/chip-router is its own Go module (go.mod/go.sum are checked in next to
# main.go), so the image builds from pinned dependencies only: no `replace`
# editing and no `go mod tidy` at build time. The router binary does not import
# the parent chiprouter module (the deploy-side library) or `framework`.
WORKDIR /src/cmd/chip-router

COPY cmd/chip-router/go.mod cmd/chip-router/go.sum ./
RUN go mod download

COPY cmd/chip-router/ ./

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG CTF_LOG_LEVEL=info
ENV CTF_LOG_LEVEL=${CTF_LOG_LEVEL}

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -mod=readonly -trimpath -ldflags="-s -w" -o /out/chip-router .

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/chip-router /chip-router

EXPOSE 50050 50051

ENTRYPOINT ["/chip-router"]
