# CLL Private Test Examples

These examples require either internal images or access to our repositories.


## Job Distributor
Checkout Job Distributor repository and build an image
```
docker build -t job-distributor:0.9.0 -f e2e/Dockerfile.e2e .
```

## Test Environment Setup

Test environment includes:
- `Anvil` blockchain
- NodeSet with 5 nodes
- JobDistributor
- `gin` mockserver for fakes

This setup allows you to also set chain fork, see [config](jd_nodeset.toml), change `docker_cmd_params`.

You can also load JobDistributor database dump using [config](jd_nodeset.toml), see `jd_sql_dump_path` field

Run the tests locally
```
CTF_CONFIGS=jd_nodeset.toml go test -v -run TestJDNodeSet
```

## Connect to the JD database
```
PGPASSWORD=thispasswordislongenough psql -h 0.0.0.0 -p 14000 -U chainlink -d job-distributor-db
```