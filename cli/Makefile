.ONESHELL:

EXECUTABLES = git go find pwd
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
IMAGE_BASE="ghcr.io/zoomoid/assignments"
BINARY=assignmentctl
VERSION?=$(shell grep 'Version' cmd/root.go | awk '{ print $$4 }' | head -n 1 | tr -d '"')
BUILD=`git rev-parse HEAD | head -c 8`
WINDOWS_PLATFORMS=windows/amd64 windows/386 windows/arm windows/arm64
UNIX_PLATFORMS=linux/amd64 linux/arm linux/arm64 linux/386 darwin/amd64
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

ARTIFACTS_DIR=bin

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-s -w -X main.Version=${VERSION}"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: test build

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

install-dev:
	CGO_ENABLED=0 go build $(LDFLAGS) -o /usr/local/bin/$(BINARY) main.go

build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o ./bin/$(BINARY) main.go

release: $(UNIX_PLATFORMS) $(WINDOWS_PLATFORMS)

clean:
	rm -rf bin/

$(UNIX_PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -v -o $(ARTIFACTS_DIR)/$(BINARY)_$(os)_$(arch)/assignmentctl
	cd $(ARTIFACTS_DIR)/$(BINARY)_$(os)_$(arch)/
	tar -czf $(BINARY)_$(VERSION)_$(os)_$(arch).tar.gz assignmentctl
	mv $(BINARY)_$(VERSION)_$(os)_$(arch).tar.gz ../
	cd ../
	rm -rf $(BINARY)_$(os)_$(arch)/

$(WINDOWS_PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -v -o $(ARTIFACTS_DIR)/$(BINARY)_$(os)_$(arch)/assignmentctl.exe
	cd $(ARTIFACTS_DIR)/$(BINARY)_$(os)_$(arch)/
	zip $(BINARY)_$(VERSION)_$(os)_$(arch).zip assignmentctl.exe
	mv $(BINARY)_$(VERSION)_$(os)_$(arch).zip ../
	cd ../
	rm -rf $(BINARY)_$(os)_$(arch)/

.PHONY: clean release build

containers: ubuntu-cli-container alpine-cli-container runner-container

ubuntu-cli-container:
	docker buildx build \
	-t $(IMAGE_BASE)/cli:$(VERSION) \
	-t $(IMAGE_BASE)/cli:latest \
	-t assignments/cli:latest \
	-f build/cli/ubuntu/Dockerfile \
	.

alpine-cli-container:
	docker buildx build \
	-t $(IMAGE_BASE)/cli:alpine-$(VERSION) \
	-t $(IMAGE_BASE)/cli:alpine \
	-t assignments/cli:alpine \
	-f build/cli/alpine/Dockerfile \
	.

runner-container: ubuntu-cli-container
	docker buildx build \
	-t $(IMAGE_BASE)/runner:$(VERSION) \
	-t $(IMAGE_BASE)/runner:latest \
	-t assignments/runner:latest \
	-f build/runner/Dockerfile \
	--build-arg=IMAGE=assignments/cli:latest \
	build/runner/

