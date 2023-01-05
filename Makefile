#! /usr/bin/make -f
# Go related variables.
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin


# Go files.
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

# Common commands.
all: fmt test
development: precommit-install githooks-install

gofmt:
	gofmt -s -w ${GOFMT_FILES}

run:
	go run ./apis/main.go

build:
	go build ./apis/main.go

test:
	@echo "  >  Running unit tests."
	GOBIN=$(GOBIN) go test -cover -race -coverprofile=coverage.txt -covermode=atomic -v ./...