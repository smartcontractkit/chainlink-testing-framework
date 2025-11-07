# K8s Deployment (Deprecated)

[![Documentation](https://img.shields.io/badge/Documentation-MDBook-blue?style=for-the-badge)](https://smartcontractkit.github.io/chainlink-testing-framework/lib/k8s/KUBERNETES.html)

## K8s Framework Sunsetting

We were using this `K8s` framework specifically in `sdlc` cluster but after migration to `main.stage` we use it only on-demand for some legacy tests.

### Building Base Image for K8s Tests
Is rarely required but sometimes we need infrequent updates.

Go to SSO [home](https://sso.smartcontract.com/app/UserHome) and find AWS tile, use creds for `secure-sdlc` account in your terminal, then build and push the image.

You can find `sdlc` registry id in AWS UI.

```bash
make build
make push registry=<sdlc_registry_id>
```

