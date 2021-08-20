FROM golang:alpine as builder

RUN mkdir /app
WORKDIR /app
RUN apk add --update --no-cache ca-certificates git
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /bin/explorer cmd/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/explorer /bin/explorer
ENTRYPOINT ["/bin/explorer"]