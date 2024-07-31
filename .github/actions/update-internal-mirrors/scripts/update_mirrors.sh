#!/bin/bash

set -e

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

cd "$SCRIPT_DIR"/../ || exit 1

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

    repository_name=$(echo "$docker_image" | cut -d: -f1)
    image_tag=$(echo "$docker_image" | cut -d: -f2)

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

    docker pull "$docker_image"
    docker tag "$docker_image" "$ecr_image"
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
    elif [[ $image_name == ghcr.io* ]]; then
        # Handle GitHub Container Registry images
        images=$(fetch_images_from_gh_container_registry "${image_name}" "${image_expression}" "${page_size}")
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

# Function to fetch images from Github Container Registry with pagination support
fetch_images_from_gh_container_registry() {
        local image_name="$1"
        local image_expression="$2"
        local max_image_count="$3"

        local org
        local package

        org=$(echo "$image_name" | awk -F'[/:]' '{print $2}')
        package=$(echo "$image_name" | awk -F'[/:]' '{print $3}')

        if [ -z "$org" ] || [ -z "$package" ]; then
            >&2 echo "Error: Failed to extract organisation and package name from $image_name. Please provide the image name in the format ghcr.io/org/package."
            exit 1
        fi

        if [ -z "$GHCR_TOKEN" ]; then
            >&2 echo "Error: $GHCR_TOKEN environment variable is not set."
            exit 1
        else
            >&2 echo "::debug::GHCR_TOKEN is set"
        fi

        local url="https://api.github.com/orgs/$org/packages?package_type=container"
        >&2 echo "::debug::url: $url"

        local image_count=0
        local images=""

        while [ -n "$url" ]; do
            response=$(curl -s -H "Authorization: Bearer $GHCR_TOKEN" \
                              -H "Accept: application/vnd.github.v3+json" \
                              "$url")

             >&2 echo "::debug::response: $response"

            if ! echo "$response" | jq empty > /dev/null 2>&1; then
                 >&2 echo "Error: Received invalid JSON response."
                exit 1
            fi

            if echo "$response" | jq -e 'if type == "object" then (has("message") or has("status")) else false end' > /dev/null; then
                message=$(echo "$response" | jq -r '.message // empty')
                status=$(echo "$response" | jq -r '.status // empty')

                if [ -n "$status" ] && [ "$status" -eq "$status" ] 2>/dev/null && [ "$status" -gt 299 ]; then
                     >&2 echo "Error: Request to get containers failed with status $status and message: $message"
                    exit 1
                fi
            fi

            packages=$(echo "$response" | jq -r --arg package "$package" '.[] | select(.name == $package) | .name')

            if [ -z "$packages" ]; then
                 >&2 echo "Error: No matching packages found."
                exit 1
            fi

            for package in $packages; do
                versions_url="https://api.github.com/orgs/$org/packages/container/$package/versions"
                while [ -n "$versions_url" ]; do
                    versions_response=$(curl -s -H "Authorization: token $GHCR_TOKEN" \
                                             -H "Accept: application/vnd.github.v3+json" \
                                             "$versions_url")

                    if ! echo "$versions_response" | jq empty > /dev/null 2>&1; then
                         >&2 echo "Error: Received invalid JSON response for versions."
                        exit 1
                    fi

                    tags=$(echo "$versions_response" | jq -r --arg regex "$image_expression" '
                        .[] |
                        select(.metadata.container.tags | length > 0) |
                        .metadata.container.tags[] as $tag |
                        select($tag | test($regex)) |
                        $tag
                    ')

                    while read -r tag; do
                        if [ "$image_count" -lt "$max_image_count" ]; then
                            images+="$tag"$'\n'
                            ((image_count++))
                        else
                            break 2
                        fi
                    done <<< "$tags"

                    if [ "$image_count" -ge "$max_image_count" ]; then
                        images=$(echo "$images" | grep -v '^\s*$')
                        echo "$images"
                        return
                    fi

                    versions_url=$(curl -sI -H "Authorization: token $GHCR_TOKEN" \
                                            -H "Accept: application/vnd.github.v3+json" \
                                            "$versions_url" | awk -F'[<>]' '/rel="next"/{print $2}')
                done
            done

            url=$(curl -sI -H "Authorization: token $GHCR_TOKEN" \
                        -H "Accept: application/vnd.github.v3+json" \
                        "$url" | awk -F'[<>]' '/rel="next"/{print $2}')
        done

        images=$(echo "$images" | grep -v '^\s*$')
        echo "$images"
}

push_images_in_list() {
    local -a image_list=("$@")
    local prefix="library/"

    for docker_image in "${image_list[@]}"; do
        echo "---"
        echo "Checking if $docker_image exists in ECR..."

        docker_image="${docker_image#"$prefix"}"

        if ! check_image_in_ecr "$docker_image"; then
            echo "$docker_image does not exist in ECR. Mirroring image..."

            pull_tag_push "$docker_image"
        else
            echo "$docker_image already exists in ECR. Skipping..."
        fi
    done
}

# Run the code
if [ $# -eq 1 ]; then
    push_images_in_mirror_json
    echo "Update from mirror.json completed."
elif [ $# -eq 3 ] || [ $# -eq 4 ]; then
    push_latest_images_for_expression_from_dockerhub "$2" "$3" "$4"
    echo "Update from dockerhub completed."
else
    usage
fi
