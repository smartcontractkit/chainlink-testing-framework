#!/bin/sh
values_file="Values.yaml"

if ! helm lint . -f ./$values_file; then
  echo "Helm lint failed. Exiting."
  exit 1
fi

helm package .

chart_package=$(ls -1 | grep '.tgz')

if [ -z "$chart_package" ]; then
  echo "No Helm chart package found."
  exit 1
fi

deployment_name=$1
if [ -z "$deployment_name" ]; then
  deployment_name="--generate-name"
fi

helm install "$deployment_name" "$chart_package" -f ./$values_file
