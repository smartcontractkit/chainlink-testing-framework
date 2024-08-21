#!/bin/bash

set -o pipefail
set +e

# Check if required parameters are provided
if [ "$#" -ne 6 ]; then
    echo "Usage: $0 DOCKERFILE_PATH TESTS_ROOT_PATH IMAGE_TAG ECR_REGISTRY_NAME ECR_REGISTRY_REPO_NAME DOCKER_CMD_EXEC_PATH"
    exit 1
fi

DOCKERFILE_PATH="$1"
TESTS_ROOT_PATH="$2"
IMAGE_TAG="$3"
ECR_REGISTRY_NAME="$4"
ECR_REGISTRY_REPO_NAME="$5"
DOCKER_CMD_EXEC_PATH="$6"

# Build Docker image
cd "$DOCKER_CMD_EXEC_PATH" && docker build --platform linux/amd64 -t "$IMAGE_TAG" -f "$DOCKERFILE_PATH" --build-arg TESTS_ROOT="$TESTS_ROOT_PATH" .

# Authenticate Docker with ECR
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin "$ECR_REGISTRY_NAME"

# Tag the Docker image with ECR registry name
docker tag "$IMAGE_TAG" "$ECR_REGISTRY_NAME/$ECR_REGISTRY_REPO_NAME:$IMAGE_TAG"

# Push Docker image to ECR
docker push "$ECR_REGISTRY_NAME/$ECR_REGISTRY_REPO_NAME:$IMAGE_TAG"

# Verify push success
# shellcheck disable=SC2181
if [ $? -eq 0 ]; then
    echo "Image successfully pushed to ECR."
else
    echo "Failed to push image to ECR."
    exit 1
fi
