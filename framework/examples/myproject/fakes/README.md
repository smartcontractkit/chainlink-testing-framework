# Example Fakes

This example demonstrate how to write fake services for testing.

Copy this example into your project, write the logic of fake, build and upload it and run.

## Install

To handle some utility command please install `Taskfile`
```
brew install go-task
```

## Private Repositories (Optional)

If your tests are in a private repository please generate a new SSH key and add it on [GitHub](https://github.com/settings/keys). Don't forget to click `Configure SSO` in UI
```
task new-ssh
```

## Usage

Build it and run locally when developing fakes
```
task build -- ${product-name}-${tag} # ex. myproduct-1.0
task run
```

Test it
```
curl "http://localhost:9111/static-fake"
curl "http://localhost:9111/dynamic-fake"
```
Publish it
```
task publish -- $tag
```

Image name is `795953128386.dkr.ecr.us-west-2.amazonaws.com/ccip-fakes:test`