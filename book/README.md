# Chainlink Testing Framework Docs

We use [mdBook](https://github.com/rust-lang/mdBook) for our docs. They can be found [here](https://smartcontractkit.github.io/chainlink-testing-framework/).

## Development

First [install Rust](https://doc.rust-lang.org/cargo/getting-started/installation.html), then [mdbook](https://github.com/rust-lang/mdBook).

```sh
# Install mdBook
cargo install mdbook && \
cargo install mdbook-alerts && \
cargo install mdbook-cmdrun
# Run the mdBook
make run
```

Visit [localhost:9999](http://localhost:9999) to see your local version.
