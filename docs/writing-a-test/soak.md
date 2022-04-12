---
layout: default
title: Soak Tests
nav_order: 2
parent: Writing a Test
---

# Soak Tests

Soak tests refer to running longer tests, that can take anywhere from hours to days to see how the application fares over long stretches of time. See some examples in our [soak test suite](https://github.com/smartcontractkit/integrations-framework/tree/main/suite/soak).

The test framework is designed around you launching a test environment to a K8s cluster, then running the test from your personal machine. Your personal machine coordinates the chainlink nodes, reads the blockchain, etc. This works fine for running tests that take < 5 minutes, but soak tests often last days or weeks. So the tests become dependent on your local machine maintaining power and network connection for that time frame. This quickly becomes untenable. So the solution is to launch a `remote-test-runner` container along with the test environment.

## Writing the Test

Since the test is being run from a `remote-test-runner` instead of your local machine, setting up and tearing down the test environment is a little different. So is connecting to things like the blockchain networks, chainlink nodes, and the mock adapter. Most other interactions should be as normal though.

```go
// Connects to the soak test resources from the `remote-test-runner`
env, err := environment.DeployOrLoadEnvironmentFromConfigFile(
  tools.ChartsRoot,      // Default location of helm charts to look for
  "/root/test-env.json", // Default location for the soak-test-runner container
)
log.Info().Str("Namespace", env.Namespace).Msg("Connected to Soak Environment")

// Run test logic

// Teardown remote suite
if err := actions.TeardownRemoteSuite(keeperBlockTimeTest.TearDownVals()); err != nil {
  log.Error().Err(err).Msg("Error tearing down environment")
}
log.Info().Msg("Soak Test Concluded")
```

## Running the Test

The soak tests are triggered by the [soak_runner_test.go](https://github.com/smartcontractkit/integrations-framework/blob/main/suite/soak/soak_runner_test.go) tests, or with `make test_soak`. When running, the test will check for a local config file: `remote_runner_config.yaml`. If it's not already created, it will generate one with some default values, and then inform you that you should modify those values.

```yaml
test_regex: '@soak-ocr' # The regex of the test name to run
test_directory: /Users/adam/Projects/integrations-framework/suite/soak/tests # The directory where the go tests you want the remote runner to run
# Slack values are covered below
```

Modify these values that make sense for the tests you want to run. Once the values are modified, you can run the test again. The soak runner test will then compile the tests that you pointed to by the `test_directory` into a `remote.test` executable. This executable, including your local `framework.yaml` and `networks.yaml` configs are uploaded to the remote test runner. Make sure to read the section below to take advantage of slack integration.

## Watching the Test

The rest of the `remote_runner_config.yaml` file holds various Slack bot params to notify you when the test finishes.

```yaml
slack_api_key: abcdefg # A Slack API key to upload test results with. This API Key needs to have `file:write` permissions
slack_channel: C01xxxxx # The Slack Channel ID (open your Slack channel details and copy the ID there)
slack_user_id: U01xxxxx # Your Slack member ID https://zapier.com/help/doc/common-problems-slack
```

The Slack Bot will need to have the following:

* Permission for [files:write](https://api.slack.com/scopes/files:write)
* Permission for [chat:write](https://api.slack.com/scopes/chat:write)
* The bot must be invited into the channel you want it to notify in: `/invite @botname` HINT: If you get the error `not_in_channel` this is likely what you need to set up.

## After the Test

The test environment **will stay active until you manually delete it from your Kubernetes cluster**. This keeps the test env alive so you can view the logs when the test is done. You can do so by [using kubectl](https://www.dnsstuff.com/how-to-tail-kubernetes-and-kubectl-logs), something like [Lens](https://k8slens.dev/), or use the `chainlink-soak-xxyyx.yaml` as an environment file with our handy [helmenv](https://github.com/smartcontractkit/helmenv) tool.

Using the `helmenv` cli tool, you can input the generated file to download logs like so.

`envcli dump -e chainlink-long-xxyyx.yaml -a soak_test_logs -db chainlink`
