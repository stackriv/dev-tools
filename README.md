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

| Tool                          | Description                                                                             |
|-------------------------------|-----------------------------------------------------------------------------------------|
| 🚫 **.gitignore Generator**   | Generate a `.gitignore` file by selecting your languages, frameworks and tools          |
| 📄 **License Generator**      | Generate an open source license (MIT, Apache 2.0, GPL v3, BSD, ISC, MPL 2.0)            |
| 📦 **package.json Generator** | Generate a `package.json` file for your Node.js project                                 |
| 🔐 **.env Generator**         | Generate a `.env` template with presets for common stacks                               |
| 🐳 **Dockerfile Generator**   | Generate a production-ready Dockerfile with multi-stage build and non-root user support |

## Features

- ✅ Syntax highlighting on all outputs (highlight.js)
- ✅ One-click copy to clipboard
- ✅ Clean sidebar navigation
- ✅ Multi-arch support (`linux/amd64`, `linux/arm64`)
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
    restart: unless-stopped
```

## Supported Technologies

### .gitignore Generator
Languages: `Go`, `Python`, `Node`, `Java`, `Rust`, `C++`
Frameworks: `Angular`, `React`, `Next.js`, `Vue`, `Laravel`, `Django`
Tools: `Docker`, `Terraform`, `Linux`, `macOS`, `Windows`, `VSCode`, `JetBrains`

### License Generator
`MIT`, `Apache 2.0`, `GPL v3`, `BSD 2-Clause`, `BSD 3-Clause`, `ISC`, `Mozilla Public License 2.0`

### .env Generator
Presets: `Node.js`, `Go`, `Django`, `Database`, `JWT Auth`, `SMTP Mail`, `S3 / MinIO`, `Docker`

### Dockerfile Generator
Languages: `Go`, `Node.js`, `Python`, `Java`, `Rust`, `PHP`
Options: multi-stage build, non-root user, healthcheck

## API Reference

All generators expose a REST API endpoint.

### .gitignore
```http
POST /api/gitignore
Content-Type: application/json

{ "technologies": ["Go", "Docker", "VSCode"] }
```

### License
```http
POST /api/license
Content-Type: application/json

{ "type": "mit", "author": "John Doe", "year": "2026", "project": "my-project" }
```

### package.json
```http
POST /api/packagejson
Content-Type: application/json

{
  "name": "my-package",
  "version": "1.0.0",
  "description": "...",
  "author": "John Doe",
  "license": "MIT",
  "type": "module",
  "includeScripts": true,
  "includeEngines": false,
  "includePrivate": false
}
```

### .env
```http
POST /api/env
Content-Type: application/json

{ "presets": ["nodejs", "database", "jwt"], "appName": "my-app" }
```

### Dockerfile
```http
POST /api/dockerfile
Content-Type: application/json

{
  "lang": "go",
  "langVersion": "1.23",
  "port": "8080",
  "workdir": "/app",
  "multistage": true,
  "nonroot": true,
  "healthcheck": false
}
```

## Versioning

| Tag      | Description                  |
|----------|------------------------------|
| `latest` | Always the most recent build |
| `1.0`    | Latest patch of 1.0.x        |

## Built With

- [Go](https://golang.org/) — Backend + HTML templating
- [highlight.js](https://highlightjs.org/) — Syntax highlighting
- [Alpine Linux](https://alpinelinux.org/) — Base Docker image

## Part of Stackriv

Dev Tools is part of the **Stackriv** homelab infrastructure suite — a collection of self-hosted tools for managing personal and professional servers.

---

Made with ❤️ by [Stackriv](https://stackriv.dev)