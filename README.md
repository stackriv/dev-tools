# ⚙️ Stackriv Dev Tools

> A lightweight self-hosted suite of developer tools built with Go — part of the Stackriv infrastructure suite.

![Docker Pulls](https://img.shields.io/docker/pulls/proverbes12x/dev-tools)
![Docker Image Size](https://img.shields.io/docker/image-size/proverbes12x/dev-tools/latest)
![Go Version](https://img.shields.io/badge/Go-1.25-blue)
![License](https://img.shields.io/badge/license-MIT-green)

## Overview

**Stackriv Dev Tools** is a self-hosted web application that provides a collection of generators to help you scaffold your projects faster. It features a clean dark-themed interface with syntax highlighting on all outputs.

No external dependencies, no telemetry, no cloud — just a single Go binary running on your server.

## Tools

### 🔧 Generators
| Tool                   | Description                                                                   |
|------------------------|-------------------------------------------------------------------------------|
| 🚫 **.gitignore**      | Generate a `.gitignore` for your language, framework and tools                |
| 📄 **License**         | Generate MIT, Apache 2.0, GPL v3, BSD, ISC or MPL 2.0 license                 |
| 📦 **package.json**    | Generate a `package.json` for your Node.js project                            |
| 🔐 **.env**            | Generate a `.env` template with presets (DB, JWT, SMTP, S3, Docker...)        |
| 🐳 **Dockerfile**      | Generate a production-ready Dockerfile (multi-stage, non-root, healthcheck)   |
| 📝 **README**          | Generate a professional `README.md` with badges and sections                  |
| 🗂️ **Docker Compose** | Generate a `compose.yaml` for your stack (Postgres, Redis, Caddy, Grafana...) |

### 🛠️ Utilities
| Tool                   | Description                                                           |
|------------------------|-----------------------------------------------------------------------|
| 🔑 **UUID Generator**  | Generate random UUIDs v4 — single or bulk, with format options        |
| 🔄 **Base64**          | Encode or decode Base64 strings (standard and URL-safe)               |
| #️⃣ **Hash Generator** | Generate MD5, SHA1, SHA256 and SHA512 hashes                          |
| 🔓 **JWT Decoder**     | Decode and inspect a JWT token without any external service           |
| 🔍 **Regex Tester**    | Test regular expressions with match highlighting and group capture    |
| ✍️ **Markdown Editor** | Write Markdown with live split-screen preview and syntax highlighting |

### 📡 Infrastructure
| Tool                  | Description                                                       |
|-----------------------|-------------------------------------------------------------------|
| ⏰ **Cron Builder**    | Build and validate cron expressions with next 5 run preview       |
| 🌐 **DNS Lookup**     | Resolve A, AAAA, MX, TXT, NS and CNAME records for any domain     |
| 🔒 **SSL Checker**    | Check certificate validity, issuer, expiration and TLS version    |
| 📡 **Uptime Monitor** | Check availability and response time of your services in parallel |
| 📋 **Log Viewer**     | View and filter logs from your Docker containers                  |

### 🧾 Business
| Tool                     | Description                                                |
|--------------------------|------------------------------------------------------------|
| 🧾 **Invoice Generator** | Generate professional HTML invoices with PDF print support |

## Features

- ✅ Syntax highlighting on all code outputs (highlight.js)
- ✅ One-click copy to clipboard on every tool
- ✅ Sidebar navigation with scrollable menu
- ✅ Split-screen Markdown editor with live preview
- ✅ Invoice generator with PDF print / save support
- ✅ Multi-arch Docker image (`linux/amd64`, `linux/arm64`)
- ✅ Tiny image size (~15 MB)
- ✅ No external API calls — everything runs locally
- ✅ Self-hosted — your data never leaves your server

## Quick Start

```bash
docker run -d \
  --name dev-tools \
  -p 8090:8090 \
  proverbes12x/dev-tools:latest
```

Then open `http://localhost:8090` in your browser.

## Docker Compose

```yaml
services:
  dev-tools:
    image: proverbes12x/dev-tools:latest
    container_name: dev-tools
    ports:
      - "127.0.0.1:8090:8090"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock  # required for Log Viewer
    restart: unless-stopped
```

> **Note:** The Docker socket mount is only required for the **Log Viewer** tool. Remove it if you don't need that feature.

## Supported Technologies

### .gitignore Generator
**Languages:** Go, Python, Node, Java, Rust, C++
**Frameworks:** Angular, React, Next.js, Vue, Laravel, Django
**Tools:** Docker, Terraform, Linux, macOS, Windows, VSCode, JetBrains

### License Generator
MIT, Apache 2.0, GPL v3, BSD 2-Clause, BSD 3-Clause, ISC, Mozilla Public License 2.0

### .env Generator
Presets: `Node.js`, `Go`, `Django`, `Database`, `JWT Auth`, `SMTP Mail`, `S3 / MinIO`, `Docker`

### Dockerfile Generator
**Languages:** Go, Node.js, Python, Java, Rust, PHP
**Options:** multi-stage build, non-root user, healthcheck

### Docker Compose Generator
PostgreSQL, MySQL, Redis, MongoDB, Nginx, Caddy, MinIO, RabbitMQ, Grafana, Prometheus, Portainer

### Invoice Generator
**Currencies:** USD, EUR, GBP, XOF (FCFA), CAD
**Features:** multiple line items, tax rate, notes, PDF print/save

## API Reference

All tools expose a REST API endpoint.

```http
POST /api/gitignore
{ "technologies": ["Go", "Docker", "VSCode"] }
 
POST /api/license
{ "type": "mit", "author": "John Doe", "year": "2026" }
 
POST /api/packagejson
{ "name": "my-app", "version": "1.0.0", "license": "MIT", "includeScripts": true }
 
POST /api/env
{ "presets": ["nodejs", "database", "jwt"], "appName": "my-app" }
 
POST /api/dockerfile
{ "lang": "go", "langVersion": "1.23", "port": "8080", "multistage": true, "nonroot": true }
 
POST /api/readme
{ "name": "my-project", "author": "John", "github": "pro12x/my-project", "sections": ["overview", "features", "installation"] }
 
POST /api/compose
{ "services": [{"name": "postgres"}, {"name": "redis"}], "network": "mynet" }
```

### Utilities

```http
GET  /api/uuid?count=5
 
POST /api/base64
{ "input": "Hello World", "action": "encode" }
 
POST /api/hash
{ "input": "my secret" }
 
POST /api/jwt
{ "token": "eyJ..." }
 
POST /api/regex
{ "pattern": "[a-z]+", "input": "hello world", "flags": { "global": true } }
 
POST /api/markdown
{ "content": "# Hello\n\nThis is **markdown**." }
```

### Infrastructure

```http
POST /api/cron
{ "expression": "0 3 1 * *" }
 
POST /api/dns
{ "domain": "example.com", "types": ["A", "MX", "TXT"] }
 
POST /api/ssl
{ "domain": "example.com" }
 
POST /api/uptime
{ "targets": [{ "name": "My App", "url": "https://example.com" }] }
 
POST /api/logs
{ "container": "caddy", "lines": 100, "search": "error", "since": "1h" }
 
GET  /api/containers
```

### Business

```http
POST /api/invoice
{
  "number": "INV-2026-001",
  "date": "2026-06-01",
  "currency": "EUR",
  "from": { "name": "ZENDEV Labs", "email": "contact@zendev.com" },
  "to": { "name": "Client Corp", "email": "client@corp.com" },
  "items": [{ "description": "Development", "quantity": 10, "unitPrice": 150 }],
  "taxRate": 20
}
```

## Versioning

| Tag      | Description                           |
|----------|---------------------------------------|
| `latest` | Always the most recent build          |
| `1.0.0`  | Specific realease - never overwritten |

## Built With

- [Go](https://golang.org/) — Backend + HTML templating
- [highlight.js](https://highlightjs.org/) — Syntax highlighting
- [Alpine Linux](https://alpinelinux.org/) — Base Docker image

## Part of Stackriv

Dev Tools is part of the **Stackriv** homelab infrastructure suite — a collection of self-hosted tools for managing personal and professional servers.

---

Made with ❤️ by [Stackriv](https://stackriv.dev)