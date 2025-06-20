# Makefile for Vigilant

# Variables
BINARY_NAME=vigilant

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run

.PHONY: all build run clean

all: build

# Build the application binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) .

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	$(GORUN) .

# Clean the build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME) 