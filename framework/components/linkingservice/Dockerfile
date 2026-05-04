FROM golang:1.25.3 AS build

WORKDIR /src/framework

COPY framework/go.mod framework/go.sum ./
COPY wasp/go.mod wasp/go.sum /src/wasp/
RUN go mod download

COPY framework /src/framework
COPY wasp /src/wasp
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/linking-service ./components/linkingservice/cmd/linking-service

FROM gcr.io/distroless/base-debian12

COPY --from=build /out/linking-service /linking-service
EXPOSE 18124 18125
ENTRYPOINT ["/linking-service"]
