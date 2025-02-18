# Finding the Root Cause of E2E Test Flakes

## Introduction
When end-to-end tests fail intermittently, the underlying issues can stem from resource constraints, environment setup, or test design—among other possibilities. This guide helps engineers systematically diagnose and address E2E test flakiness, reducing the time spent on guesswork and repeated failures.

---

## 1. GitHub Runners' Hardware
GitHub provides **hosted runners** with specific CPU, memory, and disk allocations. If your tests require more resources than these runners can provide, you may encounter intermittent failures.

By default, we run tests on **`ubuntu-latest`**, as it is **free for public repositories** and the **most cost-effective option for private repositories**. However, this runner has limited resources, which can lead to intermittent failures in resource-intensive tests.

> **Note:** `ubuntu-latest` for **private repositories** has weaker hardware compared to `ubuntu-latest` for **public repositories**. You can learn more about this distinction in [GitHub's documentation](https://docs.github.com/en/actions/using-github-hosted-runners/using-github-hosted-runners/about-github-hosted-runners#standard-github-hosted-runners-for-public-repositories).

### 1.1 Available GitHub Runners
Below are the some of the GitHub-hosted runners available in our organization:

| Runner Name | CPU | Memory | Disk |
|------------|-----|--------|------|
| `ubuntu-22.04-4cores-16GB` | 4 cores | 16 GB RAM | 150 GB SSD |
| `ubuntu-latest-4cores-16GB` | 4 cores | 16 GB RAM | 150 GB SSD |
| `ubuntu-22.04-8cores-32GB` | 8 cores | 32 GB RAM | 300 GB SSD |
| `ubuntu-latest-8cores-32GB` | 8 cores | 32 GB RAM | 300 GB SSD |
| `ubuntu-22.04-8cores-32GB-ARM` | 8 cores | 32 GB RAM | 300 GB SSD |


### 1.2 Tips for Low-Resource Environments
- **Profile your tests** to understand their CPU and memory usage.  
- **Optimize**: Only spin up what you need.  
- **If resources are insufficient**, consider redesigning your tests to run in smaller, independent chunks.
- **If needed**, you can configure CI workflows to use a higher-tier runner, but this comes at an additional cost.

---

## 2. Reproducing Flakes
Flaky tests don't fail on every run, so you need to execute them multiple times to isolate problems.

### 2.1 Repeat Runs
For E2E tests, run them 5–10 times consecutively to expose intermittent issues.  

### 2.2 Flaky Unit Tests in the Core Repository
For unit tests in the core repository, you can use a dedicated command to detect flakiness in an updated test:


```sh
cd chainlink-core/
make run_flakeguard_validate_tests
```


## 3. Testing Locally Under CPU and Memory Constraints

If CPU throttling or resource contention is suspected, here's how you can approach testing under constrained resources:

1. **Spin up Docker containers locally with limited CPU or memory.**  
2. **Mimic GitHub's environment** (use the same OS, similar resource limits).  
3. **Run E2E tests** repeatedly to see if flakiness correlates with resource usage.  
4. **Review logs and metrics** for signs of CPU or memory starvation.


### Setting Global Limits (Docker Desktop)
If you are using **Docker Desktop** on **macOS or Windows**, you can globally limit Docker's resource usage:

1. Open **Docker Desktop**.
2. Navigate to **Settings** → **Resources**.
3. Adjust the sliders for **CPUs** and **Memory**.
4. Click **Apply & Restart** to enforce the new limits.

This setting caps the **total** resources Docker can use on your machine, ensuring all containers run within the specified constraints.


### Observing Test Behavior Under Constraints
- **Run your E2E tests repeatedly** with different global resource settings.
- Watch for flakiness: If tests start failing more under tighter limits, suspect CPU throttling or memory starvation.
- **Examine logs/metrics** to pinpoint if insufficient resources are causing sporadic failures.

By setting global limits, you can simulate resource-constrained environments similar to CI/CD pipelines and detect potential performance bottlenecks in your tests.


## 4. Common Pitfalls and “Gotchas”
1. **Resource Starvation**: Heavy tests on minimal hardware lead to timeouts or slow responses.  
2. **External Dependencies**: Network latency, rate limits, or third-party service issues can cause sporadic failures.  
3. **Shared State**: Race conditions arise if tests share databases or global variables in parallel runs.  
4. **Timeouts**: Overly tight time limits can fail tests on slower environments.


## 5. Key Takeaways
Tackle flakiness systematically:
1. **Attempt local reproduction** (e.g., Docker + limited resources).  
2. **Run multiple iterations** on GitHub runners.  
3. **Analyze logs and metrics** to see if resource or concurrency issues exist.  
4. **Escalate** to the infra team only after confirming the issue isn't in your own test code or setup.
