# Raison d'etre

This chart allows to run a geth node as a non-root user, which is esential for running it on more secure clusters.

# Deploying with ingress

By default ingress is disabled. To enable it you need to override a couple of values in the values.yaml file. You can easily do it from command-line when installing the chart.

Currently ingress created for CRIB doesn't work, even though there are no errors or warnings in Kubernetes. Hopefuly soon we will have some eyes on it.

Sample command:
```bash
export RELEASE_NAME="your-release-name"
export NAMESPACE="your-namespace"
export INGRESS_BASE_DOMAIN="your-ingress-base-domain"
export INGRESS_CERT="your-ingress-certificate"
export INGRESS_CIDRS="your-ingress-cidrs"

helm install "${RELEASE_NAME}" . -f ./values.yaml \
--set ingress.annotation_certificate_arn="${INGRESS_CERT}"\
--set "ingress.hosts[0].host"="${NAMESPACE}-geth-http.${INGRESS_BASE_DOMAIN}"\
--set "ingress.hosts[1].host"="${NAMESPACE}-geth-ws.${INGRESS_BASE_DOMAIN}"\
--set "ingress.annotation_group_name"="${NAMESPACE}"\
--set "ingress.enabled"=true\
--set "networkPolicy.ingress.allowCustomCidrs=${INGRESS_CIDRS}"
```

# Limitations
Seems that Geth in dev mode doesn't fully support configuring `chainId`. Flag `--networkid` partially sets it, the setting fron `genesis.json` is ignored and result is a bit inconsistent:
```
INFO [02-06|19:28:20.759] Initialising Ethereum protocol           network=2337 dbversion=<nil>
INFO [02-06|19:28:20.759] Writing custom genesis block
INFO [02-06|19:28:20.760] Persisted trie from memory database      nodes=13 size=1.91KiB time="80.052µs" gcnodes=0 gcsize=0.00B gctime=0s livenodes=0 livesize=0.00B
INFO [02-06|19:28:20.760]
INFO [02-06|19:28:20.760] ---------------------------------------------------------------------------------------------------------------------------------------------------------
INFO [02-06|19:28:20.760] Chain ID:  1337 (unknown)
```

jsonRPC method `eth_chainId` returns `{"jsonrpc":"2.0","id":1,"result":"0x539"}`, which means that `chainId` is set to `1337` in hex.