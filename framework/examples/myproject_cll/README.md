# CLL Private Test Examples

These examples require either internal images or access to our repositories.


## Job Distributor
Checkout Job Distributor repository and build an image
```
docker build -t job-distributor:0.9.0 -f e2e/Dockerfile.e2e .
```

Run the tests locally
```
CTF_CONFIGS=jd.toml go test -v -run TestJD
```

## Job Distributor from staging dump
Get `.sql` dump from staging or production service and then run
```
CTF_CONFIGS=jd_dump.toml go test -v run TestJD
```

## Connect to the database
```
PGPASSWORD=thispasswordislongenough psql -h 0.0.0.0 -p 14000 -U chainlink -d job-distributor-db
```