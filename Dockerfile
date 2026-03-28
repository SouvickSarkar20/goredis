# ── Stage 1: Build the React frontend ───────────────────────────────────────
# We use the official Node image to run npm install and npm run build.
# After this stage we only care about the web/dist/ folder it produces.
FROM node:18-alpine AS frontend-builder

WORKDIR /app/client

# Copy package files first — Docker caches this layer.
# If package.json hasn't changed, npm install is skipped on next build.
COPY client/package.json client/package-lock.json ./

RUN npm install

# Now copy the rest of the client source
COPY client/ .

# Build the React app — outputs to ../web/dist (as configured in vite.config.js)
RUN npm run build


# ── Stage 2: Build the Go binary ────────────────────────────────────────────
# We use the official Go image to compile the server.
# After this stage we only care about the compiled binary.
FROM golang:1.22-alpine AS backend-builder

WORKDIR /app

# Copy go.mod first — cached layer, only re-downloads if dependencies change
COPY go.mod ./
RUN go mod download

# Copy all Go source code
COPY . .

# Copy the built React app from Stage 1 into web/dist/
# The Go binary needs to find this folder at runtime
COPY --from=frontend-builder /app/client/../web/dist ./web/dist

# Build a static binary — CGO_ENABLED=0 means pure Go, no C libraries needed
# This is important for Alpine which doesn't have glibc by default
# -o goredis names the output file
RUN CGO_ENABLED=0 GOOS=linux go build -o goredis .


# ── Stage 3: The final minimal image ────────────────────────────────────────
# Alpine is a 5MB Linux distro — just enough to run a binary.
# We copy ONLY what's needed to run: the binary and the web/dist folder.
# Everything else (Go compiler, Node, npm, source code) is left behind.
FROM alpine:latest

WORKDIR /app

# Copy the compiled Go binary from Stage 2
COPY --from=backend-builder /app/goredis .

# Copy the built frontend from Stage 2
COPY --from=backend-builder /app/web/dist ./web/dist

# Tell Docker this container listens on these ports.
# This is documentation — the actual port mapping happens in docker-compose.yml
EXPOSE 6379
EXPOSE 8080

# Run the binary when the container starts
CMD ["./goredis"]