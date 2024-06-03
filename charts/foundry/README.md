## Anvil node for testing
See example [values.yaml](values.yaml) and [Dockerfile](Dockerfile) if you need to publish a custom image

Runs anvil node with an optional `--fork-url` and mines new block every `1s` by default

Change the URL in [values.yaml](values.yaml) and deploy

```
anvil:
  host: '0.0.0.0'
  port: '8545'
  blockTime: 1
  forkURL: 'https://goerli.infura.io/v3/...'
  forkBlockNumber: "10448829"
  forkRetries: "5"
  forkTimeout: "45000"
  forkComputeUnitsPerSecond: "330"
  # forkNoRateLimit: "true"
```

By default ingress is disabled, so remember to enable it in `values.yaml`.
Sample command:
```bash
export RELEASE_NAME="your-release-name"
export NAMESPACE="your-namespace"
export INGRESS_BASE_DOMAIN="your-ingress-base-domain"
export INGRESS_CERT_ARN="your-ingress-certificate"
export INGRESS_CIDRS="allowed-cidrs"

helm install "${RELEASE_NAME}" . -f ./values.yaml \
--set ingress.annotation_certificate_arn="${INGRESS_CERT_ARN}" \
--set "ingress.hosts[0].host"="${NAMESPACE}-anvil.${INGRESS_BASE_DOMAIN}" \
--set "ingress.annotation_group_name"="${NAMESPACE}" \
--set "ingress.enabled"=true \
--set "networkPolicyDefault.ingress.allowCustomCidrs"=true \
--set "networkPolicyDefault.ingress.customCidrs"="${INGRESS_CIDRS}"
# to override default chain id uncomment the following line
# --set "anvil.chainId"="2337"
```

You can find documentation for anvil [here](https://book.getfoundry.sh/reference/anvil/)

Also additional use cases can be found in the [testcases](../../client/rpc_test.go) folder:
1. (anvil) test drop transaction
2. (anvil) test we can shrink the block and control transaction inclusion
3. (anvil) test we can change next block base fee per gas and make tx pass or fail
4. (anvil) test we can mine sub-second blocks
5. (anvil) test we can mine blocks with strictly N+ transactions
