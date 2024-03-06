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

Before running tests, you must create a Docker image containing the test binary. To do this, execute the `create-test-image` command and provide the path to the test folder you wish to package.

```sh
go run ./cmd/main.go create-test-image --image-registry-url <staging-ecr-registry-url> --image-tag "<image-tag>" "<path-to-test-folder>"
```

Where `image-tag` should be a descriptive name for your test, such as "mercury-load-tests".


### Running the Test in Kubernetes

If a Docker image containing the test binary is available in an image registry (such as staging ECR), use `run` command to execute the test in K8s.

```
go run ./cmd/main.go run -c "<path-to-test-runner-toml-config>"
```

The TOML config should specify the test runner configuration as follows:

```
namespace = "wasp"
update_image = true
image_registry_url = "" # URL to the ECR containing the test binary image, e.g., staging ECR URL
image_name = "k8s-test-runner" 
image_tag = ""  # The image tag to use, like "mercury-load-tests" (see readme above)
wasp_jobs = "1"
keep_jobs = true
wasp_log_level = "debug"
test_name = "TestMercuryLoad/all_endpoints"
test_timeout = "24h"
test_config_env_name = "LOAD_TEST_BASE64_TOML_CONTENT"
test_config_file_path = "/Users/lukasz/Documents/test-configs/load-staging-testnet.toml"
```

Where:
- `test_name` is the name of the test to run (must be included in the test binary).
- `test_config_env_name` is the name of the environment variable used to provide the test configuration for the test (optional).
- `test_config_file_path` is the path to the configuration file for the test (optional).
