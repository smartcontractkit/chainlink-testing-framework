# Troubleshooting Guide

This page lists common issues, error messages, and solutions for the Chainlink Testing Framework (CTF).

---

## 1. Docker Issues

**Problem:** Docker containers fail to start or are slow.
- **Solution:**
  - Ensure Docker is running and has enough resources (CPU, RAM).
  - Use OrbStack for better performance on macOS.
  - Run `ctf d rm` to clean up old containers.
  - Check for port conflicts with `lsof -i :PORT`.

---

## 2. Configuration Errors

**Problem:** `missing required field` or `validation failed`.
- **Solution:**
  - Check your TOML config for typos and required fields.
  - Use `ctf validate config.toml` to validate your config.
  - Ensure all `validate:"required"` fields are present in your Go struct.

---

## 3. Caching Problems

**Problem:** Component not updating or stale state.
- **Solution:**
  - Set `use_cache = false` in your TOML or `export CTF_DISABLE_CACHE=true`.
  - Delete the `.ctf_cache` directory and rerun tests.

---

## 4. Observability Issues

**Problem:** No logs or metrics in Grafana.
- **Solution:**
  - Check that Loki and Prometheus containers are running (`docker ps`).
  - Ensure your components are configured to send logs/metrics.
  - Check Grafana data source settings.

---

## 5. Test Failures

**Problem:** Tests fail with `connection refused`, `timeout`, or `not ready` errors.
- **Solution:**
  - Increase test timeouts (`go test -timeout 10m`).
  - Use `require.Eventually` to wait for readiness.
  - Check Docker resource limits and logs.

---

## 6. CI/CD Integration

**Problem:** Tests pass locally but fail in CI.
- **Solution:**
  - Ensure all environment variables are set in CI.
  - Use `CTF_DISABLE_CACHE=true` for clean runs.
  - Check CI logs for missing dependencies or permissions.

---

## 7. Miscellaneous

- **Problem:** `permission denied` errors on files or Docker.
  - **Solution:** Run with appropriate permissions or fix file ownership.
- **Problem:** `address already in use`.
  - **Solution:** Free the port or change the config.

---

## Getting Help
- Check the [README](https://github.com/smartcontractkit/chainlink-testing-framework#readme)
- Search [GitHub Issues](https://github.com/smartcontractkit/chainlink-testing-framework/issues)
- Ask in [Discussions](https://github.com/smartcontractkit/chainlink-testing-framework/discussions)
- Create a new issue with detailed logs and config 