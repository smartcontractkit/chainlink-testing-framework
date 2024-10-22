# Prysm containers

Docker images provided by Prysm are bare-bone and do not even contain executable shell, which makes debugging impossible. To overcome it we need to build our own image with debug tools installed. The following Dockerfile will do the job for `beacon chain`:

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

# Local Kubernetes @ Docker Desktop (MacOS)

It's very handy and easy to setup, but comes with some hurdles: `hostpath` storage class doesn't work as expected in relation to retention policy. Whether you chose `Delete` or `Recycle` old data won't be removed and your chain will start correctly only the first time, when there is no data yet (consecutive runs will fail, because they will try to generate genesis.json based on previous chain states).

My hacky workaround is to generate a new host directory based on current timestamp every time I start the service. Not idea, but works. It has one main drawback, though: disk gets bloated, since old data is never removed. Now... you can delete it manually, but it's not as straight-fowardward as you might think, because that directory is not directly accessible from your MacOS machine, because it runs inside Docker's VM. Here's how to go about it:

```
docker run -it --rm --privileged --pid=host justincormack/nsenter1
```

(Other options are described [here](https://gist.github.com/BretFisher/5e1a0c7bcca4c735e716abf62afad389))
Then you need to find find a folder where `rootfs` for Docker VM is located, in my case it was `/containers/services/02-docker/rootfs` (search for a folder containing whatever hostpath you have the persistent volume mounted at, or inspect one of your pods, check the volume mount and find it in the VM, just remember that even if your volume is supposedly mounted in `/data/shared` on the host in reality that `/data/shared` folder is still relative to wherever `rootfs` folder is located).

Once you find it, you can delete the old data and start the service again. Or you can use it for debugging to easily inspect the content of shared volumes (useful in case of containers that are bare-bone and don't even have `bash` installed, like Prysm).
