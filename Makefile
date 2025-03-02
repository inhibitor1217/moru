MODULE_NAME := $(shell go list -m)

# Consider the location of Makefile as the root project directory
PROJECT_ROOT := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
PROJECT_NAME := $(shell basename $(PROJECT_ROOT))

# Artifacts
TARGET_DIR := $(PROJECT_ROOT)/build

# Application
TARGET ?= morusrv
STAGE ?= dev
VERSION := $(shell git describe --exact-match --tags HEAD 2>/dev/null || git rev-parse --abbrev-ref HEAD)

# Go
GO := go
GOVERSION := $(shell go version | awk '{print $$3}')
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
LDFLAGS := -ldflags=""

# Android (for cross-compilation)
ANDROID_SDK ?= $(HOME)/Library/Android/sdk
NDK_BIN ?= $(ANDROID_SDK)/ndk/28.0.13004108/toolchains/llvm/prebuilt/darwin-x86_64/bin

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

build-so-android: build-so-android-x86_64 build-so-android-arm64

build-so-android-x86_64: init
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=amd64 \
	CC=${NDK_BIN}/x86_64-linux-android35-clang \
	${GO} build ${LDFLAGS} \
		-buildmode=c-shared \
		-o ${TARGET_DIR}/jniLibs/x86_64/lib${TARGET}.so \
		${PROJECT_ROOT}/cmd/${TARGET}

build-so-android-arm64: init
	CGO_ENABLED=1 \
	GOOS=android \
	GOARCH=arm64 \
	CC=${NDK_BIN}/aarch64-linux-android35-clang \
	${GO} build ${LDFLAGS} \
		-buildmode=c-shared \
		-o ${TARGET_DIR}/jniLibs/arm64-v8a/lib${TARGET}.so \
		${PROJECT_ROOT}/cmd/${TARGET}

# ------------------------------------------------------------------------------
# Development Tools
# ------------------------------------------------------------------------------

run:
	${TARGET_DIR}/${TARGET}.${GOOS}.${GOARCH}

dev: build run

# ------------------------------------------------------------------------------
# Utilities
# ------------------------------------------------------------------------------

generate: generate-go generate-proto

generate-go:
	${GO} generate ./...

generate-proto:
	protoc -I ./moru-proto \
		--go_out=proto \
		--go_opt=module=github.com/inhibitor1217/moru/proto \
		./moru-proto/**/*.proto

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
