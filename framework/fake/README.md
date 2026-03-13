## CRE Fakes

## Prerequisites

```bash
brew install just
```

## Deploying

```bash
# local build
just build
# local run
just run
# push to remote ECR
just push <main.stage_registry> <any_tag>
# deploy
just deploy <main.stage_registry> <any_tag> <namespace_with_cre_services>
```
