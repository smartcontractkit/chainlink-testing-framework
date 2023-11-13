#!/bin/sh
# Delete PVC and PV
PVC_NAME="chain-state-claim"
PV_NAME="chain-state-storage"
kubectl delete pvc $PVC_NAME
kubectl delete pv $PV_NAME

param=$1

if [ -n "$param" ]; then
  values_file="Values-$param.yaml"
else
  values_file="Values.yaml"
fi

# Run helm lint
if ! helm lint . -f ./$values_file; then
  echo "Helm lint failed. Exiting."
  exit 1
fi

# Package the Helm chart in the current directory
helm package .

# Get the name of the generated chart package
chart_package=$(ls -1 | grep '.tgz')

if [ -z "$chart_package" ]; then
  echo "No Helm chart package found."
  exit 1
fi

deployment_name=$2
if [ -z "$deployment_name" ]; then
  # If deployment name is empty, use --generate-name
  deployment_name="--generate-name"
fi

# Install the newly generated chart package
helm install "$deployment_name" "$chart_package" -f ./$values_file
