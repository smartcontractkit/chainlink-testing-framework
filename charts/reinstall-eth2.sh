#!/bin/bash

RELEASE_NAME="eth2-prysm-geth"

# Run helm list and grab the name of the first deployment
deployment_name=$(helm list -o json | jq -r '.[0].name')

if [ -z "$deployment_name" ]; then
  echo "No Helm deployments found."
  exit 1
fi

# Uninstall the first deployment
helm uninstall "$deployment_name"

PVC_NAME="chain-state-claim"
PV_NAME="chain-state-storage"

kubectl delete pvc $PVC_NAME
kubectl delete pv $PV_NAME

# Run helm lint
if ! helm lint ./geth-prysm -f ./geth-prysm/Values.yaml; then
  echo "Helm lint failed. Exiting."
  exit 1
fi

# Package the Helm chart in the current directory
helm package ./geth-prysm

# Get the name of the generated chart package
chart_package=$(ls -1 | grep '.tgz')

if [ -z "$chart_package" ]; then
  echo "No Helm chart package found."
  exit 1
fi

# Install the newly generated chart package
helm install $RELEASE_NAME "$chart_package" -f ./geth-prysm/Values.yaml

