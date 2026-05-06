FROM golang:1.25.3 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG CTF_LOG_LEVEL=info
ENV CTF_LOG_LEVEL=${CTF_LOG_LEVEL}

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags="-s -w" -o /out/chip-router ./cmd/chip-router

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/chip-router /chip-router

EXPOSE 50050 50051

ENTRYPOINT ["/chip-router"]
