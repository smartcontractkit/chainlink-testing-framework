.PHONY: build-base
build:
	docker build --platform linux/amd64 -t test-base-image -f Dockerfile.base .

push:
	aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin $(registry).dkr.ecr.us-west-2.amazonaws.com
	docker tag test-base-image:latest $(registry).dkr.ecr.us-west-2.amazonaws.com/test-base-image:latest
	docker push $(registry).dkr.ecr.us-west-2.amazonaws.com/test-base-image:latest
