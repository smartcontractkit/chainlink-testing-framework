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
--nodes 3
```

Keep in mind that `refs` should be present in regsitry you are testing against, the first (oldest) `ref` should also have a valid end-to-end test that works.

In CI we detect SemVer tags automatically, whenever a new tag appears we select last 3, rollback to the oldest and perform upgrade process.

```bash
ctf compat backward \
--registry <sdlc_ecr_registry> \
--buildcmd "just cli" \
--envcmd "cl r" \
--testcmd "cl test ocr2 TestSmoke/rounds" \
--nodes 3 \
--versions-back 3
```

In case you have multiple DONs in your product and names of nodes are different please use `--node-name-template custom-cl-node-%d` option

## Modelling Node Operators Cluster

It is possible to fetch versions node operators are currently running and model DON upgrade sequence locally. Logic is the same, get all the versions, rollback to the oldest one, setup product, verify, try to upgrade all the versions running the oldest test for each upgrade.

```bash
ctf compat backward \
--registry <sdlc_ecr_registry>\
--buildcmd "just cli" \
--envcmd "cl r" \
--testcmd "cl test ocr2 TestSmoke/rounds" \
--nop northwestnodes \
--versions-back 3 \
--nodes 3
```
