# Compatibility Testing

## Prerequisites

Authorize in our SDLC ECR registry first. Get the creds and run

```bash
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin <sdlc_ecr_registry>
```

## Testing Upgrade Sequence

We have a simple tool to check compatibility for CL node clusters. The example command will filter and sort the available tags, rollback and install the oldest version, and then begin performing automatic upgrades to verify that each subsequent version remains compatible with the previous one.

`buildcmd`, `envcmd`, `testcmd` can be arbitrary bash commands.

```bash
ctf compat backward \
--registry <sdlc_ecr_registry> \
--buildcmd "just cli" \
--envcmd "cl r" \
--testcmd "cl test ocr2 TestSmoke/rounds" \
--refs 2.32.0 \
--refs 2.33.0 \
--refs 2.34.0 \
--refs 2.35.0 \
--upgrade-nodes 3
```

Keep in mind that `refs` should be present in regsitry you are testing against, the first (oldest) `ref` should also have a valid end-to-end test that works.

In CI we detect SemVer tags automatically, whenever a new tag appears we select last 3, rollback to the oldest and perform upgrade process.

```bash
ctf compat backward \
--registry <sdlc_ecr_chainlink_registry> \
--buildcmd "just cli" \
--envcmd "cl r env.toml,products/ocr2/basic.toml" \
--testcmd "cl test ocr2 TestSmoke/rounds" \
--strip-image-suffix v \
--upgrade-nodes 2 \
--versions-back 2
```

In case you have multiple DONs in your product and names of nodes are different please use `--node-name-template custom-cl-node-%d` option

## Modelling Node Operators Cluster (WIP)

It is possible to fetch versions node operators are currently running and model DON upgrade sequence locally. When `product` is specified, `compat` will fetch the current versions from the RANE SOT data source and model the upgrade sequence for versions found on real DONs up to the latest one, each node one at a time.

```bash
ctf compat backward \
--registry <sdlc_ecr_chainlink_registry> \
--buildcmd "just cli" \
--envcmd "cl r env.toml,products/ocr2/basic.toml" \
--testcmd "cl test ocr2 TestSmoke/rounds" \
--product data-feeds \
--no-git-rollback \
--don_nodes 5
```

The tool will check out earliest Git `ref` and setup environment and tests.

If you don't have tests on this tag you can use `--no-git-rollback` to skip the rollback step.

Since not all the versions from SOT are currently having corresponding Git tags or images, you can provide refs directly using `--refs` flag, useful for testing.

```bash
ctf compat backward \
--registry <sdlc_ecr_chainlink_registry>\
--buildcmd "just cli" \
--envcmd "cl r env.toml,products/ocr2/basic.toml" \
--testcmd "cl test ocr2 TestSmoke/rounds" \
--product data-feeds \
--refs "2.36.1-rc.0" \
--refs "2.36.1-beta.0" \
--refs "2.36.1-beta.2" \
--refs "2.37.0-rc.0" \
--refs "2.37.0-beta.0" \
--refs "2.38.0-rc.0" \
--refs "2.38.0-beta.0" \
--no-git-rollback \
--don_nodes 5
```
