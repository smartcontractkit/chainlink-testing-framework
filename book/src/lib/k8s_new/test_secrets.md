# Kubernetes - Test secrets

We all have our secrets, don't we? It's the same case with tests... and since some of our repositories are public
we need to take special precautions in protecting them.

> [!WARNING]
> Before continuing you should read [test secrets part of the CTF config documentation](../config/config.md).

In general, your `remote runner` will need to have access to the same secrets as your local test would. Forunatelly,
they are forwarded automatically and securily as long as their name has prefix `E2E_TEST_`.

So all that you need to do is to pass them to `docker run` command:
```bash
docker run \
    --rm \
    -v ~/.aws:/root/.aws:ro \
    -v ~/.kube/config:/root/.kube/config:ro \
    -e DETACH_RUNNER=true \
    -e E2E_TEST_MY_SECRET=my-secret \
    -e ENV_JOB_NAME="<AWS_ACCOUNT>.dkr.ecr.<AWS_REGION>.amazonaws.com/<AWS_REPOSITORY>/link-remote-runner-test:latest" \
    -e AWS_PROFILE=<your-prfile> \
    -e KUBECONFIG=/root/.kube/config \
    <AWS_ACCOUNT>.dkr.ecr.<AWS_REGION>.amazonaws.com/<AWS_REPOSITORY>/link-remote-runner-test:latest
```
And it will be available in the remore runner.

> [!WARNING]
> **This way of passing the secrets should by no means be used in the CI, as it will most surely result in their exposure**.
> Running `k8s` tests in the CI is completely out of the scope of this documentation. You should use dedicated actions
> and reusable workflows.