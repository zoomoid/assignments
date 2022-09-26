VERSION?=$(shell grep 'Version' ../cli/cmd/root.go | awk '{ print $$4 }' | head -n 1 | tr -d '"')
IMAGE_BASE="ghcr.io/zoomoid/assignments/ci"
BUILD_PLATFORMS?=linux/amd64
BUILD_ARGS?=--load

all: gitlab github

gitlab: # build gitlab container
	docker buildx build \
	-t $(IMAGE_BASE)/gitlab:$(VERSION) \
	-t $(IMAGE_BASE)/gitlab:latest \
	--build-arg=BASE_IMAGE=ghcr.io/zoomoid/assignments/cli:alpine \
	-f gitlab/Dockerfile \
	.

github: # build github container
	docker buildx build \
	-t $(IMAGE_BASE)/github:$(VERSION) \
	-t $(IMAGE_BASE)/github:latest \
	--build-arg=BASE_IMAGE=ghcr.io/zoomoid/assignments/cli:alpine \
	-f github/Dockerfile \
	.

push: gitlab github
	docker push $(IMAGE_BASE)/github:latest
	docker push $(IMAGE_BASE)/github:$(VERSION)
	docker push $(IMAGE_BASE)/gitlab:latest
	docker push $(IMAGE_BASE)/gitlab:$(VERSION)

version:
	@echo $(VERSION)

.PHONY: default gitlab github
