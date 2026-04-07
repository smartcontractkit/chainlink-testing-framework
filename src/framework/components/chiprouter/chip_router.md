# Chip Router

`chiprouter` is a small CTF component that owns the fixed ChIP ingress port and fans incoming telemetry out to registered downstream subscribers.

It exists to keep the local CRE topology simple:
- Chainlink nodes always publish to a single ingress owner on `50051`
- lightweight test sinks subscribe behind the router
- real ChIP / Beholder subscribes behind the same router

That removes the old split where some tests bound ingress directly while others started real ChIP.

## Ports

The component exposes:
- admin HTTP: `50050`
- ingress gRPC: `50051`

In the local CRE topology, real ChIP / Beholder typically subscribes downstream on `50053`.

## Image Contract

The component runs whatever image is provided in `chip_router.image`.

The expected local CRE convention is:
- env TOMLs use a local alias such as `chip-router:<commit-sha>`
- setup/pull logic is responsible for making that alias exist locally
- remote ECR image names stay in setup/pull config and are retagged locally to the alias

## Runtime Behavior

The router:
- exposes a health endpoint on `/health`
- accepts subscriber registration over its admin API
- forwards published ChIP ingress requests to all registered subscribers
- is best-effort per subscriber, so one failing downstream does not block others

Host-based downstream subscribers should register host-reachable endpoints. In local CRE, host-local sink endpoints are normalized to the Docker host gateway before registration.
