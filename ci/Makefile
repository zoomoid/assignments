VERSION?=$(shell grep 'Version' ../cli/cmd/root.go | awk '{ print $$4 }' | head -n 1 | tr -d '"')
IMAGE_BASE="ghcr.io/zoomoid/assignments/ci"
BUILD_PLATFORMS?=linux/amd64
BUILD_ARGS?=--load

.PHONY: default gitlab github

all: gitlab github

gitlab: # build gitlab container
	docker buildx build \
	-t $(IMAGE_BASE)/gitlab:$(VERSION) \
	-f gitlab/Dockerfile \
	.

github: # build github container
	docker buildx build \ 
	-t $(IMAGE_BASE)/github:$(VERSION) \
	-f github/Dockerfile \
	.

version:
	@echo $(VERSION)