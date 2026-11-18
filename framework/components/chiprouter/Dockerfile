FROM golang:1.26.6 AS builder

# go.mod requires a published `framework` release (no `replace`), so the build
# needs only the checked-in go.mod/go.sum: no go.mod editing and no
# `go mod tidy` at build time. The router binary itself does not import
# `framework` (only the deploy-side library does).
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG CTF_LOG_LEVEL=info
ENV CTF_LOG_LEVEL=${CTF_LOG_LEVEL}

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -mod=readonly -trimpath -ldflags="-s -w" -o /out/chip-router ./cmd/chip-router

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/chip-router /chip-router

EXPOSE 50050 50051

ENTRYPOINT ["/chip-router"]
