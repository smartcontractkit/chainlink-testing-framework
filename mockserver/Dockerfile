FROM golang:alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -o mockserver .

FROM scratch

COPY --from=builder /app/mockserver /mockserver

EXPOSE 8080
ENTRYPOINT [ "/mockserver" ]
