VERSION?=$(shell grep 'Version' ../cli/cmd/root.go | awk '{ print $$4 }' | head -n 1 | tr -d '"')
IMAGE_BASE="ghcr.io/zoomoid/assignments/ci"
BUILD_PLATFORMS?=linux/amd64
BUILD_ARGS?=--load

all: gitlab github

gitlab: # build gitlab container
	docker buildx build \
	-t $(IMAGE_BASE)/gitlab:$(VERSION) \
	--build-arg=BASE_IMAGE=assignments/cli:alpine \
	-f gitlab/Dockerfile \
	.

github: # build github container
	docker buildx build \
	-t $(IMAGE_BASE)/github:$(VERSION) \
	--build-arg=BASE_IMAGE=assignments/cli:alpine \
	-f github/Dockerfile \
	.

push: gitlab github
	docker push $(IMAGE_BASE)/github --all-tags
	docker push $(IMAGE_BASE)/gitlab --all-tags

version:
	@echo $(VERSION)

.PHONY: default gitlab github
