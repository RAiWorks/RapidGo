# RGo Framework

A modern Go web framework inspired by Laravel and CodeIgniter — combining
**Laravel-style developer experience** with **Go performance**.

## Overview

RGo is an opinionated, full-stack Go web framework built on proven
libraries (Gin, GORM, Cobra) that provides everything needed for web
applications, REST APIs, and real-time WebSocket services.

### Key Features

- **MVC + Services + Helpers** architecture
- **Service Container & Providers** — Laravel-style IoC for extensibility
- **Multi-database support** — PostgreSQL, MySQL, SQLite via GORM
- **Session management** — DB, Redis, File, Memory, Cookie backends
- **Built-in validation** — zero-dependency validator + struct-based (go-playground)
- **Built-in crypto** — AES-256-GCM, HMAC-SHA256, secure random tokens
- **JWT & session-based authentication**
- **CLI scaffolding** — `make:controller`, `make:model`, `make:service`, etc.
- **Middleware registry** — aliases, groups, custom middleware
- **WebSocket support** — via `coder/websocket`
- **Caching** — Redis and in-memory backends
- **Mail** — SMTP via `go-mail`
- **File storage** — local disk and S3
- **Events / hooks** — pub-sub event dispatcher
- **i18n / localization** — JSON translation files
- **Caddy integration** — embedded or reverse proxy (optional)
- **Docker support** — multi-stage builds (optional)
- **Graceful shutdown**, health checks, CSRF, CORS, rate limiting

## Project Status

**Phase:** Architecture & Documentation

The framework is currently in the architecture design phase.
Architecture blueprint and documentation strategy are maintained
locally in `reference/docs/` (not tracked in version control).

## Repository

- **GitHub:** [RAiWorks/RGo](https://github.com/RAiWorks/RGo)

## Tech Stack

| Component | Library |
|-----------|---------|
| HTTP Router | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io) |
| CLI | [Cobra](https://github.com/spf13/cobra) |
| Config | [godotenv](https://github.com/joho/godotenv) / [Viper](https://github.com/spf13/viper) |
| JWT | [golang-jwt](https://github.com/golang-jwt/jwt) |
| WebSocket | [coder/websocket](https://github.com/coder/websocket) |
| Redis | [go-redis](https://github.com/redis/go-redis) |
| Mail | [go-mail](https://github.com/wneessen/go-mail) |
| S3 | [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2) |

## License

TBD
