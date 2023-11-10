#!/bin/bash

# Run helm list and grab the name of the first deployment
deployment_name=$(helm list -o json | jq -r '.[0].name')

if [ -z "$deployment_name" ]; then
  echo "No Helm deployments found."
  exit 1
fi

# Uninstall the first deployment
helm uninstall "$deployment_name"
