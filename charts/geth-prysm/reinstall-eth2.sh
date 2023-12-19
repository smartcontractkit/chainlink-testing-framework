#!/bin/bash
read -p "Enter the Helm release name to uninstall (press Enter to uninstall the first one found, if none is given): " release_name

if [ -z "$release_name" ]; then
  read -p "No release name provided. Are you sure you want to uninstall the first Helm deployment? (y/n): " confirm_uninstall
  if [ "$confirm_uninstall" != "y" ]; then
    echo "Aborted uninstallation."
    exit 0
  fi
  release_name=$(helm list -o json | jq -r '.[0].name')
fi

helm uninstall "$release_name"

param=$1
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

deployment_name=$2
if [ -z "$deployment_name" ]; then
  deployment_name="--generate-name"
fi

helm install "$deployment_name" "$chart_package" -f ./$values_file