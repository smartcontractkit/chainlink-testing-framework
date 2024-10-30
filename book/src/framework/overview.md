## Chainlink Testing Framework Harness

This module includes the CTFv2 harness, a lightweight, modular, and data-driven framework designed for combining off-chain and on-chain components while implementing best practices for end-to-end system-level testing:

- **Non-nil configuration**: All test variables must have defaults, automatic validation.
- **Component isolation**: Components are decoupled via input/output structs, without exposing internal details.
- **Modular configuration**: No arcane knowledge of framework settings is required; the config is simply a reflection of the components being used in the test. Components declare their own configuration—'what you see is what you get.'
- **Replaceability and extensibility**: Since components are decoupled via outputs, any deployment component can be swapped with a real service without altering the test code.
- **Caching**: any component can use cached configs to skip environment setup for faster test development
- **Integrated observability stack**: use `ctf obs up` to spin up a local observability stack.

