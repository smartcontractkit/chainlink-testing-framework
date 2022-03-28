# Soak Tests

The test framework is designed around you launching a test environment to a K8s cluster, then running the test from your personal machine. Your personal machine coordinates the chainlink nodes, reads the blockchain, etc. This works fine for running tests that take < 5 minutes, but soak tests often last days or weeks. So the tests become dependent on your local machine maintaining power and network connection for that time frame. This quickly becomes untenable. So the solution is to launch a `remote-test-runner` container along with the test environment.

## Running the Test

The soak tests are triggered by the [soak_runner_test.go](./suite/soak/soak_runner_test.go) tests. When running, the test will check for a local config file: `remote_runner_config.yaml`. If it's not already created, it will generate one with some default values, and then inform you that you should modify those values.

```yaml
test_regex: '@soak-ocr' # The regex of the test name to run
test_directory: /Users/adam/Projects/integrations-framework/suite/soak/tests # The directory where the go tests you want the remote runner to run
# Slack values are covered below
```

Modify these values that make sense for the tests you want to run. Once the values are modified, you can run the test again. The soak runner test will then compile the tests that you pointed to by the `test_directory` into a `remote.test` executable. This executable, including your local `framework.yaml` and `networks.yaml` configs are uploaded to the remote test runner. Make sure to read the section below to take advantage of slack integration.

## Watching the Test

The rest of the `remote_runner_config.yaml` file holds various Slack bot params to notify you when the test finishes.

```yaml
slack_webhook_url: https://hooks.slack.com/services/XXX # A slack webhook URL to send notification on the test
slack_api_key: abcdefg # A Slack API key to upload test results with. This API Key needs to have `file:write` permissions
slack_channel: C01xxxxx # The Slack Channel ID (open your Slack channel details and copy the ID there)
slack_user_id: U01xxxxx # Your Slack member ID https://zapier.com/help/doc/common-problems-slack
```

The Slack Bot will need to have the following:

* A webhook, [they're easy to setup](https://api.slack.com/messaging/webhooks)
* Permission for [files:write](https://api.slack.com/scopes/files:write)
* The bot must be invited into the channel you want it to notify in: `/invite @botname` HINT: If you get the error `not_in_channel` this is likely what you need to set up.

## After the Test

The test environment **will stay active until you manually delete it from your Kubernetes cluster**. This keeps the test env alive so you can view the logs when the test is done. You can do so by [using kubectl](https://www.dnsstuff.com/how-to-tail-kubernetes-and-kubectl-logs), something like [Lens](https://k8slens.dev/), or use the `chainlink-soak-xxyyx.yaml` as an environment file with our handy [helmenv](https://github.com/smartcontractkit/helmenv) tool.

Using the `helmenv` cli tool, you can input the generated file to download logs like so.

`envcli dump -e chainlink-long-xxyyx.yaml -a soak_test_logs -db chainlink`
