FROM golang:alpine AS builder

ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -o mockserver .

FROM scratch

COPY --from=builder /app/mockserver /mockserver

EXPOSE 8080
ENTRYPOINT [ "/mockserver" ]
