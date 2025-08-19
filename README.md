# Nornir Workflow System

```ascii
    ╭───╮  ╭───╮  ╭───╮
    │ o │  │ o │  │ o │
    │\│/│  │\│/│  │\│/│
    │ │ │  │ │ │  │ │ │
    │/ \│  │/ \│  │/ \│
    ╰─⟡─╯  ╰─⟡─╯  ╰─⟡─╯
    Urðr   Verðandi Skuld
    Past   Present Future
```

## Overview

A distributed workflow orchestration system built with Go, connecting services through a REST/gRPC architecture.

## Build and Deploy

The project uses a Makefile for all build and deployment operations:

### Protocol Buffers

```bash
# Generate protobuf code
make proto

# Clean generated files
make proto-clean
```

### Testing

```bash
# Run all tests
make test
```

### Docker Build

```bash
# Build both services
make build
```

Images built:
- `nornir-gateway:latest`
- `nornir-worker:latest`

### Deployment

```bash
# Deploy to Kubernetes with custom values
make helm-install

# Build and deploy locally
make helm-install-local

# Quick deploy with local images
make deploy
```

## Project Structure

```
nornir/
├── proto/          # Protocol buffer definitions
├── gateway-service # REST API gateway
├── worker-service  # gRPC workflow processor
└── charts/         # Helm charts for Kubernetes deployment
```

## API Usage

**Start a Workflow:**
```http
POST /workflows
Content-Type: application/json

{
    "name": "my-workflow"
}
```

**Response:**
```json
{
    "id": "<uuid>",
    "status": "STARTED"
}
```

## Development Requirements

- Go 1.21+
- Docker
- Protocol Buffers compiler
- Kubernetes
- Helm

## Configuration

The system supports configuration through Helm values:
- Gateway service image and tag
- Worker service image and tag
- Deployment parameters

## License

MIT License