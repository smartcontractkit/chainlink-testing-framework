## Releasing Go modules

The Chainlink repository contains multiple independent modules. To release any of them, follow these steps:

- Merge all PRs that include the functionality you want to release (this can apply to multiple modules).
- For each module being released, add a release notes file in the `.changeset` folder named `vX.X.X.md`. The format is flexible, but ensure to describe all changes and any necessary upgrade procedures.
- Follow tagging strategy described below, push the tag and check the [release page](https://github.com/smartcontractkit/chainlink-testing-framework/releases)

### Tagging strategy

Use only [lightweight tags](https://git-scm.com/book/en/v2/Git-Basics-Tagging)

**Do not move tags between commits. If something need to be fixed increment patch or minor version.**

Major versions require `$pkg/$subpkg/vX.X.X-alpha` to be released before to analyze the scope of breaking changes across other repositories via Dependabot.
```
git tag k8s-test-runner/v2.0.0-test-release-alpha && git push --tags
```

Minor changes can be tagged as `$pkg/$subpkg/vX.X.X` and released after verifying [release page](https://github.com/smartcontractkit/chainlink-testing-framework/releases)
```
git tag k8s-test-runner/v0.6.0-test-release-alpha && git push --tags
```

### Minor and patch releases
- Tag the main branch using the format `$pkg/$subpkg/vX.X.X`
- Push the tags and visit https://github.com/smartcontractkit/chainlink-testing-framework/releases to check the release.

There should be no breaking changes in patch or minor releases, check output of `Breaking changes` on the release page, if there are - fix them and publish another patch version.

### Major releases
- Append `vX` to Go module path in `go.mod`, example:
```
module github.com/smartcontractkit/chainlink-testing-framework/wasp/v2
```
- Tag the main branch using the format `$pkg/$subpkg/v2.X.X-alpha`
- Push the tags and visit https://github.com/smartcontractkit/chainlink-testing-framework/releases to check the release.
- Check Dependabot pipeline to analyze scope of changes across other repositories

### Binary releases
If your module have `cmd/main.go` we build binary automatically for various platforms and attach it to the release page.

## Debug Release Pipeline
Since some components of pipeline are relying on published Go modules index and Dependabot we have a test script to verify the release pipeline:

To test release for any module use `$pkg/$subpkg/v1.999.X-test-release` tags, they are retracted so consumers can't accidentally install them
```
nix develop
python ./scripts/test-package-release.py -tag k8s-test-runner/v1.999.0-test-release -package ./k8s-test-runner
```

To remove all the test tags use
```
nix develop
python3 ./scripts/test-package-release.py -remove-test-tags
```

[Pipeline for releasing Go modules](.github/workflows/release-go-module.yml)
[Dependabot summary pipeline](.github/workflows/dependabot-consumers-summary.yaml)

## Check breaking changes locally
We have a simple wrapper to check breaking changes for all the packages. Commit all your changes and run:
```
./scripts/breaking-changes.sh
```