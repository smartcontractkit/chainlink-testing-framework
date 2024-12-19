# Test helpers

There two test helper containers:
* Killgrave
* Mockserver

Both represent HTTP mocking solutions, with Mockserver used in k8s and Killgrave outside of it. Since both will soon
be replaced by a single, in-house solution, their usage is discouraged.

That worked is begin done in [this PR](https://github.com/smartcontractkit/chainlink-testing-framework/pull/1246) and tracked in [this ticket](https://smartcontract-it.atlassian.net/browse/TT-1608).