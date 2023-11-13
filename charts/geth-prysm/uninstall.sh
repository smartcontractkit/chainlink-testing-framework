#!/bin/bash

# Ask for release name to uninstall
read -p "Enter the Helm release name to uninstall (press Enter to uninstall the first one found, if none is given): " release_name

if [ -z "$release_name" ]; then
  read -p "No release name provided. Are you sure you want to uninstall the first Helm deployment? (y/n): " confirm_uninstall
  if [ "$confirm_uninstall" != "y" ]; then
    echo "Aborted uninstallation."
    exit 0
  fi
  # Run helm list and grab the name of the first deployment
  release_name=$(helm list -o json | jq -r '.[0].name')
fi

# Uninstall the specified release
helm uninstall "$release_name"