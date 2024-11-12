## Framework

The primary focus of the Chainlink Testing Framework is to reduce the complexity of end-to-end testing, making complex system-level tests appear straightforward.
It enables tests to run in any environment and serves as a single source of truth for system behavior as defined by requirements.

### Features
- **Straightforward and sequential test composition**: Tests are readable and give you precise control over key aspects in a strict step-by-step order.

- **Modular configuration**: No arcane knowledge of framework settings is required; the config is simply a reflection of the components being used in the test. Components declare their own configuration— `what you see is what you get`.

- **Component isolation**: Components are decoupled via input/output structs, without exposing internal details.

- **Replaceability and extensibility**: Since components are decoupled via outputs, any deployment component can be swapped with a real service without altering the test code.

- **Quick local environments**: A common setup can be launched in just `15` seconds 🚀 [*](#cached).

- **Caching**: Any component can use cached configs to skip setup for even faster test development.

- **Integrated observability stack**: get all the info you need to develop end-to-end tests: metrics, logs, traces, profiles.

###### * If all the images are cached, you are using [OrbStack](https://orbstack.dev/) with M1/M2/M3 chips and have at least 8CPU dedicated to Docker

