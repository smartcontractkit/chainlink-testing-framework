#!/usr/bin/env bash
docker build -f load_test.Dockerfile --build-arg BUILD_ROOT=$1 --platform linux/amd64 -t $2 .
docker tag wasp-mercury-test $1/wasp-mercury-test:$4
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin $1
docker push $1/wasp-mercury-test:$4