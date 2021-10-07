#!/usr/bin/env bash
GOOS=linux go build -o ./app . || exit 1;
docker build -t env-cleaner .
docker tag env-cleaner:latest "${REGISTRY}"/env-cleaner
docker push "${REGISTRY}"/env-cleaner
