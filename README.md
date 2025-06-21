# Vigilant

A K9s-like terminal UI application built with Go.

## Features

- Terminal-based user interface using Bubble Tea
- Table-based pod display with Kubernetes integration
- Keyboard navigation
- Pod description view with kubectl integration
- Simple and clean interface with cyberpunk theme

## Installation

### Prerequisites

- Go 1.24 or later
- kubectl configured with access to a Kubernetes cluster

### Build and run

```bash
git clone https://github.com/kevholditch/vigilant.git
cd vigilant
go mod tidy
go run .
```

### Build binary

```bash
go build -o vigilant .
./vigilant
```

## Usage

```bash
# Run the application
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

## Project Structure

```
vigilant/
├── internal/              # Private application code
│   ├── app/              # Main application logic
│   ├── models/           # Data models
│   ├── theme/            # UI theme and styling
│   └── views/            # UI view components
├── main.go               # Application entry point
├── go.mod                # Go module file
└── README.md             # This file
```

## Development

This is a K9s-like terminal UI application with the following features:

1. **Pod List View**: Displays all pods in the cluster with status, namespace, and other details
2. **Pod Description View**: Shows detailed pod information using `kubectl describe pod`
3. **Theme System**: Cyberpunk-themed UI with consistent styling
4. **Keyboard Navigation**: Full keyboard support for navigation and interaction

You can extend it by:

1. Adding more Kubernetes resources (services, deployments, etc.)
2. Implementing resource management operations
3. Adding configuration options
4. Implementing real-time updates

## License

This project is licensed under the MIT License - see the LICENSE file for details. 