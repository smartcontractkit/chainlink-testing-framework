FROM golang:1.24 AS builder

# SSH setup, copy the key you've created for this guide
RUN apt-get update && \
    apt-get install -y openssh-client git && \
    mkdir -p ~/.ssh && \
    chmod 700 ~/.ssh && \
    touch ~/.ssh/known_hosts && \
    chmod 600 ~/.ssh/known_hosts
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts
RUN --mount=type=secret,id=ssh_key \
    cp /run/secrets/ssh_key ~/.ssh/id && \
    chmod 600 ~/.ssh/id && \
    echo "Host github.com\n  IdentityFile ~/.ssh/id\n  IdentitiesOnly yes" > ~/.ssh/config && \
    chmod 600 ~/.ssh/config
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

ENV GOPRIVATE=github.com/smartcontractkit/*

WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=secret,id=ssh_key go mod download
COPY ../.. .
RUN --mount=type=secret,id=ssh_key CGO_ENABLED=0 GOOS=linux go build -o /fake main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /fake /fake
CMD ["/fake"]