#!/bin/bash

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

cd "$SCRIPT_DIR"/../ || exit 1

# Check if an argument is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <ecr-registry-url>"
    exit 1
fi

# The first argument is the ECR registry URL
ECR_REGISTRY_URL="$1"

# Function to check if the image exists in ECR
check_image_in_ecr() {
    local docker_image="$1"
    local repository_name image_tag

    # Extract the repository name and tag from the docker image string
    repository_name=$(echo "$docker_image" | cut -d: -f1)
    image_tag=$(echo "$docker_image" | cut -d: -f2)

    # If the image tag is empty, it means the image name did not include a tag, and we'll use "latest" by default
    if [[ -z "$image_tag" ]]; then
        image_tag="latest"
    fi

    # The repository name in ECR is just the part before the colon, we don't need to replace it with an underscore
    if aws ecr describe-images --repository-name "$repository_name" --image-ids imageTag="$image_tag" > /dev/null 2>&1; then
        return 0 # Image exists
    else
        return 1 # Image does not exist
    fi
}

# Function to pull, tag, and push the image to ECR
pull_tag_push() {
    local docker_image=$1
    local ecr_image="$ECR_REGISTRY_URL/${docker_image}"
    
    # Pull the image from Docker Hub
    docker pull "$docker_image"
    
    # Tag the image for ECR
    docker tag "$docker_image" "$ecr_image"
    
    # Push the image to ECR
    docker push "$ecr_image"
}

# Read the JSON file into a bash array
docker_images=()
while IFS= read -r line; do
    docker_images+=("$line")
done < <(jq -r '.[]' ./mirror/mirror.json)

# Iterate over the images
for docker_image in "${docker_images[@]}"; do
    echo "---"
    echo "Checking if $docker_image exists in ECR..."

    # Check if the image exists in ECR
    if ! check_image_in_ecr "$docker_image"; then
        echo "$docker_image does not exist in ECR. Mirroring image..."
        # Pull, tag, and push the image to ECR
        pull_tag_push "$docker_image"
    else
        echo "$docker_image already exists in ECR. Skipping..."
    fi
done

echo "Mirroring process completed."