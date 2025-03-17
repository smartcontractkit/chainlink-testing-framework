# Finding the Root Cause of Test Flakes in Go

Flaky tests can arise from many sources and can be frustrating to fix. Here's a non-exhaustive guide to help you find and resolve common causes for flakes in Go. But first, to answer a common question...

## The Test Only Flakes 0.xx% of the Time, Why Bother Fixing It?

You bother fixing it because of **MATH!**

Let's imagine a large repo with 10,000 tests, and let's imagine only 100 (1%) of them are flaky. Let's further imagine that each of those flaky tests has a chance of flaking 1% of the time. If you are a responsible dev that requires all of your tests to pass in CI before you merge, flaky tests have now become a massive headache. 

$$P(\text{at least one flaky test}) = 1 - (1 - 0.01)^{100}$$

$$P(\text{at least one flaky test}) \approx 63.40\%$$

Even a small percentage of tests with a small chance of flaking can cause massive damage to dev velocity.

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

Sometimes you can only discover the truth by going directly to the source. Before you do so, please double check what `runs_on` systems your workflows use. If you're only using `ubuntu-latest` runners, these runs should be free, or at least very cheap. `8-core`, `16-core`, and `32-core` workflows can become very expensive, very quickly. Please use caution and discretion when running these workflows repeatedly.

#### 6.1 CI Resource Constraints

It is sometimes the case that tests only fail in CI environments because those environments are underpowered. **This is more rare than you think, be cautious of [System 1 thinking](https://en.wikipedia.org/wiki/Thinking,_Fast_and_Slow) here.** You can diagnose this with [this excellent GitHub workflow telemetry action](https://github.com/catchpoint/workflow-telemetry-action) that can give you detailed stats on how many resources your tests are consuming. (This is also handy if you're looking to optimize your CI runtimes or costs.) If your tests are flaking due to low resources, consider other options before just increasing the power of the CI runners. [Increasing the power of a GitHub Actions workflow by a single tier doubles its cost](https://docs.github.com/en/billing/managing-billing-for-your-products/managing-billing-for-github-actions/about-billing-for-github-actions#per-minute-rates-for-x64-powered-larger-runners). If your workflow runs often, you can burn a lot of cash quickly. You can otherwise try strategies like:

* Splitting the tests into different workflows, each running on `ubuntu-latest`
* Moving more resource-hungry tests to run only on nightly cadences
* Try removing `t.Parallel()` from subtests, as too many tests trying to run at once will often hurt stability and runtimes on smaller machines

### 7. Fix It!

Maybe you've found the source of the flake and are now drilling down into the reasons why. Whatever those reasons might be, I urge you to, at least briefly, reframe the problem and ask if the test is actually working as intended, and it is revealing flaky behavior in your application instead. Consider that you might have found a rare bug, rather than a rare flake.

### 8. Give Up

It's not my favorite answer, but sometimes this truly is the solution. It's hard to know exactly when you should abandon hope, but maybe the below steps can help you figure it out.

#### 8.1 Evaluate the Importance of the Test

Ask yourself these questions to help figure out if it's worth working on this flake further, and to help you figure out what to do next.

* What does the test actually check? Is it a critical path?
* Is the test flaking because it's a bad test? Or it's trying to test behavior that shouldn't or can't be tested?
* Can you write a new test that checks the same behavior, but doesn't fall to the same issues?
* Can you come back to this later? Maybe in a week or two you'll have new ideas, or maybe the underlying system will change in ways that this flake is no longer an issue?

#### 8.2 Turn it Off

Assuming you're ready to declare defeat, it's time to turn off the test. How you do this depends on the test, your team, and the answers to the questions above. If you've determined the test isn't particularly important and isn't worth running anymore, you should just delete it. 

If the test does check crucial behavior and you'd like to return to it, please consider this general flow:

1. Make a ticket noting the flake, and write down as much as you've figured out so far. Ideally have the ticket setup to remind you or your team in a couple of weeks that it's still there.
2. Mark the test with `t.Skip("Flaky:<Ticket-Name>")` so that Flakeguard can properly identify tests that are being skipped for flaky issues.

## Chainlink E2E Tests

At CLL, we have specially designed E2E tests that run in Docker and Kubernetes environments. They're more thorough validations of our systems, and much more complex than typical unit tests.

### 1. Find Flakes

If you don't already have some Flakeguard results to work with, you can run flakeguard on E2E tests easily with the following command from the `chainlink-core/` directory:

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

### 2. Explore Logs and Debug

E2E tests are complex beasts, orchestrating many systems together at once. As such, the best strategy is usually to mark the time of the test failure, and get busy cross-referencing that timestamp across the logs of each component so you can see what was happening in each component in the test. The [CTF Debug Docs](https://smartcontractkit.github.io/chainlink-testing-framework/framework/components/debug.html) can help increase the logs you're collecting, and help use more traditional debugging tools like [Delve](https://github.com/go-delve/delve).

### 3. Remember You're (Sort of) in the Real World

E2E tests are meant to closely simulate real-world situations and deployments, so causes that never pop up in unit tests become common hazards:

* **Networking issues**: A blip in networking connections can throw off timings, and make expected test states happen in an unexpected order.
* **Timing issues**: Many E2E tests are set to fail after not receiving an expected result after x seconds. It's possible blips in infrastructure and networking slow it down, and expanding this timeout can help stabilize things.
* **GitHub Actions degraded performance**: This is a common scapegoat for when you're frustrated, but it's worth [checking GitHub's status page](https://www.githubstatus.com/) for any incidents that might have occurred while your test was running.

### 4. Check the Test's Resources

E2E tests are much more resource hungry than your typical unit test suite. Much of what I mentioned in our general tips for [CI resources](####-6.1-ci-resource-constraints) applies here. You can turn on CI telemetry easily in most of our E2E tests by setting [collect_test_telemetry](https://github.com/smartcontractkit/.github/blob/main/.github/workflows/run-e2e-tests.yml#L174), like so:

```yaml
uses: smartcontractkit/.github/.github/workflows/run-e2e-tests.yml@version
with:
  collect_test_telemetry: true
```