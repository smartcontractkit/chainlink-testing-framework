FROM golang:1.25.3 AS build

WORKDIR /src/framework

COPY framework/go.mod framework/go.sum ./
COPY wasp/go.mod wasp/go.sum /src/wasp/
RUN go mod download

COPY framework /src/framework
COPY wasp /src/wasp
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/vault-jwt-issuer ./components/vaultjwtissuer/cmd/vault-jwt-issuer

FROM gcr.io/distroless/base-debian12

COPY --from=build /out/vault-jwt-issuer /vault-jwt-issuer
EXPOSE 18123
ENTRYPOINT ["/vault-jwt-issuer"]
