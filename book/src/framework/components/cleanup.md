# Components Cleanup

Managing state is challenging, especially in end-to-end testing, we use [ryuk](https://golang.testcontainers.org/features/garbage_collector/#ryuk) and following simple rules:
- If `TESTCONTAINERS_RYUK_DISABLED=true`, no cleanup occurs — containers, volumes, and networks remain on your machine.

  Feel free to use `ctf d rm` to remove containers when you are ready.
- If `TESTCONTAINERS_RYUK_DISABLED` is unset, the test environment will be automatically cleaned up a few seconds after the test completes.


Keep in mind that all components are mapped to [static ports](state.md), so without cleanup, only one environment can run at a time.

This design choice simplifies debugging.

A simplified command is available to prune unused volumes, containers, and build caches. Use it when you’re running low on space on your machine.
```
ctf d c
```

<div class="warning">

The framework manages cleanup for both on-chain and off-chain Docker components. However, if your test involves actions like configuring Chainlink jobs, it's best practice to make these actions idempotent, so they can be applied reliably in any environment.

</div>