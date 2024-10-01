## Releasing Go modules

The Chainlink Testing Framework (CTF) repository contains multiple independent modules. To release any of them, we follow some best practices about breaking changes.

### Release strategy

Use only [lightweight tags](https://git-scm.com/book/en/v2/Git-Basics-Tagging)

**Do not move tags between commits. If something need to be fixed increment patch or minor version.**

Steps to release:

- When all your PRs are merged to `main` check the `main` branch [breaking changes badge](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/rc-breaking-changes.yaml)
- If there are no breaking changes for external methods, create a branch and explain all your module changes in `vX.X.X.md` committed under `.changeset` dir in your module. If changes are really short, and you run the [script](#check-breaking-changes-locally) locally you can push `.changeset` as a part of your final feature PR
- If there are accidental breaking changes, and it is possible to make them backward compatible - fix them
- If there are breaking changes, and we must release them change `go.mod` path, add prefix `/vX`, merge your PR(s)
- When all the changes are merged, and there are no breaking changes in the [pipeline](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/rc-breaking-changes.yaml) then proceed with releasing
- Tag `main` branch in format `$pkg/$subpkg/vX.X.X` according to your changes and push it, example:
  ```
  git tag $pkg/$subpkg/v1.1.0 && git push --tags
  git tag $pkg/$subpkg/v1.1.1 && git push --tags
  git tag $pkg/$subpkg/v2.0.0 && git push --tags
  ```
- Check the [release page](https://github.com/smartcontractkit/chainlink-testing-framework/releases)

### Binary releases

If your module have `cmd/main.go` we build binary automatically for various platforms and attach it to the release page.

## Debugging release pipeline and `gorelease` tool

Checkout `dummy-release-branch` and release it:

- `git tag dummy-module/v0.X.0`
- Add `vX.X.X.md` in `.changeset`
- `git push --no-verify --force && git push --tags`
- Check [releases](https://github.com/smartcontractkit/chainlink-testing-framework/releases)

Pipelines:

- [Main branch breaking changes](https://github.com/smartcontractkit/chainlink-testing-framework/actions/workflows/rc-breaking-changes.yaml)
- [Pipeline for releasing Go modules](.github/workflows/release-go-module.yml)

## Check breaking changes locally

We have a simple wrapper to check breaking changes for all the packages. Commit all your changes and run:

```
go run ./tools/breakingchanges/cmd/main.go
go run ./tools/breakingchanges/cmd/main.go --subdir wasp # check recursively starting with subdir
go run ./tools/breakingchanges/cmd/main.go --ignore tools,wasp,havoc,seth # to ignore some packages
```
