## Preparing to Run Tests on Staging

Ensure you complete the following steps before executing tests on the staging environment:

1. **Connect to the VPN**

2. **AWS Login with Staging Profile**

   Authenticate to AWS using your staging profile, specifically with the `StagingEKSAdmin` role. Execute the following command:

   ```sh
   aws sso login --profile staging
   ```

3. **Verify Authorization**

   Confirm your authorization status by listing the namespaces in the staging cluster. Run `kubectl get namespaces`. If you see a list of namespaces, this indicates successful access to the staging cluster.

## Running Tests

### Creating an Image with the Test Binary

Before running tests, you must create a Docker image containing the test binary. To do this, execute the `create-test-image` command and provide the path to the test folder you wish to package. This command:

1. Compiles test binary under `<path-to-test-folder>`
2. Creates a docker image with the test binary
3. Pushes the docker image to the image registry (e.g. Staging ECR)

```sh
go run ./cmd/main.go create-test-image --image-registry-url <staging-ecr-registry-url> --image-name "<image-name>" --image-tag "<image-tag>" "<path-to-test-folder>"
```

Where `image-tag` should be a descriptive name for your test, such as "mercury-load-tests".

### Running the Test in Kubernetes

If a Docker image containing the test binary is available in an image registry (such as staging ECR), use `run` command to execute the test in K8s.

```
go run ./cmd/main.go run -c "<path-to-test-runner-toml-config>"
```

The TOML config should specify the test runner configuration as follows:

```
namespace = "e2e-tests"
rbac_role_name = "" # RBAC role name for the chart
rbac_service_account_name = "" # RBAC service account name for the chart
image_registry_url = "" # URL to the ECR containing the test binary image, e.g., staging ECR URL
image_name = "k8s-test-runner"
image_tag = ""  # The image tag to use, like "mercury-load-tests" (see readme above)
job_count = "1"
test_name = "TestMercuryLoad/all_endpoints"
test_timeout = "24h"
test_config_base64_env_name = "LOAD_TEST_BASE64_TOML_CONTENT"
test_config_file_path = "/Users/lukasz/Documents/test-configs/load-staging-testnet.toml"
resources_requests_cpu = "1000m"
resources_requests_memory = "512Mi"
resources_limits_cpu = "2000m"
resources_limits_memory = "1024Mi"
[envs]
WASP_LOG_LEVEL = "info"
TEST_LOG_LEVEL = "info"
MERCURY_TEST_LOG_LEVEL = "info"
[metadata.labels]
"chain.link/component" = "test-runner"
"chain.link/product" = "<your-product-name>"
"chain.link/team" = "<name–of-the-team-you're-running-the-test-for>"
"chain.link/cost-center" = "test-tooling-<testType>-test"
```

> [NOTE]
> Make sure to quote labels with "/" as otherwise parsing them will fail.

Where:

- `test_name` is the name of the test to run (must be included in the test binary).
- `test_config_env_name` is the name of the environment variable used to provide the test configuration for the test (optional).
- `test_config_file_path` is the path to the configuration file for the test (optional).

## Using K8s Test Runner on CI

### Example

This example demonstrates the process step by step. First, it shows how to download the Kubernetes Test Runner. Next, it details the use of the Test Runner to create a test binary specifically for the Mercury "e2e_tests/staging_prod/tests/load" test package. Finally, it describes executing the test in Kubernetes using a customized test runner configuration.

```
- name: Download K8s Test Runner
    run: |
        mkdir -p k8s-test-runner
        cd k8s-test-runner
        curl -L -o k8s-test-runner.tar.gz https://github.com/smartcontractkit/chainlink-testing-framework/releases/download/v0.2.4/test-runner.tar.gz
        tar -xzf k8s-test-runner.tar.gz
        chmod +x k8s-test-runner-linux-amd64
```

Alternatively, you can place the k8s-test-runner package within your repository and unpack it:

```
- name: Unpack K8s Test Runner
    run: |
        cd e2e_tests
        mkdir -p k8s-test-runner
        tar -xzf k8s-test-runner-v0.0.1.tar.gz -C k8s-test-runner
        chmod +x k8s-test-runner/k8s-test-runner-linux-amd64
```

Then:

```
- name: Build K8s Test Runner Image
    if: github.event.inputs.test-type == 'load' && github.event.inputs.rebuild-test-image == 'yes'
    run: |
        cd e2e_tests/k8s-test-runner

        ./k8s-test-runner-linux-amd64 create-test-image --image-registry-url "${{ secrets.AWS_ACCOUNT_ID_STAGING }}.dkr.ecr.${{ secrets.AWS_REGION }}.amazonaws.com" --image-tag "mercury-load-test" "../staging_prod/tests/load"

- name: Run Test in K8s
    run: |
        cd e2e_tests/k8s-test-runner

        cat << EOF > config.toml
        namespace = "e2e-tests"
        rbac_role_name = "" # RBAC role name for the chart
        rbac_service_account_name = "" # RBAC service account name for the chart
        image_registry_url = "${{ secrets.AWS_ACCOUNT_ID_STAGING }}.dkr.ecr.${{ secrets.AWS_REGION }}.amazonaws.com"
        image_name = "k8s-test-runner"
        image_tag = "mercury-load-test"
        job_count = "1"
        chart_path = "./chart"
        test_name = "TestMercuryLoad/all_endpoints"
        test_timeout = "24h"
        resources_requests_cpu = "1000m"
        resources_requests_memory = "512Mi"
        resources_limits_cpu = "2000m"
        resources_limits_memory = "1024Mi"
        test_config_base64_env_name = "LOAD_TEST_BASE64_TOML_CONTENT"
        test_config_base64 = "${{ steps.conditional-env-vars.outputs.LOAD_TEST_BASE64_TOML_CONTENT }}"
        [envs]
        WASP_LOG_LEVEL = "info"
        TEST_LOG_LEVEL = "info"
        MERCURY_TEST_LOG_LEVEL = "info"
        [metadata.labels]
        "chain.link/component" = "test-runner"
        "chain.link/product" = "data-streamsv0.3"
        "chain.link/team" = "Data Streams"
        "chain.link/cost-center" = "test-tooling-load-test"
        EOF

        ./k8s-test-runner-linux-amd64 run -c config.toml
```

## Release

Run `./package <version>`
