## Framework

The Framework module a lightweight, modular, and data-driven framework designed for combining off-chain and on-chain components while implementing best practices for end-to-end system-level testing:

- **Modular configuration**: No arcane knowledge of framework settings is required; the config is simply a reflection of the components being used in the test. Components declare their own configuration—'what you see is what you get.'
- **Component isolation**: Components are decoupled via input/output structs, without exposing internal details.
- **Replaceability and extensibility**: Since components are decoupled via outputs, any deployment component can be swapped with a real service without altering the test code.
- **Caching**: any component can use cached configs to skip setup for faster test development.
- **Integrated observability stack**: get all the info you need to develop end-to-end tests: metrics, logs, traces, profiles.

