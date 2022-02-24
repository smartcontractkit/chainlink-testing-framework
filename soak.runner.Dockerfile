FROM golang:alpine as builder

# Setup app folder
RUN mkdir /app
ADD . /app
WORKDIR /app

# Build soak tests
RUN go mod download && \
    CGO_ENABLED=0 go test -c ./suite/soak/tests -o soak.test && \
    chmod +x ./soak.test
  
FROM scratch
COPY --from=builder /app/soak.test /app/soak.test
COPY --from=builder /app/framework.yaml /app/framework.yaml
COPY --from=builder /app/networks.yaml /app/networks.yaml
