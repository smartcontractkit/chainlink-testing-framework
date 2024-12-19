# Kubernetes - Using remote runner

In this chapter we will explain how you can run a test in `k8s` and what changes to the test logic are needed.

First, some introduction. The general logic of runnnig these tests is as follows:
* create a `k8s` environment from the local machine (local meaining both your local machine or a CI runner)
    * that environment will launch a `remote runner` (thanks to the `ENV_JOB_IMAGE` environment variable)
* `remote runner` will execute **the same testing code again, from the beginning**
    * while being smart enough to notice that environment is already deployed, so it doesn't deploy it again
* once `remote runner` has finished to run it's copy of the test it will return control to local test execution
    * local test will exit early to avoid running the test logic again
* if running in `detached mode` control will be returned to the local test as soon as remote test has started
    * local test will exit as soon as control is returned
    * `remote runner` will keep running in `k8s` until the test has finished


I admit it does sound a bit complex as it requires some changes to the test logic, to make sure that it doesn't
really execute twice. In the next steps we will see where and why we should add certain ifs and early exists
to prevent that from happening.

> [!NOTE]
> We are working on a new `k8s` test runner, which doesn't require adding any `k8s`-specific logic and allows
> to run a the same test both remotely and locally, but it's not ready yet.
>
> You can find some of its documentation [here](../../k8s-test-runner/k8s-test-runner.md), though.

# Requirements
In general, `remote runner` requires a Docker image with your test. There are various ways to build it, some of them
already automated in our repositories, but for the sake of this documentation we will do everything from the scratch.

> [!NOTE]
> Current approach is that the CTF repository builds a base testing image for each CTF release using [this action](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/k8s-publish-test-base-image.yaml).
> That image contains `kubectl`, `helm` and a couple of other dependencies, so that you can use it as a base image, when building
> your final image. It means that all you need to do to get a final image is copy your tests compiled as a Go binary and
> instruct the entrypoint to execute them.
> An example of a Chainlink repo that does it, can be found [here](https://github.com/smartcontractkit/chainlink/actions/workflows/on-demand-ocr-soak-test.yml).

# Step 1: Build a Docker image with your tests
Let's define a Dockerfile first:
```docker
# base test for all k8s test runs
FROM golang:1.23-bullseye

ARG GOARCH
ARG GOOS
ARG BASE_URL
ARG HELM_VERSION
ARG HOME
ARG KUBE_VERSION
ARG NODE_VERSION

# compile Go binary targetting linux/amd64, since that's what our k8s runners use
ENV GOOS="linux"
ENV GOARCH=amd64

ENV BASE_URL="https://get.helm.sh"
ENV HELM_VERSION="3.10.3"
ENV KUBE_VERSION="v1.25.5"
ENV NODE_VERSION=18

# install all dependencies
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

# kubectl config is configured to use AWSCLI v2, so we need to install it
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

# compile specific test folder
WORKDIR /go/testdir/k8s/examples/link
RUN go test -c . -o link

# run compiled test
ENTRYPOINT ["./link"]
```

As mentioned before, it installs requied dependencies, such as:
* `kubectl`
* `nodejs`
* `helm`
* `AWS CLI v2`

`kubectl` is strictly necessary as we use it for interacting with `k8s` for creation of test environments and the `remote runner` itself. `Helm` runtime and
addition of repositories is required, because our environment is using Helm charts created by us and published in our Helm repository.

In our very case our `kubectl` is configured to use `AWS CLI` to communicate with AWS `k8s` cluster and thus we will need that dependency as well.
You might need to adjust it to reflect your config, if you're not using AWS.

> [!NOTE]
> If you'd like to use Helm charts committed in this repository, and included in the image, instead of published ones you should set `LOCAL_CHARTS=true`.

# Step 2: Build the image
Now it's time to build the image. One important thing to keep in mind here is to build it from the root folder of the CTF repository (as it needs to have access
to the `lib` subfolder in order to copy it; in other words: that folder needs to be present in the default build context).

```bash
docker build -f lib/k8s/examples/link/Dockerfile \
    --platform=linux/amd64 \
    -t link-test:latest .
```

> [!NOTE]
> Notice the `--platform=linux/amd64` parameter, which is necessary if you aren't building the image on a linux machine with x86 CPU architecture.

# Step 3: Test the image
Since we haven't modified the code yet to support `remote runner` execution, we can only run it in local mode. Let's do that to make sure that everything works as expected.

For that you need to have access to a `k8s` cluster. It doesn't matter whether it's a local one (e.g. running on Docker or `k3d`) or a remote one, running in AWS or Azure or Google Cloud.
All that matters is that `kubectl` on your local machine is configured to use that cluster and that you are logged in, if your service requires authorization. That's because we will
share your host machine's `kubectl` configuration with the container.

Now that you have `kubectl` configured and you are logged in to the AWS (since that's what we are using the example), lets run our tests using the image we just built:
```bash
docker run \
    --rm \
    -v ~/.aws:/root/.aws:ro \
    -v ~/.kube/config:/root/.kube/config:ro \
    -e AWS_PROFILE=<your-prfile> \
    -e KUBECONFIG=/root/.kube/config \
    link-test:latest
```
We have told the Docker to mount your local `.aws` directory and your `kubectl` config directory `.kube/config` in the `/root` folder of the Docker container, because contrary to best practices
this container is running as `root`. **This is not how a production image should run!**.

We did mount these directories so that the container can use your host machine's `kubectl` configuration and your active AWS session.

The test should run inside the Docker container and finish after a couple of minutes. If it did, we know that dockerisation did work. Time to make the final changes that will allows us to run it detached mode
using `remote runner`.

# Step 4: Make the test `remote runner`-compatible
As mentioned in the beginning, now we need to divide the test logic into two parts:
* one that will run both on local and remote runner
* one that will run only in the remote runner

First part should be the idempotent part, meaning that we can run it multiple times without any side effects. In most cases that will be the creation of environment. It needs to always run on the local machine,
if only because it creates the `remote runner` job in the `k8s`. The runner is smart enough to not create the environment again.

Once the environment has been created, and `remote runner` has started the local test should exit (after waiting or not, for the remote to finish, depending on the `DETACH_RUNNER` environment variable).

In our simple case, the code that creates the environment looked like this:
```go
err = testEnv.Run()
if err != nil {
    t.Fatal("Error running environment: ", err)
}

log.Info().
    Bool("Remote runner?", testEnv.WillUseRemoteRunner()).
    Msg("Started environment")

// if test is running inside K8s, nothing to do, default network urls are correct
if !testEnv.Cfg.InsideK8s {
    // rewrite the urls
    // ...
```

That means we want the local version to exit after environment creation, and we will achieve that in the following way:
```go
err = testEnv.Run()
if err != nil {
    t.Fatal("Error running environment: ", err)
}

// ===================== NEW CODE START ========================
if testEnv.WillUseRemoteRunner() {
    log.Info().
        Msg("Stopping local execution as test will continu in the remote runner")
    return
}
// ======================= NEW CODE END ========================

log.Info().
    Bool("Remote runner?", testEnv.WillUseRemoteRunner()).
    Msg("Started environment")

// if test is running inside K8s, nothing to do, default network urls are correct
if !testEnv.Cfg.InsideK8s {
    // rewrite the urls
```

So basically we need to use `testEnv.WillUseRemoteRunner()` function to check whether this test will also execute in remote runner and in the case it does stop execution of code
that should run only once (or is not idempotent). Usually that means:
* after environment setup
* before rewritting URLs to use forwarded versions
* before test clean up (we want to execute it once the `remote runner` has finished, not the local one)

If you go back to the test we have written previously you will see that the part that should execute only in the `remote runner` is the part that deploys contracts and then maybe
goes on to generate some load. Here we will enhance it, so that it deploys 5 contracts instead of 1 and sleeps between each attempt, so that you have enough time to observe detached
remote mode in action :-)

```go
sethClient, err := seth.NewClientBuilder().
		WithRpcUrl(nodeNetwork.URLs[0]).
		WithPrivateKeys([]string{nodeNetwork.PrivateKeys[0]}).
		Build()
if err != nil {
    t.Fatal("Error creating Seth client", err)
}

for i := 0; i < 5; i++ {
    log.Info().
        Msgf("Deploying LinkToken contract, instance %d/%d", i+1, 5)

    linkTokenAbi, err := link_token_interface.LinkTokenMetaData.GetAbi()
    if err != nil {
        t.Fatal("Error getting LinkToken ABI", err)
    }
    linkDeploymentData, err := sethClient.DeployContract(sethClient.NewTXOpts(), "LinkToken", *linkTokenAbi, common.FromHex(link_token_interface.LinkTokenMetaData.Bin))
    if err != nil {
        t.Fatal("Error deploying LinkToken contract", err)
    }
    linkToken, err := link_token_interface.NewLinkToken(linkDeploymentData.Address, sethClient.Client)
    if err != nil {
        t.Fatal("Error creating LinkToken contract instance", err)
    }

    totalSupply, err := linkToken.TotalSupply(sethClient.NewCallOpts())
    if err != nil {
        t.Fatal("Error getting total supply of LinkToken", err)
    }

    if totalSupply.Cmp(big.NewInt(0)) <= 0 {
        t.Fatal("Total supply of LinkToken should be greater than 0")
    }

    time.Sleep(15 * time.Second)
}
```

# Step 5: Push image to image registry (optional)
In case you won't be using a local `k8s` cluster you should push your image to an image registry that's available to your cluster.

We will use AWS ECR for that. First, we will login:
```bash
aws --profile <YOUR_PROFILE> ecr get-login-password \
    --region <AWS_REGION> | docker login \
    --username AWS \
    --password-stdin <AWS_ACCOUNT>.dkr.ecr.<AWS_REGION>.amazonaws.com
```

Then tag our image:
```bash
docker tag link-test:latest \
    <AWS_ACCOUNT>.dkr.ecr.<AWS_REGION>.amazonaws.com/<AWS_REPOSITORY>/link-remote-runner-test:latest
```

And finally push it:
```bash
docker push \
    <AWS_ACCOUNT>.dkr.ecr.<AWS_REGION>.amazonaws.com/<AWS_REPOSITORY>/link-remote-runner-test:latest
```

# Step 6: Run the test with `remote runner`
Now it's time to put everything together and run our test with `remote runner` in the Kubernetes. We will run it in `detached` mode, which means that local execution will detach itself
from the process as soon as the remote test has started. This will allow us to continue with our work while the test is running.

> [!NOTE]
> When using detached mode the namespace is never deleted at the end of the test. Instead it keeps on running so that you can access the `remote runner` and check test logs to make sure
> whether everything went as expected.
>
> Of course, you could manage test results in some other way, for example by storing test logs on external storage or notifying you about the result and then removing the namespace.

In order to do that we will need to set two environment variables:
* `DETACH_RUNNER` - set to `true` to detach as soon as environment is up and running
* `ENV_JOB_IMAGE` - Docker image to use for `remote runner`

> [!NOTE]
> [Here](test_secrets.md) you can read how to pass secrets to your test.

Now our Docker command becomes:
```bash
docker run \
    --rm \
    -v ~/.aws:/root/.aws:ro \
    -v ~/.kube/config:/root/.kube/config:ro \
    -e DETACH_RUNNER=true \
    -e ENV_JOB_NAME="<AWS_ACCOUNT>.dkr.ecr.<AWS_REGION>.amazonaws.com/<AWS_REPOSITORY>/link-remote-runner-test:latest" \
    -e AWS_PROFILE=<your-prfile> \
    -e KUBECONFIG=/root/.kube/config \
    <AWS_ACCOUNT>.dkr.ecr.<AWS_REGION>.amazonaws.com/<AWS_REPOSITORY>/link-remote-runner-test:latest
```

If everything went well, you should see something along these lines at the bottom of the test output:
```bash
10:55:26.60 INF Waiting for remote runner to complete
10:55:26.60 INF Started environment Remote runner?=true
10:55:26.60 INF Exiting as test will use remote runner
PASS
```

That's it. It's been quite a complex journey, but hopefully now you know how to write a `k8s` test that can be executed both from a local machine
and with a detached `remote runner`.

> [!NOTE]
> You can find this example [here](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/lib/k8s/examples/link). It includes both
> the test code and the `Dockerfile`.