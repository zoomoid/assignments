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

build:
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.VERSION=$(VERSION)" -o ./bin/assignmentctl main.go
