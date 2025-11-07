# K8s Deployment (Deprecated)

[![Documentation](https://img.shields.io/badge/Documentation-MDBook-blue?style=for-the-badge)](https://smartcontractkit.github.io/chainlink-testing-framework/lib/k8s/KUBERNETES.html)

## K8s Framework Sunsetting

We were using this `K8s` framework specifically in `sdlc` cluster but after migration to `main.stage` we use it only on-demand for some legacy tests.

### Building Base Image for K8s Tests
This is rarely required, but sometimes infrequent updates are needed.

Go to SSO [home](https://sso.smartcontract.com/app/UserHome) and find `AWS -> secure-sdlc`, copy the registry ID
```
<sdlc_registry_id> | secure-sdlc@smartcontract.com
```
use creds for this account in your terminal, paste the <sdlc_registry_id> then build and push the image.

```bash
make build
make push registry=<sdlc_registry_id>
```

