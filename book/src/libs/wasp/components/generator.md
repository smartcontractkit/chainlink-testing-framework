# WASP - Generator

A **Generator** is a component that encapsulates all the characteristics of load generation, including:
* Load type:
  * RPS (requests per second)
  * VUs (virtual users)
* [Load schedule](./schedule.md)
* Call logic
* Response data
* [Sampling](./sampler.md)
* Timeouts

> [!WARNING]
> RPS load type can only be used with a `Gun`, while VUs can only be used with a `VirtualUser`.

---

### Choosing Between `Gun` and `VirtualUser`

#### **`Gun`**
* Best for **stateless protocols** (e.g., HTTP).
* Simplistic in nature; ideal for executing a single operation that does not require setup or teardown.
* Operates using an **open model**, where:
  * The number of requests is fixed.
  * The load adjusts to meet the target RPS, regardless of the system's response time.
  * There's no feedback from the system.
* Recommended for scenarios focused on measuring **throughput**.

#### **`VirtualUser`**
* Designed for **stateful protocols** (e.g., `WebSocket`) or workflows involving multiple operations (e.g., authenticating, executing tasks, and logging out).
* More complex, with dedicated methods for setup and teardown.
* Operates using a **closed model**, where:
  * New iterations start only after the previous one completes.
  * The RPS fluctuates based on the system's response time. Longer response times reduce RPS.
  * Feedback from the system is used to adjust the load.
---

### Closed vs. Open Models

* A `Gun` follows an **open model**:
  - It controls the rate of requests being sent.
  - The system's response time does not impact the load generation rate.

* A `VirtualUser` follows a **closed model**:
  - It controls the rate of receiving responses.
  - The system's response time directly impacts the load generation rate. If the system slows down, iterations take longer, reducing the RPS.

---

### Summary

In simpler terms:
* A **`Gun`** limits the load during the **sending** phase, making it ideal for throughput measurements.
* A **`VirtualUser`** limits the load during the **receiving** phase, reflecting the system's performance under load.

---

This distinction helps you decide which tool to use based on the protocol type and the goals of your test.
