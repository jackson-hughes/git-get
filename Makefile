GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GIT_COMMIT=$(shell git rev-list -1 HEAD)
GIT_URL=$(shell git config --get remote.origin.url)
GIT_TAG=$(shell git describe --tags $(shell git rev-list --tags --max-count=1))

.PHONY: build clean test

build:
	GOARCH=$(GOARCH) GOOS=$(GOOS) go build -o . -ldflags " \
	-X main.vcsCommit=${GIT_COMMIT} \
	-X main.vcsURL=${GIT_URL}" ./...

install:
	GOARCH=$(GOARCH) GOOS=$(GOOS) go install -ldflags " \
	-X main.vcsCommit=${GIT_COMMIT} \
	-X main.vcsURL=${GIT_URL}" ./...

clean:
	go clean -r
	rm -f ./git-get

test:
	go test ./... -cover -race
	govulncheck ./...
