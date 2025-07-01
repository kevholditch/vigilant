# Makefile for Vigilant

# Variables
BINARY_NAME=vigilant
MAIN_FILE=main.go

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run

# Detect if we're in CI (GitHub Actions sets CI=true)
CI ?= false

# Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

# Tool Binaries
ENVTEST ?= $(LOCALBIN)/setup-envtest

# Tool Versions
#ENVTEST_VERSION is the version of controller-runtime release branch to fetch the envtest setup script
ENVTEST_VERSION ?= $(shell go list -m -f "{{ .Version }}" sigs.k8s.io/controller-runtime | awk -F'[v.]' '{printf "release-%d.%d", $$2, $$3}')
#ENVTEST_K8S_VERSION is the version of Kubernetes to use for setting up ENVTEST binaries
ENVTEST_K8S_VERSION ?= $(shell go list -m -f "{{ .Version }}" k8s.io/api | awk -F'[v.]' '{printf "1.%d", $$3}')

.PHONY: all build run clean test test-ci setup-envtest envtest

all: build

# Build the application binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_FILE)

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	$(GORUN) $(MAIN_FILE)

# Clean the build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)

# Test target for local development
test: setup-envtest
	@echo "Running tests..."
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" go test ./... -v

# Test target for CI environments
test-ci:
	@echo "Running tests in CI environment..."
	@echo "Installing envtest binaries for Kubernetes version $(ENVTEST_K8S_VERSION)..."
	@GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@$(ENVTEST_VERSION)
	@$(LOCALBIN)/setup-envtest use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path
	@echo "Setting up environment variables..."
	@export KUBEBUILDER_ASSETS="$(shell $(LOCALBIN)/setup-envtest use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" && \
	echo "KUBEBUILDER_ASSETS=$$KUBEBUILDER_ASSETS" && \
	echo "Checking for etcd binary..." && \
	ls -la $$KUBEBUILDER_ASSETS && \
	echo "Running tests..." && \
	go test ./... -v

# Alternative test target for CI environments (simpler approach)
test-ci-simple:
	@echo "Running tests in CI environment (simple approach)..."
	@echo "Installing envtest binaries for Kubernetes version $(ENVTEST_K8S_VERSION)..."
	@GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@$(ENVTEST_VERSION)
	@echo "Downloading Kubernetes binaries for current platform..."
	@$(LOCALBIN)/setup-envtest use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path
	@echo "Running tests with KUBEBUILDER_ASSETS and PATH..."
	@export KUBEBUILDER_ASSETS="$$($(LOCALBIN)/setup-envtest use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" && \
	export PATH="$$KUBEBUILDER_ASSETS:$$PATH" && \
	go test ./... -v

setup-envtest: envtest
	@echo "Setting up envtest binaries for Kubernetes version $(ENVTEST_K8S_VERSION)..."
	@$(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path || { \
		echo "Error: Failed to set up envtest binaries for version $(ENVTEST_K8S_VERSION)."; \
		exit 1; \
	}

envtest: $(ENVTEST)
$(ENVTEST): $(LOCALBIN)
	$(call go-install-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest,$(ENVTEST_VERSION))

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef

# Download and install a specific version of envtest binaries
# Usage: make download-envtest VERSION=1.29.0 PLATFORM=darwin-arm64
download-envtest:
	@if [ -z "$(VERSION)" ] || [ -z "$(PLATFORM)" ]; then \
		echo "Usage: make download-envtest VERSION=<version> PLATFORM=<platform>"; \
		echo "Example: make download-envtest VERSION=1.29.0 PLATFORM=darwin-arm64"; \
		exit 1; \
	fi
	@mkdir -p bin/k8s/$(VERSION)-$(PLATFORM)
	@echo "Downloading envtest binaries for $(VERSION) on $(PLATFORM)..."
	@curl -L https://github.com/kubernetes-sigs/kubebuilder/releases/download/v$(VERSION)/kubebuilder_$(VERSION)_$(PLATFORM).tar.gz | tar -xz -C bin/k8s/$(VERSION)-$(PLATFORM) --strip-components=1
	@echo "Creating symlink..."
	@ln -sf $(VERSION)-$(PLATFORM) bin/k8s/latest
	@echo "Envtest binaries downloaded and symlinked to bin/k8s/latest"

# Download and install the latest envtest binaries for the current platform
download-envtest-latest:
	@echo "Detecting platform..."
	@PLATFORM=$$(case "$$(uname -s)" in \
		Darwin) \
			case "$$(uname -m)" in \
				arm64) echo "darwin-arm64" ;; \
				x86_64) echo "darwin-amd64" ;; \
				*) echo "darwin-amd64" ;; \
			esac ;; \
		Linux) \
			case "$$(uname -m)" in \
				x86_64) echo "linux-amd64" ;; \
				arm64) echo "linux-arm64" ;; \
				*) echo "linux-amd64" ;; \
			esac ;; \
		*) echo "linux-amd64" ;; \
	esac); \
	VERSION=$$(curl -s https://api.github.com/repos/kubernetes-sigs/kubebuilder/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' | sed 's/v//'); \
	echo "Latest version: $$VERSION"; \
	echo "Platform: $$PLATFORM"; \
	$(MAKE) download-envtest VERSION=$$VERSION PLATFORM=$$PLATFORM

# Define a function to create symlinks
define create-symlink
	@if [ -d "bin/k8s/$(1)-$(2)" ]; then \
		echo "Creating symlink for $(1)-$(2)"; \
		ln -sf $(1)-$(2) bin/k8s/$(1); \
	else \
		echo "Directory bin/k8s/$(1)-$(2) does not exist"; \
		exit 1; \
	fi
endef

# Create symlinks for specific versions
symlink-1.33.0:
	$(call create-symlink,1.33.0,darwin-arm64) 