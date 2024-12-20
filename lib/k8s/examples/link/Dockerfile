# base test for all k8s test runs
FROM golang:1.23-bullseye

ARG BASE_URL
ARG HELM_VERSION
ARG HOME
ARG KUBE_VERSION
ARG NODE_VERSION

ENV GOOS="linux"
ENV BASE_URL="https://get.helm.sh"
ENV HELM_VERSION="3.10.3"
ENV KUBE_VERSION="v1.25.5"
ENV NODE_VERSION=18

RUN apt-get update && apt-get install -y ca-certificates wget curl git gnupg zip && \
    mkdir -p /etc/apt/keyrings && \
    curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && \
    echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_VERSION.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list && \
    apt-get update && apt-get install -y nodejs && \
    curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
    chmod +x ./kubectl && \
    mv ./kubectl /usr/local/bin && \
    case `uname -m` in \
        x86_64) ARCH=amd64; ;; \
        armv7l) ARCH=arm; ;; \
        aarch64) ARCH=arm64; ;; \
        ppc64le) ARCH=ppc64le; ;; \
        s390x) ARCH=s390x; ;; \
        *) echo "un-supported arch, exit ..."; exit 1; ;; \
    esac && \
    wget ${BASE_URL}/helm-v${HELM_VERSION}-linux-${ARCH}.tar.gz -O - | tar -xz && \
    mv linux-${ARCH}/helm /usr/bin/helm && \
    chmod +x /usr/bin/helm && \
    rm -rf linux-${ARCH} && \
    npm install -g yarn && \
    apt-get clean all && \
    helm repo add chainlink-qa https://raw.githubusercontent.com/smartcontractkit/qa-charts/gh-pages/ && \
    helm repo add bitnami https://charts.bitnami.com/bitnami && \
    helm repo update

# This is needed, because kubectl config used is configured to use AWSCLI v2
RUN case `uname -m` in \
        x86_64) AWS_ARCH=x86_64; ;; \
        armv7l) AWS_ARCH=armv7l; ;; \
        aarch64) AWS_ARCH=aarch64; ;; \
        *) echo "un-supported arch, exit ..."; exit 1; ;; \
    esac && \
    curl https://awscli.amazonaws.com/awscli-exe-linux-${AWS_ARCH}.zip -o "awscliv2.zip" && \
    unzip awscliv2.zip && \
    ./aws/install && \
    rm -rf awscliv2.zip

COPY lib/ testdir/
WORKDIR /go/testdir/k8s/examples/link

ENV GOARCH=amd64
ENV GOOS=linux

RUN go test -c . -o link

ENTRYPOINT ["./link"]