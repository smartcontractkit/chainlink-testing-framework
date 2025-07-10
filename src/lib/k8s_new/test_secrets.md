# Kubernetes - Test Secrets

We all have our secrets, don't we? It's the same case with tests... and since some of our repositories are public, we need to take special precautions to protect them.

> [!WARNING]
> Before continuing, you should read the [test secrets section of the CTF configuration documentation](../config/config.md).

## Overview

In general, your `remote runner` will need access to the same secrets as your local test. Fortunately, these secrets are forwarded automatically and securely as long as their names have the prefix `E2E_TEST_`.

To make a secret available to the `remote runner`, simply pass it to the `docker run` command:

```bash
docker run \
    --rm \
    -v ~/.aws:/root/.aws:ro \
    -v ~/.kube/config:/root/.kube/config:ro \
    -e DETACH_RUNNER=true \
    -e E2E_TEST_MY_SECRET=my-secret \
    -e ENV_JOB_NAME="<image-url>" \
    -e AWS_PROFILE=<your-profile> \
    -e KUBECONFIG=/root/.kube/config \
    <image-url>
```

The secret will then be available to the `remote runner` during its execution.

---

## Important Considerations

> [!WARNING]
> **Do not use this method of passing secrets in CI environments.** Exposing secrets in this way can compromise their security.
>
> When running `k8s` tests in CI pipelines, use dedicated actions or reusable workflows designed to handle secrets securely.