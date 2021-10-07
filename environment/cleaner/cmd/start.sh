#!/usr/bin/env bash
kubectl run -i env-cleaner --image="${REGISTRY}"/env-cleaner