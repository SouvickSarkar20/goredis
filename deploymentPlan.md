# 🏗️ GoRedis Deployment Plan (Droplet Edition)

This document explains our professional deployment strategy for GoRedis on a DigitalOcean Droplet.

---

## 🧭 Overall Strategy
We use the **Droplet + Docker Compose** model, which is the equivalent of "ECS on EC2" on AWS. This gives us full control over the persistence of our Redis records and the lifecycle of our SSL certificates.

### Why not App Platform (Serverless)?
We chose the Droplet approach because:
1.  **Persistence**: We need a permanent disk to store our AOF records (/data). 
2.  **TCP Support**: We need Port 6379 (Redis) open for non-HTTP clients. 
3.  **Certbot Control**: We have full ownership of our SSL certificate logic.

---

## 🐳 Our Containers (The "Services")

### 1. `goredis` (The Core)
*   **Role**: All-in-one Go server and React frontend.
*   **Build**: Custom multi-stage `Dockerfile` (Node for React, Go for Server, Alpine for Final).
*   **Networking**: Exposed on **Port 6379** (TCP) for direct Redis traffic.
*   **Persistence**: Mounted to `./data` on the host disk to keep AOF records safe.

### 2. `nginx` (The Ingress)
*   **Role**: Reverse proxy and SSL termination.
*   **Networking**: Exposed on Port **80** (HTTP) and Port **443** (HTTPS). 
*   **Logic**:
    *   Redirects all Port 80 traffic to Port 443.
    *   Solves the Let's Encrypt challenge via `/.well-known/acme-challenge/`.
    *   Proxies all secure traffic to the internal `goredis:8080` port.

### 3. `certbot` (The Guard)
*   **Role**: Automated SSL certificate manager.
*   **Lifecyle**: 
    *   On startup, it runs in a 12-hour loop.
    *   It only requests/renews certificates if they are within 30 days of expiration.
    *   Uses the shared `nginx/www` folder to communicate with Nginx during the ACME challenge.

---

## 🚀 How we Deployed

### Step A: Provisioning
1.  Created a DigitalOcean Droplet with the **"Docker on Ubuntu"** marketplace image.
2.  Used Cloud Firewall to open Ports 22 (SSH), 80 (HTTP), 443 (HTTPS), and 6379 (Redis).

### Step B: The Bootstrap
Because Nginx fails if SSL certs are missing, we used a three-step bootstrap:
1.  Started a Port-80-only Nginx container.
2.  Ran a manual Certbot container to fetch the first certificate from Let's Encrypt.
3.  Restarted Nginx with the full secure config.

### Step C: Verification
1.  `docker compose ps` to verify all services are "Up".
2.  Navigate to `https://goredis.me` (SSL confirmed).
3.  Run `redis-cli` against the Droplet IP (Redis TCP confirmed).

---

**Summary**: This architecture is robust, portable, and allows for automated data backups (by simply copying the `./data` folder).
