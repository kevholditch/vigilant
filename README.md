# Vigilant

A K9s-like terminal UI application built with Go.

## Features

- Terminal-based user interface using tview
- Table-based resource display
- Keyboard navigation
- Simple and clean interface

## Installation

### Prerequisites

- Go 1.24 or later

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

- `q` - Quit the application
- Arrow keys - Navigate the table
- More controls coming soon...

## Project Structure

```
vigilant/
├── internal/              # Private application code
│   └── app/              # Main application logic
├── main.go               # Application entry point
├── go.mod                # Go module file
└── README.md             # This file
```

## Development

This is a basic skeleton for a K9s-like terminal UI application. You can extend it by:

1. Adding more views and pages
2. Implementing resource management
3. Adding configuration options
4. Implementing real data sources

## License

This project is licensed under the MIT License - see the LICENSE file for details. 