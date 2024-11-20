# WASP - Profile

A **Profile** allows you to combine load from different generators. Each generator operates according to its own schedule, and you can even mix RPS and VU load types, as shown in [this example](https://github.com/smartcontractkit/chainlink-testing-framework/blob/main/wasp/examples/profiles/node_mixed_test.go).

> [!NOTE]  
> The `error` returned by the `.Run(true)` function of a `Profile` might indicate the following:
> * load generator did not start and instead returned an error
> * or **if** `Profile` was created with `WithGrafana` option that enabled `CheckDashboardAlertsAfterRun` that at some alerts fired.
> To check for errors during load generation, you need to call the `Errors()` method on each `Generator` within the `Profile`.

---

### Timing Considerations

It is not possible to have different generators in the same `Profile` start load generation at different times.  
If you need staggered start times, you must use separate `Profiles` and handle timing in your `Go` code.  
An example demonstrating this can be found [here](../how-to/parallelise_load.md).
