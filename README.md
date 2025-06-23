[![Build Status](https://github.com/kevholditch/vigilant/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/kevholditch/vigilant/actions/workflows/go.yml)

# Vigilant

Vibe coding the missing Kubernetes UI!

## About

Vigilant is an experimental terminal-based Kubernetes management tool built with Go and Bubble Tea. This project serves as a learning exercise in building intuitive terminal UIs and exploring AI-assisted development workflows.

## Features

- Terminal-based user interface using Bubble Tea
- Table-based pod display with Kubernetes integration
- Keyboard navigation
- Pod description view with kubectl integration
- Simple and clean interface with cyberpunk theme
- Real-time cluster information display

## Installation

### Prerequisites

- Go 1.24 or later
- kubectl configured with access to a Kubernetes cluster

### Build and run

```bash
git clone https://github.com/kevholditch/vigilant.git
cd vigilant
go mod tidy
make run
```

### Build binary

```bash
make build
./vigilant
```

## Usage

```bash
# Run the application
make run
# or
./vigilant
```

### Controls

#### Pod List View
- `q` - Quit the application
- `↑/↓` or `j/k` - Navigate through pods
- `d` - Describe selected pod (opens pod description view)

#### Pod Description View
- `Esc` - Return to pod list view
- `↑/↓` or `j/k` - Scroll through pod description

## Development

This is a hobby project exploring terminal UI development with the following features:

1. **Pod List View**: Displays all pods in the cluster with status, namespace, and other details
2. **Pod Description View**: Shows detailed pod information using `kubectl describe pod`
3. **Theme System**: Cyberpunk-themed UI with consistent styling
4. **Keyboard Navigation**: Full keyboard support for navigation and interaction
5. **Cluster Information**: Real-time display of cluster name, Kubernetes version, and node counts

### Available Make Commands

```bash
make build      # Build the binary
make run        # Run the application
make test       # Run all tests (includes setting up envtest binaries)
make clean      # Clean build artifacts
```

### Testing

The project includes comprehensive tests using the Kubernetes controller-runtime envtest framework:

```bash
make test       # Run all tests with proper envtest setup
```

The test suite includes BDD-style scenarios for the header controller with given-when-then structure.

### Future Ideas

You can extend it by:

1. Adding more Kubernetes resources (services, deployments, etc.)
2. Implementing resource management operations
3. Adding configuration options
4. Implementing real-time updates
5. Adding resource filtering and search capabilities
6. Implementing multi-cluster support

## Contributing

This is a personal hobby project, but feel free to fork and experiment with your own ideas!

## License

This project is licensed under the MIT License - see the LICENSE file for details. 