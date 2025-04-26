# CLL Private Test Examples

These examples require either internal images or access to our repositories.


## Job Distributor
Checkout Job Distributor repository and build an image
```
docker build -t job-distributor:0.9.0 -f e2e/Dockerfile.e2e .
```

Run the tests locally
```
CTF_CONFIGS=jd.toml go test -v -count 1 -run TestJD
CTF_CONFIGS=jd_fork.toml go test -v -count 1 -run TestJDFork
```