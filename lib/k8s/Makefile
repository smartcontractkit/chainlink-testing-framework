.PHONY: build
build:
	docker build --platform linux/amd64 -t test-base-image -f Dockerfile.base .

.PHONY: push
push:
	@if [ -z "$(registry)" ]; then echo "Error: registry variable is not set. Please provide it with 'make push registry=...'. See README.md explaining where to get it."; exit 1; fi
	aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin $(registry).dkr.ecr.us-west-2.amazonaws.com
	docker tag test-base-image:latest $(registry).dkr.ecr.us-west-2.amazonaws.com/test-base-image:latest
	docker push $(registry).dkr.ecr.us-west-2.amazonaws.com/test-base-image:latest
