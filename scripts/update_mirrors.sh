#!/bin/bash

set -e

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

cd "$SCRIPT_DIR"/../ || exit 1

# Check if any arguments are provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <ecr-registry-url>"
    exit 1
fi

# The first argument is the ECR registry URL
ECR_REGISTRY_URL="$1"

usage() {
    echo "Need either 1 or 3 arguments depending on which operation mode you want"
    echo "Usage: $0 <ecr-registry-url> [image-name] [image-expression]"
    echo "Only add registry url to update from mirror.json file"
    echo "Add both image expression and image name to update from list in dockerhub."
    echo "expression example for something like postgres: '^15\.[0-9]+$'"
    echo "Optional argument 4 is the number of images to check in dockerhub, useful for images like 50 that have lots of images between official version tags."
    exit 1

}

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

push_images_in_mirror_json() {
    # Read the JSON file into a bash array
    docker_images=()
    while IFS= read -r line; do
        docker_images+=("$line")
    done < <(jq -r '.[]' ./scripts/mirror.json)

    push_images_in_list "${docker_images[@]}"
}

push_latest_images_for_expression_from_dockerhub() {
    local image_name=$1
    local image_expression=$2
    local page_size=$3
    local images

    # check if we have a number for the page size, if not, set it to 50
    if [ -z "$page_size" ] || ! [[ "$page_size" =~ ^[0-9]+$ ]] || [ "$page_size" -eq 0 ]; then
        page_size=50
    fi

    set +e
    if [[ $image_name == gcr.io* ]]; then
        # Handle GCR images
        images=$(gcloud container images list-tags gcr.io/prysmaticlabs/prysm/validator --limit="${page_size}" --filter='tags:v*' --format=json | jq -r '.[].tags[]' | grep -E "${image_expression}")
    else
        images=$(fetch_images_from_dockerhub "${image_name}" "${image_expression}" "${page_size}")
    fi
    set -e

    if [ -z "$images" ]; then
        echo "No images were found matching the expression. Either something went wrong or you need to increase the page size to greater than ${page_size}."
        exit 1
    else
        echo "Images found:"
        echo "$images"
    fi

    image_list=()

    # Convert newline-separated string to an array
    while IFS= read -r line; do
        image_list+=("${image_name}:${line}")
    done <<< "$images"
    push_images_in_list "${image_list[@]}"
}

# Function to fetch images from Docker Hub with pagination support
fetch_images_from_dockerhub() {
    local image_name=$1
    local image_expression=$2
    local desired_count=$3
    local max_page_size=100  # Fixed maximum page size
    local loop_count=0
    local max_loops=$(( (desired_count + max_page_size - 1) / max_page_size ))  # Calculate max loops needed
    local images=""
    local next_url="https://hub.docker.com/v2/repositories/${image_name}/tags/?page_size=${max_page_size}"

    while [ "$next_url" != "null" ] && [ "$loop_count" -lt "$max_loops" ]; do
        response=$(curl -s "$next_url")
        new_images=$(echo "$response" | jq -r '.results[].name' | grep -E "${image_expression}" | grep -v '^\s*$')
        if [ -n "$new_images" ]; then
            images+="$new_images"$'\n'
        fi
        next_url=$(echo "$response" | jq -r '.next')
        loop_count=$((loop_count + 1))
    done

    # Ensure no trailing newlines
    images=$(echo "$images" | grep -v '^\s*$')
    echo "$images"
}

push_images_in_list() {
    local -a image_list=("$@")
    local prefix="library/"
    # Iterate over the images
    for docker_image in "${image_list[@]}"; do
        echo "---"
        echo "Checking if $docker_image exists in ECR..."

        # Check if the image is a standard libary image and needs the library/ prefix removed
        docker_image="${docker_image#"$prefix"}"

        # Check if the image exists in ECR
        if ! check_image_in_ecr "$docker_image"; then
            echo "$docker_image does not exist in ECR. Mirroring image..."
            # Pull, tag, and push the image to ECR
            pull_tag_push "$docker_image"
        else
            echo "$docker_image already exists in ECR. Skipping..."
        fi
    done

}

# Run the code
if [ $# -eq 1 ]; then
    push_images_in_mirror_json
    echo "Update from mirror.json comleted."
elif [ $# -eq 3 ] || [ $# -eq 4 ]; then
    push_latest_images_for_expression_from_dockerhub "$2" "$3" "$4"
    echo "Update from dockerhub completed."
else
    usage
fi
