MODULE_NAME := $(shell go list -m)

# Consider the location of Makefile as the root project directory
PROJECT_ROOT := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
PROJECT_NAME := $(shell basename $(PROJECT_ROOT))

# Artifacts
TARGET_DIR := $(PROJECT_ROOT)/build

# Application
TARGET ?= moru
STAGE ?= dev
VERSION := $(shell git describe --exact-match --tags HEAD 2>/dev/null || git rev-parse --abbrev-ref HEAD)

# Go
GO := go
GOVERSION := $(shell go version | awk '{print $$3}')
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
LDFLAGS := -ldflags=""

# ------------------------------------------------------------------------------
# Build Targets
# ------------------------------------------------------------------------------

env:
	@echo "PROJECT_ROOT: $(PROJECT_ROOT)"
	@echo "PROJECT_NAME: $(PROJECT_NAME)"
	@echo "TARGET_DIR: $(TARGET_DIR)"
	@echo "STAGE: $(STAGE)"
	@echo "VERSION: $(VERSION)"
	@echo "GO: $(GO)"
	@echo "GOVERSION: $(GOVERSION)"
	@echo "GOOS: $(GOOS)"
	@echo "GOARCH: $(GOARCH)"
	@echo "LDFLAGS: $(LDFLAGS)"

init:
	${GO} mod download

build: init
	GOOS=${GOOS} \
	GOARCH=${GOARCH} \
	${GO} build ${LDFLAGS} \
		-o ${TARGET_DIR}/${TARGET}.${GOOS}.${GOARCH} \
		${PROJECT_ROOT}/cmd/${TARGET}

run:
	${TARGET_DIR}/${TARGET}.${GOOS}.${GOARCH}

dev: build run

generate:
	${GO} generate ./...

clean:
	rm -rf ${TARGET_DIR}

test: build
	go clean -testcache
	go test `go list ./... | grep -v /generated/`

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...
	go mod tidy
