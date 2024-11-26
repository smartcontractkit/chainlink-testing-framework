# WASP - How to Debug Loki Push Errors

To troubleshoot Loki push errors, follow these steps:

---

### Step 1: Enable Trace Logging

Set the logging level to `trace` to see all messages sent to Loki:

```bash
WASP_LOG_LEVEL=trace
```

---

### Step 2: Adjust Error Limits

If the Loki client fails to deliver a batch, the test will continue until the maximum number of errors is reached.  
You can control this using `LokiConfig.MaxErrors`:
* Set `LokiConfig.MaxErrors` to a desired limit.
* Set it to `-1` to disable error checks entirely.

Often, simply increasing `LokiConfig.Timeout` can resolve the issue.

---

### Step 3: Handle Specific Errors

If you encounter errors like the following:

```
ERR Malformed promtail log message, skipping Line=["level",{},"component","client","host","...","msg","batch add err","tenant","","error",{}]
```

Take the following actions:
1. Increase `LokiConfig.MaxStreams`.
2. Verify the validity of your Loki configuration.

---

Since `LokiConfig` is specific for each `Generator` remember to adjust the configuration for all of your generators.

By following these steps, you can effectively debug and resolve Loki push errors.