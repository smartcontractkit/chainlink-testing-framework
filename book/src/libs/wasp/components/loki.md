# WASP - Loki

Loki is a component responsible for pushing batches of statistics to a Loki instance. It operates in the background using a `promtail` client.

Key features include:
* Optional basic authentication support.
* Configurable batch sizes.
* Support for configurable backoff retries, among others.

By default, a test will fail on the first error encountered while pushing data. You can modify this behavior by setting the maximum allowed errors:
* Setting the value to `-1` disables error checks entirely.

> [!NOTE]  
> While it is technically possible to execute a WASP test without Loki, doing so means you won't have access to load-generation-related metrics.  
> Unless your sole interest is in the metrics sent by your application, it is highly recommended to use Loki.  
> For this reason, Loki is considered an integral part of the WASP stack.
