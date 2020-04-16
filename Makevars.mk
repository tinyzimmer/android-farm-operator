# Go options
GO111MODULE ?= auto
CGO_ENABLED ?= 1
GOROOT ?= `go env GOROOT`

GIT_COMMIT ?= `git rev-parse HEAD`

# Golang CI Options
GOLANGCI_VERSION ?= 1.23.8
GOLANGCI_LINT ?= _bin/golangci-lint
GOLANGCI_DOWNLOAD_URL ?= https://github.com/golangci/golangci-lint/releases/download/v${GOLANGCI_VERSION}/golangci-lint-${GOLANGCI_VERSION}-$(shell uname | tr A-Z a-z)-amd64.tar.gz

# Operator SDK options
SDK_VERSION ?= v0.16.0

UNAME := $(shell uname)
ifeq ($(UNAME), Linux)
SDK_PLATFORM := linux-gnu
endif
ifeq ($(UNAME), Darwin)
SDK_PLATFORM := apple-darwin
endif

# Image Options
OPERATOR_IMAGE ?= quay.io/tinyzimmer/android-farm-operator:${APP_VERSION}
EMULATOR_IMAGE ?= quay.io/tinyzimmer/android-emulator:${EMULATOR_PLATFORM}-slim
ADB_IMAGE ?= quay.io/tinyzimmer/adbmon:latest
REDIR_IMAGE ?= quay.io/tinyzimmer/goredir:latest
EMULATOR_PLATFORM ?= android-29
EMULATOR_CLI_TOOLS ?= 6200805
ANDROID_TOOLS_VERSION ?= 29.0.6-r0

# Operator SDK
OPERATOR_SDK ?= _bin/operator-sdk
OPERATOR_SDK_URL ?= https://github.com/operator-framework/operator-sdk/releases/download/${SDK_VERSION}/operator-sdk-${SDK_VERSION}-x86_64-${SDK_PLATFORM}

# Kind Options
KIND ?= _bin/kind
KIND_VERSION ?= v0.7.0
KIND_DOWNLOAD_URL ?= https://github.com/kubernetes-sigs/kind/releases/download/${KIND_VERSION}/kind-$(shell uname)-amd64
KUBERNETES_VERSION ?= v1.18.0
METALLB_VERSION ?= v0.9.3
HELM_ARGS ?= --set operator.api.enabled=false

# Gendocs
REFDOCS ?= _bin/refdocs
REFDOCS_CLONE ?= $(dir ${REFDOCS})/gen-crd-api-reference-docs
