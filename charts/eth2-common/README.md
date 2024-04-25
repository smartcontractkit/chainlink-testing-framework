# Purpose
This package contains Helm charts that are shared between different Eth2 clients. It's meant to be used as a dependency in other charts. There's little point in using it directly.

It contains following elements/deployments:
* validators' key generation
* eth1 and eth2 genesis
* creation of other shared files (like jwt secret or keystore password files)
* Prysm beacon chain deployment
* Prysm validator deployment
* chain-ready deployment (container that waits until blocks are produced)
* persistent volume and claim definitions

# Usage
In your package you need to define it as dependency:
```yaml
# when using local filesystem version
dependencies:
  - name: eth2-common
    version: 0.1.0
    repository: "file://../eth2-common"

# when using remote version (that should be the 'prod' version)
dependencies:
  - name: eth2-common
    version: 0.1.0
    repository: "https://raw.githubusercontent.com/smartcontractkit/qa-charts/gh-pages/"
```

Some of the elements will be automatically used from dependency as soon as they used anywhere and matched by name. In our case that's true for volume and claim definitions and config maps.

Even though Prysm Beacon Chain and Validator are the same regarding of the execution client there's no way that I know of to include whole resource from dependency. That's why to define beacon chain deployment you need to use following syntax in your `prysm-beacon.deployment.yaml`
```yaml
{{- include "eth2-common.templates.deployment.prysm-beacon" .  }}
```

You need to do similar thing for prysm beacon service, validator deployment and chain ready deployment.

Remember, that you cannot reuse predefined values from the dependency, you need to define them in your package.
