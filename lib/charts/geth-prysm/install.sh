#!/bin/sh
values_file="Values.yaml"

if ! helm lint . -f ./$values_file; then
  echo "Helm lint failed. Exiting."
  exit 1
fi

helm package .

# shellcheck disable=SC2010
chart_package=$(ls -1 | grep '.tgz')

if [ -z "$chart_package" ]; then
  echo "No Helm chart package found."
  exit 1
fi

deployment_name=$1
if [ -z "$deployment_name" ]; then
  deployment_name="--generate-name"
fi

now=$(date +%s)
# shellcheck disable=SC2086
helm install "$deployment_name" "$chart_package" -f ./$values_file --set "genesis.values.currentUnixTimestamp"="$now" --set "eth2-common.genesis.values.currentUnixTimestamp"="$now" $2
