#!/usr/bin/make -f

# Build version

git_tag := $(shell git describe --tags --abbrev=0)
ifeq ($(git_tag),)
  git_tag := v0.0.0
endif
git_commit := $(shell git rev-list -1 HEAD)
app_version := $(git_tag)-$(git_commit)

# Process linker flags

ldflags = -X github.com/itiky/eth-block-proxy/cmd/eth-block-proxy/cmd.VersionGitTag=$(git_tag) \
		  -X github.com/itiky/eth-block-proxy/cmd/eth-block-proxy/cmd.VersionGitCommit=$(git_commit) \

ldflags := $(strip $(ldflags))

# Makefile rules

export GO111MODULE=on

all: install

install: go.sum install-proxy

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

install-proxy:
	@echo "--> Eth-block-proxy app build / install"
	@echo "  App version: $(app_version)"
	@go install -ldflags "$(ldflags)" -tags "$(build_tags)" ./cmd/eth-block-proxy

build-docker:
	@echo "--> Docker build"
	@docker build -t eth-block-proxy .

lint:
	@echo "--> Running Golang linter (unused variable / function warning are skipped)"
	@golangci-lint run --exclude 'unused'

tests:
	@echo "--> Running tests (no cached test results)"
	go test ./... -v --count=1
