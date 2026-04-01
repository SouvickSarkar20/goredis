# 🪵 GoRedis SSL Deployment: Issues Faced

This document chronicles the technical challenges we encountered while setting up a production-ready SSL/HTTPS environment for `goredis.me` on a DigitalOcean Droplet.

---

### ❌ Issue 1: The "Catch-22" Nginx Failure
**Symptoms**: After enabling SSL in `nginx.conf`, the Nginx container entered a crash-restart loop (`Restarting (1) 40s ago`).
**Cause**: Nginx is designed to fail if the SSL files don't exist. Since Certbot hadn't run yet, these were missing, preventing Nginx from starting. But Certbot needed Nginx to be running to solve the challenge.
**Resolution**: We used a **Bootstrap Strategy**:
1. Optimized a temporary `nginx.conf` that *only* listened on Port 80.
2. Successfully ran Certbot to generate the certificates.
3. Restored the full SSL configuration and restarted Nginx.

### ❌ Issue 2: 404 Not Found (ACME Challenge)
**Symptoms**: Certbot reported a `404: Invalid response` when trying to verify the domain.
**Cause**: Nginx was mapping the path `/.well-known/acme-challenge/` incorrectly or used an old config.
**Resolution**: We adjusted the Nginx `location` block to use an `alias` and enforced a `--force-recreate` of the Nginx container to clear the cache.

### ❌ Issue 3: Connection Timeout (Firewall)
**Symptoms**: Certbot reported a `Timeout during connect`.
**Cause**: The host's `ufw` or DigitalOcean Cloud Firewall was blocking inbound Port 80 traffic.
**Resolution**: 
1. Ran `ufw allow 80/tcp` and `ufw allow 443/tcp`.
2. Verified the DigitalOcean dashboard firewall permitted HTTP and HTTPS traffic.

### ❌ Issue 4: Entrypoint Conflict
**Symptoms**: `docker compose run certbot` was not triggering the certificate request.
**Cause**: The `docker-compose.yml` had a custom `entrypoint` for automatic renewals that blocked manual commands.
**Resolution**: Overrode the entrypoint with `--entrypoint certbot` in the command line for the initial request.

---

**Summary**: By isolating the networking, filesystem, and configuration layers, we bypassed the bootstrap "Catch-22" and secured the domain with automated renewals.
