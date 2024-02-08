# Raison d'etre

This chart allows to run a geth node as a non-root user, which is esential for running it on more secure clusters. Geth is running as Proof-of-Authority private network with a single node.

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
# to override default chain id uncomment the following line
# --set "geth.networkId"="2337"
--set "networkPolicy.ingress.allowCustomCidrs=${INGRESS_CIDRS}"
```