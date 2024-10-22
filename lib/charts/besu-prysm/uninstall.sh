#!/bin/bash
# shellcheck disable=SC2162
read -p "Enter the Helm release name to uninstall (press Enter to uninstall the first one found, if none is given): " release_name

if [ -z "$release_name" ]; then
  # shellcheck disable=SC2162
  read -p "No release name provided. Are you sure you want to uninstall the first Helm deployment? (y/n): " confirm_uninstall
  if [ "$confirm_uninstall" != "y" ]; then
    echo "Aborted uninstallation."
    exit 0
  fi
  release_name=$(helm list -o json | jq -r '.[0].name')
fi

helm uninstall "$release_name"

echo "Deleting all PVCs"

pvcs=$(kubectl get pvc --no-headers -o custom-columns=":metadata.name")

for pvc in $pvcs; do
    echo "Deleting PVC: $pvc"
    kubectl delete pvc "$pvc"
done

echo "All PVCs have been deleted."
