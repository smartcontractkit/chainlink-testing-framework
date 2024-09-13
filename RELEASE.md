## Releasing a module

The Chainlink repository contains multiple independent modules. To release any of them, follow these steps:

- Merge all PRs that include the functionality you want to release (this can apply to multiple modules).
- For each module being released, add a release notes file in the .changesets folder named `vX.X.X.md`. The format is flexible, but ensure to describe all changes and any necessary upgrade procedures.
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

### Major release
- Append `vX` to Go module path, example
```
module github.com/smartcontractkit/chainlink-testing-framework/wasp/v2
```
- Tag the main branch using the format `$pkg/$subpkg/v2.X.X-alpha`
- Push the tags and visit https://github.com/smartcontractkit/chainlink-testing-framework/releases to check the release.
- Check Dependabot pipeline to analyze dependencies

## Debug Release Pipeline
Since some components of pipeline are relying on published Go modules index and Dependabot it is hard to test it, but we have a test script for that purpose:

To test release for any module use `$pkg/$subpkg/v1.999.X-test-release` tags, they are retracted so consumers can't accidentally install them
```
nix develop
python ./scripts/test-package-release.py -tag k8s-test-runner/v1.999.0-test-release -package ./k8s-test-runner
```