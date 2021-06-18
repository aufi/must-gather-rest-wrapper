# Image URL to use all building/pushing image targets
IMG ?= quay.io/maufart/must-gather-rest-wrapper:latest
GOOS ?= `go env GOOS`
GOBIN ?= ${GOPATH}/bin
GO111MODULE = auto

ci: all

all: test manager

# Run tests
test: fmt vet
	go test ./pkg/... -coverprofile cover.out

# Build manager binary
manager: fmt vet
	go build -o bin/app github.com/aufi/must-gather-rest-wrapper/pkg/must-gather-rest-wrapper

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	go run ./pkg/must-gather-rest-wrapper/main.go

# Run go fmt against code
fmt:
	go fmt ./pkg/...

# Run go vet against code
vet:
	go vet ./pkg/...

# Build the docker image
#docker-build: test
docker-build:
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}
