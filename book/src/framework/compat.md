# Compatibility Testing

Compatibility testing verifies that your product remains functional when Chainlink nodes are upgraded from older versions to new ones. The `ctf compat backward` command automates this process: it rolls back to an older version, boots the environment, runs your tests, then upgrades nodes one step at a time — re-running tests after each upgrade — until the latest version is validated.

## How the Upgrade Process Works

```
1. Select a version sequence  (e.g. v2.32.0 → v2.33.0 → v2.34.0)
2. Check out the oldest Git tag  (so environment + test code match that release)
3. Boot environment with the oldest image  →  run tests  (baseline)
4. For each next version:
     a. Pull the new Docker image
     b. Upgrade N nodes  (stop container, swap image, restart — node DB volumes preserved)
     c. Run tests again  (mixed-version cluster)
5. Repeat until the latest version is fully deployed and tested
```

Each step preserves the database volumes so node state is carried forward, exactly as it would be in a real rolling upgrade.

---

## Golden Path: Enabling Compatibility Tests in Your Repository

### Step 1 — Create an IAM Role for ECR Access

Your GitHub Actions runner needs permission to pull Chainlink node and Job Distributor (JD) images from AWS ECR.

Request the role from your infra team. The role should follow the pattern used in [gha-smartcontractkit](https://github.com/smartcontractkit/infra/tree/master/accounts/production/global/aws/iam/roles/gha-smartcontractkit) and grant:
- `ecr-public:GetAuthorizationToken` on `us-east-1` (public ECR, for CL images)
- `ecr:GetAuthorizationToken` + `ecr:BatchGetImage` + `ecr:GetDownloadUrlForLayer` on `us-west-2` (private ECR, for JD images)

### Step 2 — Add GitHub Actions Secrets

Add the following secrets to your repository (`Settings → Secrets and variables → Actions`):

| Secret name | Description | Example value |
|---|---|---|
| `PRODUCT_IAM_ROLE` | ARN of the IAM role that grants ECR pull access. Name it with your product name, for example CCV_IAM_ROLE | `arn:aws:iam::<account_id>:role/gha-smartcontractkit-<repo>` |
| `JD_REGISTRY` | Private ECR registry ID for JD images | `<production_ecr_registry_number>.dkr.ecr.us-west-2.amazonaws.com` |
| `JD_IMAGE` | Full JD image reference (used by your environment config) | `<production_ecr_registry_number>.dkr.ecr.us-west-2.amazonaws.com/job-distributor:0.22.1` |

Using the GitHub CLI:

```bash
gh secret set PRODUCT_IAM_ROLE   # paste the IAM role ARN
gh secret set JD_REGISTRY    # paste the JD registry URL
gh secret set JD_IMAGE       # paste the JD image reference
```

### Step 3 — Copy the Compat Pipeline

Copy `devenv-compat.yml` from [chainlink/sot-upgrade-workflow](https://github.com/smartcontractkit/chainlink/blob/sot-upgrade-workflow/.github/workflows/devenv-compat.yml) into your repository at `.github/workflows/devenv-compat.yml`.

The workflow performs the following on each run:

```yaml
- name: Checkout code
  uses: actions/checkout@v5
  with:
    fetch-depth: 0          # required: ctf reads all git tags

- name: Authenticate to AWS ECR
  uses: ./.github/actions/aws-ecr-auth
  with:
    role-to-assume: ${{ secrets.PRODUCT_IAM_ROLE }}
    aws-region: us-east-1
    registry-type: public

- name: Authenticate to AWS ECR (JD)
  uses: ./.github/actions/aws-ecr-auth
  with:
    role-to-assume: ${{ secrets.PRODUCT_IAM_ROLE }}
    aws-region: us-west-2
    registry-type: private
    registries: ${{ secrets.JD_REGISTRY }}

- name: Run compatibility test
  run: |
    ctf compat backward \
      --registry <your_ecr_registry>/chainlink \
      --buildcmd "just cli" \
      --envcmd "mycli r env.toml,products/myproduct/basic.toml" \
      --testcmd "mycli test myproduct TestSmoke/rounds" \
      --strip-image-suffix v \
      --upgrade-nodes 2 \
      --versions-back 3
```

`fetch-depth: 0` is required because `ctf` reads all Git tags to build the version sequence.

### Step 4 — Add a Nightly Trigger

Compatibility tests are typically run on a nightly schedule rather than on every PR. Add a nightly workflow that points to your product configuration:

See the [chainlink nightly example](https://github.com/smartcontractkit/chainlink/blob/develop/.github/workflows/devenv-nightly-compat.yml#L42) for a complete reference.

### Step 5 — Write Your Compatibility Tests

Your tests are plain Go tests invoked via the `--testcmd` argument. They must:

1. Connect to the already-running environment (do not spin up new nodes inside the test)
2. Exercise the behaviour that must remain functional across versions
3. Be runnable at the oldest supported version AND at every intermediate mixed-version state

A minimal example for [DF1](https://github.com/smartcontractkit/chainlink/blob/develop/devenv/tests/ocr2/smoke_test.go#L22):

The test command passed to `--testcmd` is an arbitrary shell command, so you can target any subset of tests:

---

## Running the Upgrade Sequence Locally

Before pushing to CI, authenticate to the SDLC ECR registry:

```bash
aws ecr get-login-password --region us-west-2 \
  | docker login --username AWS --password-stdin <sdlc_ecr_registry>
```

Install the `ctf` CLI if you haven't already — see [Getting Started](https://smartcontractkit.github.io/chainlink-testing-framework/framework/getting_started.html).

### Explicit version list

Provide the exact versions to test in order from oldest to newest. The first ref must have a working test at that Git tag.

```bash
ctf compat backward \
  --registry <sdlc_ecr_registry> \
  --buildcmd "just cli" \
  --envcmd "mycli r env.toml" \
  --testcmd "mycli test myproduct TestSmoke/rounds" \
  --refs 2.32.0 \
  --refs 2.33.0 \
  --refs 2.34.0 \
  --refs 2.35.0 \
  --upgrade-nodes 3
```

### Automatic SemVer detection (CI mode)

In CI, omit `--refs` and let `ctf` detect the last N SemVer tags automatically:

```bash
ctf compat backward \
  --registry <sdlc_ecr_chainlink_registry> \
  --buildcmd "just cli" \
  --envcmd "mycli r env.toml,products/myproduct/basic.toml" \
  --testcmd "mycli test myproduct TestSmoke/rounds" \
  --strip-image-suffix v \
  --upgrade-nodes 2 \
  --versions-back 3
```

When a new SemVer tag is pushed, the pipeline selects the last `--versions-back` tags, rolls back to the oldest, and runs the full upgrade sequence.

---

## All `ctf compat backward` Flags

| Flag | Default | Description |
|---|---|---|
| `--registry` | `smartcontract/chainlink` | Docker image registry for Chainlink nodes |
| `--refs` | _(none)_ | Explicit version list (oldest → newest). Repeat for each version. If omitted, SemVer tags are detected from git. |
| `--versions-back` | `1` | How many previous SemVer tags to include when auto-detecting (used without `--refs`) |
| `--upgrade-nodes` | `3` | Number of nodes to upgrade at each step |
| `--don_nodes` | `5` | Total size of the DON (used with `--product` for SOT modelling) |
| `--buildcmd` | `just cli` | Command to build the devenv CLI binary |
| `--envcmd` | _(required)_ | Command to start the environment |
| `--testcmd` | _(required)_ | Command to run the compatibility tests |
| `--node-name-template` | `don-node%d` | Docker container name template for CL nodes. Use `--node-name-template custom-cl-node-%d` if your DON uses custom naming. |
| `--strip-image-suffix` | _(none)_ | Strip a prefix from refs before looking them up in the registry (e.g. `v` strips the leading `v` from `v2.34.0`) |
| `--include-refs` | _(none)_ | Only include refs matching these patterns (e.g. `rc,beta`) |
| `--exclude-refs` | _(none)_ | Exclude refs matching these patterns (e.g. `rc,beta`) |
| `--no-git-rollback` | `false` | Skip checking out the oldest Git tag. Use when your tests don't live in the same repo as the node. |
| `--skip-pull` | `false` | Skip `docker pull`; use locally cached images |
| `--product` | _(none)_ | Enable SOT DON modelling for this product name (see below) |
| `--sot-url` | `https://rane-sot-app.main.prod.cldev.sh/v1/snapshot` | RANE SOT snapshot API URL |
| `--rane-add-git-tag-prefix` | _(none)_ | Prefix to add to RANE version strings to match Git tags |

---

## Advanced: Modelling Real Node Operator DON Upgrades (WIP)

When `--product` is specified, `ctf compat backward` fetches the versions currently running on real DONs from the RANE SOT data source and models the upgrade sequence across actual node operator versions.

```bash
ctf compat backward \
  --registry <sdlc_ecr_chainlink_registry> \
  --buildcmd "just cli" \
  --envcmd "mycli r env.toml,products/myproduct/basic.toml" \
  --testcmd "mycli test myproduct TestSmoke/rounds" \
  --product data-feeds \
  --don_nodes 5
```

`--product` field is connected with `jq .nodes[].jobs[].product` field from [SOT](https://rane-sot-app.main.prod.cldev.sh/v1/snapshot) data, find your product name there.

`--no-git-rollback` can be added if your product orchestration code and tests are compatible with all these versions, otherwise, Git rollback will be performed automatically.

If SOT versions do not yet have corresponding Git tags or registry images, supply them explicitly with `--refs`:

```bash
ctf compat backward \
  --registry <sdlc_ecr_chainlink_registry> \
  --buildcmd "just cli" \
  --envcmd "mycli r env.toml,products/myproduct/basic.toml" \
  --testcmd "mycli test myproduct TestSmoke/rounds" \
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

---

## Restoring Your Working Branch After Testing

`ctf compat backward` checks out old Git tags during a run. To return to your working branch:

```bash
ctf compat restore --base_branch develop
# or just:
git checkout develop
```