<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version"/>
  <img src="https://img.shields.io/badge/Protocol-RESP-DC382D?style=for-the-badge&logo=redis&logoColor=white" alt="RESP Protocol"/>
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License"/>
  <img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker Ready"/>
</p>

# GoRedis

A high-performance, Redis-compatible in-memory database server written from scratch in **Go**. GoRedis implements the [RESP (Redis Serialization Protocol)](https://redis.io/docs/reference/protocol-spec/) and can serve as a drop-in replacement for basic Redis use-cases — complete with **AOF persistence**, a **React-based Web Dashboard**, a **custom CLI**, and a **native Go SDK**.

> **Why GoRedis?** Built for learning, experimentation, and lightweight production workloads where a full Redis deployment is overkill.

---

## ✨ Features

| Feature | Description |
|---|---|
| 🔌 **RESP Protocol** | Full RESP parser & writer — compatible with any Redis client |
| 💾 **AOF Persistence** | Append-Only File with configurable fsync (`always`, `everysec`, `no`) |
| 🗄️ **Data Structures** | Strings, Hashes, Lists, Sets — all with O(1) core operations |
| 🌐 **Web Dashboard** | Built-in React UI to monitor and interact with your server |
| ⌨️ **Custom CLI** | Lightweight Go CLI — `go install` and start querying |
| 📦 **Go SDK** | Native client library to integrate GoRedis into your applications |
| 🐳 **Docker-Ready** | One-command deployment with Docker Compose, Nginx & auto-SSL |
| 🔒 **Thread-Safe** | `sync.RWMutex`-protected store for safe concurrent access |

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     TCP Clients                         │
│              (CLI / SDK / any Redis client)              │
└──────────────────────┬──────────────────────────────────┘
                       │ TCP :6379
                       ▼
┌──────────────────────────────────────────────────────────┐
│                   RESP Parser                            │
│           resp/parser.go · resp/writer.go                │
└──────────────────────┬───────────────────────────────────┘
                       │ resp.Value
                       ▼
┌──────────────────────────────────────────────────────────┐
│                Command Router                            │
│    cmd/handler.go → string.go | hash.go | list.go | …   │
├──────────────────────┬───────────────────────────────────┤
│                      │                                   │
│    ┌─────────────────▼─────────────────┐                 │
│    │        In-Memory Store            │                 │
│    │   store/store.go (map + RWMutex)  │                 │
│    └─────────────────┬─────────────────┘                 │
│                      │ mutating cmds                     │
│    ┌─────────────────▼─────────────────┐                 │
│    │      AOF Persistence Layer        │                 │
│    │   persistence/aof.go + replay.go  │                 │
│    └───────────────────────────────────┘                 │
└──────────────────────────────────────────────────────────┘
         │
         │ :8080
         ▼
┌──────────────────────┐
│   React Web Dashboard│
│     web/server.go    │
└──────────────────────┘
```

---

## ⚡ Performance Benchmarks

Internal benchmarks run directly against the in-memory store layer, measuring raw data-structure performance without network overhead.

**Environment:**
- **OS:** Windows 11 (amd64)
- **CPU:** 11th Gen Intel® Core™ i3-1115G4 @ 3.00 GHz (4 threads)
- **Go:** 1.26.1

```
goos: windows
goarch: amd64
pkg: github.com/SouvickSarkar20/goredis/store
cpu: 11th Gen Intel(R) Core(TM) i3-1115G4 @ 3.00GHz

BenchmarkStoreSet-4       1,990,310       743.8 ns/op      293 B/op     2 allocs/op
BenchmarkStoreGet-4       9,075,974       155.6 ns/op       13 B/op     1 allocs/op
```

| Operation | Throughput | Latency |
|---|---|---|
| **SET** | ~2.0M ops/sec | ~744 ns |
| **GET** | ~9.1M ops/sec | ~156 ns |

### Run Benchmarks Yourself

```bash
# From the project root
cd store
go test -run=^$ -bench=. -benchmem -count=1 .
```

---

## 🚀 Quick Start

### Using Docker (Recommended)

The fastest way to get GoRedis running in production. This starts the Go server, Web Dashboard, and Nginx reverse proxy with automated SSL.

```bash
docker compose up -d --build
```

| Service | Endpoint |
|---|---|
| Redis Server | `localhost:6379` |
| Web Dashboard | `http://localhost:8080` |

### From Source

```bash
git clone https://github.com/SouvickSarkar20/goredis.git
cd goredis
go run cmd/goredis-server/main.go
```

---

## ⌨️ GoRedis CLI

A lightweight command-line tool to interact with your GoRedis server.

### Install

```bash
go install github.com/SouvickSarkar20/goredis/cli@latest
```

### Usage

```bash
# Connect to a remote server
cli goredis.me:6379

# Connect to localhost
cli localhost:6379
```

### Run from Source

```bash
go run cli/main.go localhost:6379
```

---

## 📦 GoRedis SDK

Integrate GoRedis into your Go applications with the native client library.

### Install

```bash
go get github.com/SouvickSarkar20/goredis
```

### Example

```go
package main

import (
    "fmt"
    goredis "github.com/SouvickSarkar20/goredis"
)

func main() {
    client, err := goredis.NewClient("localhost:6379")
    if err != nil {
        panic(err)
    }
    defer client.Close()

    // Strings
    client.Set("username", "Souvick")
    val, _ := client.Get("username")
    fmt.Println(val) // "Souvick"

    // Hashes
    client.HSet("user:1", "name", "Souvick")
    name, _ := client.HGet("user:1", "name")
    fmt.Println(name) // "Souvick"
}
```

---

## 📋 Supported Commands

### Strings
| Command | Syntax | Description |
|---|---|---|
| `SET` | `SET key value` | Set a key to a string value |
| `GET` | `GET key` | Retrieve the value of a key |
| `DEL` | `DEL key` | Delete a key |

### Hashes
| Command | Syntax | Description |
|---|---|---|
| `HSET` | `HSET key field value` | Set a field in a hash |
| `HGET` | `HGET key field` | Get a field's value from a hash |
| `HDEL` | `HDEL key field` | Delete a field from a hash |

### Lists
| Command | Syntax | Description |
|---|---|---|
| `LPUSH` | `LPUSH key value` | Push a value to the head of a list |
| `LPOP` | `LPOP key` | Remove and return the head element |

### Sets
| Command | Syntax | Description |
|---|---|---|
| `SADD` | `SADD key member` | Add a member to a set |
| `SREM` | `SREM key member` | Remove a member from a set |
| `SISMEMBER` | `SISMEMBER key member` | Check if a member exists in a set |
| `SMEMBERS` | `SMEMBERS key` | Get all members of a set |

### Utility
| Command | Syntax | Description |
|---|---|---|
| `PING` | `PING` | Test connectivity — returns `PONG` |

---

## 🛠️ Tech Stack

| Layer | Technology |
|---|---|
| **Server** | Go (stdlib `net`, `sync`, `os`) |
| **Protocol** | RESP (Redis Serialization Protocol) |
| **Persistence** | AOF (Append-Only File) with configurable fsync |
| **Frontend** | React (Vite) |
| **Deployment** | Docker, Nginx, AWS ECS (EC2 launch type) |
| **SSL** | Certbot (auto-renewal) |

---

## 📁 Project Structure

```
goredis/
├── cmd/
│   ├── goredis-server/    # Server entry point
│   ├── handler.go         # Command router
│   ├── string.go          # String command handlers
│   ├── hash.go            # Hash command handlers
│   ├── list.go            # List command handlers
│   ├── set.go             # Set command handlers
│   └── AOFCommands.go     # AOF replay helpers
├── resp/
│   ├── parser.go          # RESP protocol parser
│   └── writer.go          # RESP protocol writer
├── store/
│   ├── store.go           # Core key-value store (map + RWMutex)
│   ├── hash.go            # Hash data structure
│   ├── list.go            # List data structure
│   ├── set.go             # Set data structure
│   └── store_test.go      # Benchmarks
├── persistence/
│   ├── aof.go             # AOF writer
│   └── replay.go          # AOF replay on startup
├── web/                   # React Web Dashboard
├── cli/                   # GoRedis CLI tool
├── client.go              # Go SDK client
├── nginx/                 # Nginx + SSL config
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

---

## 🤝 Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

---

<p align="center">
  Built with ❤️ in Go
</p>