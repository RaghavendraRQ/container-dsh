# Container-DSH

A comprehensive Docker container monitoring dashboard that provides real-time container statistics through multiple interfaces, including HTTP API, CLI, GRPC, and Prometheus metrics.

## Features

- **Real-time Container Monitoring**: Track CPU usage, memory usage, network I/O, and disk I/O
- **Multiple Interface Options**:
  - **HTTP Server**: REST API with WebSocket support for real-time updates
  - **CLI Interface**: Interactive terminal-based dashboard using Bubble Tea
  - **Prometheus Integration**: Export metrics for Prometheus monitoring
  - **Protocol Buffers**: gRPC support for efficient data serialization
- **Caching**: Redis integration for improved performance
- **Cross-platform**: Built with Go, supports multiple operating systems
- **gRPC API**: Protocol buffer-based service for efficient communication
- **Docker Integration**: Direct Docker API integration for container management
- **Development Tools**: Mock mode for testing and development

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Configuration](#configuration)
- [API Reference](#api-reference)
- [Development](#development)

## 🛠 Installation

### Prerequisites

- Go 1.24 or later
- Docker (for container monitoring)
- Redis (for caching, optional for basic functionality)
- Prometheus (optional, for metrics collection)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/raghavendrarq/container-dsh
cd container-dsh

# Build the application
go build -o bin/container-dsh ./cmd/container-dsh

# Or use the Makefile
make help
```

### Using Docker

```bash
# Build the Docker image
docker build -t container-dsh .

# Run the container
docker run -p 8080:8080 -p 3000:3000 container-dsh
```

## Quick Start

### Start the HTTP Server

```bash
# Run the server mode (default)
./bin/container-dsh --mode=server

# Or run directly with Go
go run ./cmd/container-dsh --mode=server
```

The server will start on port 8080 by default.

### CLI Interface

```bash
# Run the interactive CLI
./bin/container-dsh --mode=cli
```

### Prometheus Metrics

```bash
# Start Prometheus mode
./bin/container-dsh --mode=prometheus
```

Metrics will be available at `http://localhost:8080/metrics`

## Usage

### Available Modes

Container-DSH supports multiple operation modes:

#### 1. Server Mode (default)
Starts an HTTP server with REST API and WebSocket support.

```bash
./container-dsh --mode=server
```

**Endpoints:**
- `GET /` - List all containers
- `GET /metrics` - Get metrics for all containers
- `GET /metrics/{id}` - Get metrics for a specific container
- `WS /ws/` - WebSocket endpoint for real-time updates
- `WS /ws/container` - Container-specific WebSocket updates

#### 2. CLI Mode
Interactive terminal interface for monitoring containers.

```bash
./container-dsh --mode=cli
```

**Features:**
- Navigate with arrow keys or j/k
- Select/deselect options with Enter or Space
- Quit with q or Ctrl+C

#### 3. Prometheus Mode
Export container metrics in Prometheus format.

```bash
./container-dsh --mode=prometheus
```

**Metrics Exported:**
- `container_dsh_cpu_usage` - CPU usage percentage
- `container_dsh_mem_usage` - Memory usage in bytes
- `container_dsh_net_io` - Network I/O statistics
- `container_dsh_disk_io` - Disk I/O statistics

#### 4. Mock Mode
For testing and development purposes.

```bash
./container-dsh --mode=mock
```

#### 5. Protocol Buffer Mode
gRPC server mode.

```bash
./container-dsh --mode=pb
```

#### 6. Logger Mode
*Work in Progress*

```bash
./container-dsh --mode=logger
```

### Using the Makefile

**Note**: The current Makefile has some build configuration issues. Manual building is recommended.

```bash
# Show all available targets
make help

# Generate protocol buffers
make proto

# Clean build artifacts
make clean
```

### Using the Start Script

```bash
# Run in server mode
./start.sh server

# Run in CLI mode
./start.sh cli

# Run with Prometheus (starts both Prometheus and the app)
./start.sh prometheus
```

## Configuration

### Environment Variables

Container-DSH uses the following environment variables:

- `PORT`: Server port (default: `:8080`)
- `CLIENT_URL`: Frontend client URL for CORS (required for server mode)

### Example Configuration

```bash
export PORT=":8080"
export CLIENT_URL="http://localhost:3000"
```

### Redis Configuration

Redis is used for caching container information to improve performance. The application will attempt to connect to Redis on the default port (6379), but will fall back to direct Docker API calls if Redis is unavailable.

**Note**: If Redis is not available, you may see warning messages, but the application will continue to function.

### Prometheus Configuration

If using Prometheus integration, configure `prometheus.yml`:

```yaml
global:
  scrape_interval: 1s
  evaluation_interval: 3s

scrape_configs:
  - job_name: "container-dsh"
    static_configs:
      - targets: ["localhost:8080"]
```

## API Reference

### REST Endpoints

#### GET /
Returns a list of all containers.

**Response:**
```json
[
  "container_id_1",
  "container_id_2"
]
```

#### GET /metrics
Returns metrics for all containers.

**Response:**
```json
[
  "container_id_1",
  "container_id_2"
]
```

#### GET /metrics/{id}
Returns detailed metrics for a specific container.

**Response:**
```json
{
  "id": "container_id",
  "name": "container_name",
  "status": 2,
  "cpu_usage": 25.5,
  "mem_usage": 1073741824,
  "net_io": 1024,
  "disk_io": 2048
}
```

### WebSocket Endpoints

#### WS /ws/
General WebSocket endpoint for real-time updates.

#### WS /ws/container
Container-specific WebSocket endpoint.

### gRPC API

Container-DSH also provides a gRPC service defined in `api/proto/v1/container.proto`.

**Service Definition:**
```protobuf
service ContainerService {
    rpc GetContainerStats(ContainerStatRequest) returns (Container);
}
```

**Run gRPC Server:**
```bash
./bin/container-dsh --mode=pb
```

### Container Status Codes

- `1`: Created
- `2`: Running
- `3`: Paused
- `4`: Restarting
- `5`: Removing
- `6`: Exited
- `7`: Dead

## Development

### Prerequisites for Development

- Go 1.24+
- Docker (running and accessible)
- Protocol Buffers Compiler (`protoc`) - for gRPC development
- Redis (optional, for caching features)

### Project Structure

```
container-dsh/
├── api/                    # API definitions
│   ├── gen/go/            # Generated Go code
│   └── proto/v1/          # Protocol buffer definitions
├── cmd/                   # Application entry points
│   ├── cli/               # CLI mode
│   ├── container-dsh/     # Main application
│   ├── mock/              # Mock mode
│   ├── pb/                # Protocol buffer mode
│   ├── prom/              # Prometheus mode
│   └── server/            # Server mode
├── internal/              # Internal packages
│   ├── cache/             # Redis caching
│   ├── cli/               # CLI models and logic
│   ├── config/            # Configuration management
│   └── container/         # Container operations
├── pkg/                   # Public packages
│   ├── aggr/              # Aggregator
│   ├── collector/         # Data collection
│   ├── http/              # HTTP server
│   ├── logger/            # Logging utilities
│   ├── pb/                # Protocol buffer services
│   ├── prom/              # Prometheus integration
│   └── snapshot/          # Snapshot utilities
├── Dockerfile             # Docker configuration
├── Makefile              # Build configuration
├── go.mod                # Go module definition
└── README.md             # This file
```

### Adding New Features

1. **Container Metrics**: Add new metrics in `internal/container/model.go`
2. **HTTP Endpoints**: Add routes in `pkg/http/server.go` and handlers in `pkg/http/handler.go`
3. **CLI Options**: Extend the CLI model in `internal/cli/model.go`
4. **Prometheus Metrics**: Add metrics in `pkg/prom/model.go`

### Building and Testing

```bash
# Install dependencies
go mod download

# Build the application manually (recommended)
go build -o bin/container-dsh ./cmd/container-dsh

# Test the build
./bin/container-dsh --help

# Generate protocol buffers (requires protoc)
make proto

# Build Docker image
docker build -t container-dsh .
```

## Troubleshooting

### Common Issues

**"redis connection failed"**: 
- Redis is not running or not accessible
- Solution: Install and start Redis, or ignore the warning (app will work without caching)

**"error in getting containers"**:
- Docker is not running or not accessible
- Solution: Start Docker daemon and ensure proper permissions

**Build errors with Makefile**:
- Current Makefile has configuration issues
- Solution: Use manual build command: `go build -o bin/container-dsh ./cmd/container-dsh`

**Permission denied accessing Docker**:
- User doesn't have Docker permissions
- Solution: Add user to docker group or run with sudo


## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - For the beautiful CLI interface
- [Docker](https://www.docker.com/) - For containerization support
- [Prometheus](https://prometheus.io/) - For metrics collection
- [Go](https://golang.org/) - For the robust foundation

## Support

If you encounter any problems or have questions, please:

1. Check the [Issues](https://github.com/raghavendrarq/container-dsh/issues) page
2. Create a new issue if your problem isn't already reported
3. Provide as much detail as possible including:
   - Operating system
   - Go version
   - Docker version
   - Steps to reproduce the issue

---

Made with ❤️ by [RaghavendraRQ](https://github.com/RaghavendraRQ)
