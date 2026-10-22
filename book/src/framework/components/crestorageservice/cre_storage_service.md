# CRE Storage Service

CRE Storage Service is a composite component that runs [cre-storage-service](https://github.com/smartcontractkit/cre-storage-service) with local dependencies (Postgres, LocalStack S3, a mock account registry). It is a thin `testcontainers-go` wrapper over the compose file that cre-storage-service owns, `local-stack/compose.yaml`; that file's README describes its services, images and inputs.

The stack reads the nodes it authorizes from a CapabilitiesRegistry, so it needs a chain with one deployed, reachable from the stack's Docker network.

## Configuration

```toml
[cre_storage_service]
  # file://<local path>, or an http(s):// URL readable without credentials
  compose_file = "file://../cre-storage-service/local-stack/compose.yaml"
  # existing Docker network to join, defaults to "ctf"
  # network = "ctf"
  registry_rpc_url = "http://anvil-1337:8545"
  registry_chain_selector = 3379446385462418246
  registry_address = "0x5FbDB2315678afecb367f032d93F642f64180aa3"
  don_family = "zone-a"
```

Images default to locally built `:local` tags; `CRE_STORAGE_IMAGE`, `CRE_STORAGE_MIGRATION_IMAGE` and `MOCK_ACCOUNT_REGISTRY_IMAGE` env vars override them, and `LOG_LEVEL` sets the services' log level.

## Endpoints

On the Docker network the services answer to role aliases; `Output` also gives host URLs on free ports. gRPC is served without TLS.

| Output | Role (alias) | Port | Purpose |
| --- | --- | --- | --- |
| `UserService` | `cre-deploy-api` | 50051 | workflow/artifact API for HMAC principals |
| `NodeService` | `cre-workflow-source`, `cre-artifact-api` | 50052 | node API (CSA key auth): workflow metadata, artifact URLs |
| `ArtifactStore` | `cre-artifact-store` | 4566 | S3; artifact URLs use the internal alias |

## Usage

```go
import storage "github.com/smartcontractkit/chainlink-testing-framework/framework/components/dockercompose/cre_storage_service"

out, err := storage.New(in.CREStorageService)
require.NoError(t, err)
// out.UserService.ExternalURL, out.NodeService.InternalURL, ...

// removes every stack started by New, with its volumes
require.NoError(t, storage.Stop())
```
