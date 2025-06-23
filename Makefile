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
	@echo "Running tests with KUBEBUILDER_ASSETS..."
	@KUBEBUILDER_ASSETS="$(shell $(LOCALBIN)/setup-envtest use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" go test ./... -v

# Debug target to help troubleshoot CI issues
debug-ci:
	@echo "Debugging CI environment..."
	@echo "Current directory: $(shell pwd)"
	@echo "LOCALBIN: $(LOCALBIN)"
	@echo "ENVTEST_K8S_VERSION: $(ENVTEST_K8S_VERSION)"
	@echo "ENVTEST_VERSION: $(ENVTEST_VERSION)"
	@ls -la $(LOCALBIN) || echo "LOCALBIN directory does not exist"
	@echo "Installing setup-envtest..."
	@GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@$(ENVTEST_VERSION)
	@echo "setup-envtest installed at: $(LOCALBIN)/setup-envtest"
	@ls -la $(LOCALBIN)/setup-envtest
	@echo "Downloading Kubernetes binaries..."
	@$(LOCALBIN)/setup-envtest use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path
	@echo "KUBEBUILDER_ASSETS path: $(shell $(LOCALBIN)/setup-envtest use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)"
	@echo "Checking binary directory:"
	@ls -la $(LOCALBIN)/k8s/ || echo "k8s directory does not exist"
	@find $(LOCALBIN) -name "etcd" -type f 2>/dev/null || echo "etcd binary not found"

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