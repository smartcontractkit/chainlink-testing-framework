#!/bin/bash

# Run helm list and grab the name of the first deployment
deployment_name=$(helm list -o json | jq -r '.[0].name')

if [ -z "$deployment_name" ]; then
  echo "No Helm deployments found."
  exit 1
fi

# Uninstall the first deployment
helm uninstall "$deployment_name"

PVC_NAME="chain-state-claim"
PV_NAME="chain-data-storage"

# Delete PersistentVolumeClaim (PVC)
kubectl delete persistentvolumeclaim "$PVC_NAME"
if [ $? -ne 0 ]; then
  echo "Failed to delete PVC: $PVC_NAME"
else
  echo "PVC deleted successfully: $PVC_NAME"
fi

# Delete PersistentVolume (PV)
kubectl delete persistentvolume "$PV_NAME"
if [ $? -ne 0 ]; then
  echo "Failed to delete PV: $PV_NAME"
else
  echo "PV deleted successfully: $PV_NAME"
fi

echo "PVC and PV deletion completed successfully."

# Run helm lint
if ! helm lint . -f Values.yaml; then
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

# Apply the persistent volume configuration
if kubectl apply -f pv/pv-shared-storage.yaml; then
  echo "Persistent volume configuration applied successfully."
else
  echo "Error: Failed to apply persistent volume configuration."
  exit 1
fi

# Create the persistent volume claim
if kubectl create -f pvc/pvc-chain-state-claim.yaml; then
  echo "Persistent volume claim created successfully."
else
  echo "Error: Failed to create persistent volume claim."
  exit 1
fi

# Install the newly generated chart package
helm install "$chart_package" --generate-name -f Values.yaml

