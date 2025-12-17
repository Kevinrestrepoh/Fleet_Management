# Fleet Management Platform (Real-Time Telemetry & Control)

A real-time fleet management system that ingests telemetry from thousands of simulated vehicles, allows live control commands, and streams aggregated fleet metrics to clients.

Built to explore distributed systems, gRPC streaming, and async backends using **Rust** and **Go**.

## Features

### Bidirectional gRPC Streaming
- **Vehicles stream telemetry** to the ingestion server.
- **Server streams commands** back to individual vehicles.

### Live Command Routing
- Target a specific vehicle with the following commands:
  - `PING`
  - `UPDATE_RATE`
  - `SHUTDOWN`

### Vehicle Registry
- Tracks connected vehicles using async-safe structures.

### Fleet-Wide Metrics
- **Active vehicles**
- **Low battery vehicles**
- **Average speed and engine temperature**

### Metrics Streaming via SSE
- HTTP clients receive live fleet updates via Server-Sent Events (SSE).

### Multi-Language System
- **Rust**: Ingestion and metrics server
- **Go**: Vehicle simulator and control API

---

## Example Metrics

- Active vehicles
- Low battery vehicles
- Average speed
- Average engine temperature

---

## Getting Started

Follow these steps to get the system running locally:

## 1. Install golang grpc packages

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 2. Start Ingestion Server
To run the ingestion server (written in **Rust**), execute the following command:

```bash
cargo run
```

### 3. Start Control API
To run the control API (written in Go), execute the following command:

```bash
go run main.go
```

### 3. Start Vehicle Simulator
To start the vehicle simulator (also written in Go), execute the following command:


```bash
go run main.go
```


