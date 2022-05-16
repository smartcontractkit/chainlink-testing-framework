## Remote runner
Example of running remote soak runner with a slack report
```
TEST_REGEX="@soak-ocr" SLACK_API_KEY="" SLACK_CHANNEL="" SLACK_USER_ID="" CHAINLINK_VERSION="1.4.1-rc1-nonroot" REMOTE_RUNNER_CONFIG_FILE="../runner.yaml" ginkgo -r --focus=@setup-soak soak
```