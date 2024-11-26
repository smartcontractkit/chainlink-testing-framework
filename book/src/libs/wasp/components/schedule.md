# WASP - Schedule

The **Schedule** component allows you to define the load characteristics for each [Generator](./generator.md).  
A schedule is composed of various `Segments`, each characterized by the number of requests (or VUs) and duration.

---

### Helper Functions for Defining Schedules

WASP provides several helper functions to define schedules in a human-readable way:

* **`Plain`**: Defines a segment with a stable load.
* **`Steps`**: Defines a segment with increasing or decreasing load, using a specified step size.

---

### Custom Schedules

You can also define schedules programmatically by directly composing segments to suit your specific requirements.
