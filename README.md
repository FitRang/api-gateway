# traefik

This repository contains the API Gateway configuration using Traefik, along with custom Go plugins for authentication and request middleware.

## Overview

Traefik acts as the entry point for all external traffic into the FitRang backend. It is responsible for:

* Routing requests to appropriate backend services
* Handling CORS
* Authenticating users via a custom Firebase authentication plugin
* Supporting both HTTP and WebSocket connections
* Injecting verified user identity headers into upstream services

---

## Structure

```
api-gateway/
│
├── Dockerfile            # Builds Traefik with local plugins
├── traefik.yml          # Static configuration (entrypoints, providers, plugins)
├── dynamic.yml          # Routing, services, and middleware configuration
│
└── plugins/
    └── firebaseauth/
        ├── go.mod
        ├── extractToken.go
        ├── jwks.go
        ├── plugin.go
        └── validateToken.go
```

---

## Authentication Plugin

The `firebaseauth` plugin verifies Firebase ID tokens using Google's public JWKS endpoint.

### Responsibilities

* Extract token from:

  * `Authorization: Bearer <token>` (HTTP)
  * `access_token` query param (WebSocket)
* Verify JWT signature using JWKS
* Validate:

  * issuer (`iss`)
  * audience (`aud`)
  * expiration (`exp`)
  * email verification (`email_verified`)
* Inject headers into upstream request:

  * `X-User-Email`

---

## Routed Services

| Route      | Service            | Description       |
| ---------- | ------------------ | ----------------- |
| `/graphql` | federation-service | GraphQL API       |
| `/ws`      | delivery-service   | WebSocket service |
| `/search`  | elasticsearch      | Search API        |

---

## Running locally

Build and start Traefik:

```bash
docker build -t fitrang-traefik .
docker run -p 8000:8000 fitrang-traefik
```

Or using docker compose:

```bash
docker compose up --build
```

---

## Headers added to upstream services

After successful authentication:

```
X-User-Email
```

These headers should be trusted only when coming from Traefik.

---

## Notes

* JWKS keys are cached in memory and automatically refreshed
* No private credentials are stored in this gateway
* All authentication is stateless and verified locally

---

## Purpose

This gateway centralizes routing, authentication, and middleware logic, ensuring backend services remain simple and focused on business logic.

