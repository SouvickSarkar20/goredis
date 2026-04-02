# GoRedis

GoRedis is a high-performance **Redis-compatible** server built from the ground up in Go. It supports the RESP (Redis Serialization Protocol) and can be used as a drop-in replacement for basic Redis functionality in your applications.

---

## Architecture

- **Go Server**: Handles TCP connections, RESP parsing, and command execution.
- **Data Store**: Custom in-memory storage with O(1) time complexity for core operations.
- **Persistence**: **Append-Only File (AOF)** persistence to ensure your data survives restarts.
- **Web UI**: Built-in React dashboard to monitor and interact with your server from any browser.
- **Custom CLI**: A lightweight Go command-line tool to manage your server from any terminal.

---

## Quick Start: Running the Server

### Using Docker (Production & Development)
The easiest way to run GoRedis is using **Docker Compose**. This will start the Go server, the Web UI, and an Nginx reverse proxy with automated SSL.

```bash
docker compose up -d --build
```
*   **Redis Server**: `localhost:6379`
*   **Web Dashboard**: `http://localhost:8080` (or `https://goredis.me` in production)

---

## ⚡ Using the GoRedis Ecosystem

You can interact with your GoRedis server directly from your terminal or integrate it into your apps using our official tools.

### GoRedis CLI
The fastest way to get the `goredis-cli` command on your machine without cloning the repo:

**1. Install globally (Recommended):**
```bash
go install github.com/krishsinghhura/go-redis/cli@latest
```

**2. Run it!**
```bash
cli goredis.me:6379 
```

**3. Or run from source:**
```bash
# If the repo is cloned locally:
go run cli/main.go goredis.me:6379
```

---

### GoRedis SDK
To integrate GoRedis into your own Go applications:

**1. Install the SDK:**
```bash
go get github.com/krishsinghhura/go-redis
```

**2. Use it in your code:**
```go
import "github.com/krishsinghhura/go-redis"

func main() {
    client, _ := goredis.NewClient("goredis.me:6379")
    defer client.Close()

    client.Set("username", "Krish")
    val, _ := client.Get("username")
}
```

---

## Supported Commands

GoRedis supports the most common Redis commands, including:
- **Strings**: `SET`, `GET`, `DEL`, `INCR`, `EXISTS`
- **Hashes**: `HSET`, `HGET`, `HDEL`, `HGETALL`
- **Lists**: `LPUSH`, `RPUSH`, `LPOP`, `RPOP`, `LLEN`
- **Sets**: `SADD`, `SREM`, `SCARD`, `SISMEMBER`

---

## Tech Stack
- **Server**: Go
- **Frontend**: React
- **Deployment**: Docker, Nginx, ECS(EC2 launch type)
- **Protocol**: RESP (Redis Serialization Protocol)