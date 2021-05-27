# Chainlink Testing Env

Borrowing from [here](https://github.com/smartcontractkit/chainlink-smoke-tests) to get a basic setup going.
Ideally we'd wait for an official env setup from infra, or launch this before running tests. Issue is that we need a few
things that are determined BY our test situation in order to successfully have this env. We need at least:

* A valid LINK address
* An ethereum URL
* An ethereum chain ID

Seeing how the above variables address change from test to test, we need to first deploy contracts, then deploy the 
chainlink nodes. This is intended as a stop gap.
