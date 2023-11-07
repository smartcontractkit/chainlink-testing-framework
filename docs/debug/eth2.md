# Prysm
Docker images provided by Prysm are bare-bone and do not even contain exectuable shell, which makes debugging impossible. To overcome it we need to build our own image with debug tools installed. The following Dockerfile will do the job for `beacon chain`:
```
ARG tag
FROM --platform=linux/x86_64 gcr.io/prysmaticlabs/prysm/beacon-chain:$tag as upstream

FROM --platform=linux/x86_64 debian:buster-slim
COPY --from=upstream /app /app

RUN apt-get update && apt-get install -y \
    curl \
    jq \
    nano \
    netcat-openbsd \
    iputils-ping \
    && rm -rf /var/lib/apt/lists/*


STOPSIGNAL SIGINT

ENTRYPOINT ["/app/cmd/beacon-chain/beacon-chain"]
```

And for `validator`:
```
ARG tag
FROM --platform=linux/x86_64 gcr.io/prysmaticlabs/prysm/validator:$tag as upstream

FROM --platform=linux/x86_64 debian:buster-slim
COPY --from=upstream /app /app

RUN apt-get update && apt-get install -y \
    curl \
    jq \
    nano \
    netcat-openbsd \
    iputils-ping \
    && rm -rf /var/lib/apt/lists/*


STOPSIGNAL SIGINT

ENTRYPOINT ["/app/cmd/validator/validator"]
```

And the use as follows:
```
docker build -t gcr.io/prysmaticlabs/prysm/validator:debug --build-arg tag=v4.1.0 .
```

# Lighthouse
No supported yet.