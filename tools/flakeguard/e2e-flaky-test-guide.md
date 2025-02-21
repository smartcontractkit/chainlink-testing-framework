# Finding the Root Cause of Test Flakes in Go

Flaky tests can arise from many sources and can be frustrating to fix. Here's a non-exhaustive guide to help you find and resolve common causes for flakes in Go.

## The Test Only Flakes 0.xx% of the Time, Why Bother Fixing It?

You bother to fix it because of **MATH!**

Let's imagine a large repo with 10,000 tests, and let's imagine only 100 (1%) of them are flaky. Let's further imagine that each of those flaky tests has a chance of flaking 1% of the time. If you are a responsible dev that requires all of your tests to pass in CI before you merge, flaky tests have now become a massive headache. 

$$P(\text{at least one flaky test}) = 1 - (1 - 0.01)^{100}$$

$$P(\text{at least one flaky test}) \approx 63.40\%$$

Even a few tests with a tiny chance of a flaking can cause massive damage to a repo that a lot of devs work on.

## General Tips

Ideally, if you're dealing with a flaky test, you'll already have some examples of it flaking in front of you so you can dig through logs and stack traces and figure it out that way. If that's not the case, or you'd like some more evidence, or you're just stumped, try reproducing the flake. How you reproduce the flake is often the best clue as to why its flaking.

For repos that have [flakeguard](https://github.com/smartcontractkit/chainlink-testing-framework/tree/main/tools/flakeguard) configured (like chainlink), you can try running it locally.

```sh
make run_flakeguard_validate_unit_tests
```

You can also try some more precise configurations below.

### 1. Run the Test in Isolation

As we saw above, flaky tests become issues even when their chance of flaking is tiny. You might be hunting down a flake that only happens 0.5% of the time, so you're only real solution is to run the test over and over. 

```sh
# Run just that test 1,000 times, stopping after the first failure
go test ./package -run TestName -count 1000 -failfast
```

### 2. Run the Test Package

Tests rarely run in isolation in the real world. If you can't get the flake to happen when isolated, try running the whole package on repeat.

```sh
# Run all tests in the package over and over.
go test ./package -count 1000 -failfast
```

If you get the test to fail here, but not independently, it's likely that it depends on the execution of other tests in the package. Look for global resources your test could be sharing with others, and do your best to isolate all of your unit tests.

### 3. Randomize Test Order

If that's still not doing the job, or you're still scratching your head, try randomizing the test order. Go runs tests in a deterministic order by default, but Go's idea of "deterministic" is pretty liberal.

```sh
# -shuffle randomizes test order
go test ./package -shuffle on -count 1000 -failfast
# You can supply your own int value to shuffle as a seed
go test ./package -shuffle 15 -count 1000 -failfast
```

### 4. Check for Races

If your test is failing in a situation like this, it's possible there's a race condition it's getting caught on. Go's `-race` flag isn't guaranteed to catch all races every time. Just like flakes, you sometimes just need to get lucky (unlucky?).

```sh
# Tests with -race detection take longer to run, and aren't always going to catch issues, especially in large test suites.
go test ./package -race -shuffle on -count 100 -failfast
```

### 5. Emulate Your Target System

Tests will often fail in CI, but not locally. You can try re-running the test in CI, but this might take a long time, cost a lot of money, or generally be annoying. There are a few tricks you can do to emulate CI environments locally.

#### 5.1 Play with -cpu and -parallel

You can artificially constrain or expand parallel execution directly in go. [GOMAXPROCS](https://pkg.go.dev/runtime#hdr-Environment_Variables) is set to the amount of CPUs your system has by default, and controls how many OS threads can run Go code at once. You can manipulate this value, or otherwise play with how many tests can run at once easily. This can help you figure out if resource constraints are hurting your tests.

```sh
# Use -cpu to change GOMAXPROCS. You can supply a list of values to try out different values at once
go test ./package -shuffle 15 -count 1000 -failfast -cpu 1,2,4
# Use -parallel to set the max amount of tests allowed to run in parallel at once
go test ./package -shuffle 15 -count 1000 -failfast -parallel 4
```

#### 5.2 Use Docker

Docker can help you emulate your CI environment a little better. You can lookup what type of GitHub Actions runner your CI workflow uses by matching to the lists [here](https://docs.github.com/en/actions/using-github-hosted-runners/using-github-hosted-runners/about-github-hosted-runners#standard-github-hosted-runners-for-public-repositories) and [here](https://docs.github.com/en/actions/using-github-hosted-runners/using-larger-runners/about-larger-runners#specifications-for-general-larger-runners). You can then package your Go tests in a Docker container, and run them with varying resources.

```sh
# Run the default 4-core-16GB ubuntu-latest image used for public GitHub repos
docker run -it --cpus=4 --memory="16g" ubuntu-24.04
```

You can also try using [dockexec](https://github.com/mvdan/dockexec) for convenience, but I've never personally tried it.

#### 5.3 Use act

[act](https://github.com/nektos/act) is a project that lets you emulate your GitHub Actions workflows locally. It's not perfect, and can be tricky to setup for more complex workflows, but it is a nice option for if you suspect issues are further back in the workflow, and don't want to run the full CI process.

### 6. Use Your Target System

Sometimes you can only discover the truth by going directly to the source. Before you do so, please double check what `runs_on` systems your workflows use. If you're only using `ubuntu-latest` runners, these runs should be free. `8-core`, `16-core`, and `32-core` workflows can become very expensive, very quickly. Please use caution and discretion when running these workflows repeatedly.

### 7. Fix It!

Maybe you've found the source of the flake and are now drilling down into the reasons why. Whatever those reasons might be, I urge you to, at least briefly, reframe the problem and ask if the test is actually working as intended, and it is revealing flaky behavior in your application instead. This might be an opportunity to fix a rare bug instead of force a test to conform to it.

### 8. Give Up

It's not my favorite answer, but sometimes this truly is the solution. It's hard to know exactly when this point is. I hope to eventually gather enough data on dev productivity and how flaky tests affect them that I can give you absolute rules, but until then, we'll have to go off vibes. Here's your checklist for when you feel ready to collapse in defeat.

#### 8.1 Evaluate the Importance of the Test

* What does the test actually check? Is it a critical path? 
* Is the test flaking because it's a bad test? Or it's trying to test behavior that shouldn't or can't be tested? TODO:
* 

#### 8.2 How Flaky is the Test?

Flakeguard should give you a good idea of the test's percentage chance of flaking. Remember from above that even

## Chainlink E2E Tests

At CLL, we have specially designed E2E tests that run in Docker and Kubernetes environments. They're more thorough validations of our systems, and much more complex than typical unit tests.

### 1. Find Flakes

You should already have examples thanks to flakeguard TODO:

### 2. Reproduce Flakes

For E2E tests, run them 5–10 times consecutively to expose intermittent issues. To run the tests with flakeguard validation, execute the following command from the `chainlink-core/` directory:

```sh
cd chainlink-core/
make run_flakeguard_validate_e2e_tests
```

You’ll be prompted to provide:
- **Test IDs** (e.g., `smoke/forwarders_ocr2_test.go:*,smoke/vrf_test.go:*`)

  *Note: Test IDs can be taken from the `e2e-tests.yml` file.*

- **Number of runs** (default: 5)
- **Chainlink version** (default: develop)
- **Branch name** (default: develop)

This is generally enough 

### 2. Check Resource Constraints

GitHub provides **hosted runners** with specific CPU, memory, and disk allocations. If your tests require more resources than these runners can provide, you may encounter intermittent failures.

By default, we run tests on **`ubuntu-latest`**, as it is **free for public repositories** and the **most cost-effective option for private repositories**. However, this runner has limited resources, which can lead to intermittent failures in resource-intensive tests.

> **Note:** `ubuntu-latest` for **private repositories** has weaker hardware compared to `ubuntu-latest` for **public repositories**. You can learn more about this distinction in [GitHub's documentation](https://docs.github.com/en/actions/using-github-hosted-runners/using-github-hosted-runners/about-github-hosted-runners#standard-github-hosted-runners-for-public-repositories).

### 1.1 Available GitHub Runners
Below are the some of the GitHub-hosted runners available in our organization:

| Runner Name                    | CPU     | Memory    | Disk       |
| ------------------------------ | ------- | --------- | ---------- |
| `ubuntu-22.04-4cores-16GB`     | 4 cores | 16 GB RAM | 150 GB SSD |
| `ubuntu-latest-4cores-16GB`    | 4 cores | 16 GB RAM | 150 GB SSD |
| `ubuntu-22.04-8cores-32GB`     | 8 cores | 32 GB RAM | 300 GB SSD |
| `ubuntu-latest-8cores-32GB`    | 8 cores | 32 GB RAM | 300 GB SSD |
| `ubuntu-22.04-8cores-32GB-ARM` | 8 cores | 32 GB RAM | 300 GB SSD |


### 1.2 Tips for Low-Resource Environments

- **Profile your tests** to understand their CPU and memory usage.  
- **Optimize**: Only spin up what you need.  
- **If resources are insufficient**, consider redesigning your tests to run in smaller, independent chunks.
- **If needed**, you can configure CI workflows to use a higher-tier runner, but this comes at an additional cost.
- **Run with debug logs** or Delve debugger. For more details, check out the [CTF Debug Docs.](https://smartcontractkit.github.io/chainlink-testing-framework/framework/components/debug.html)

---

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
