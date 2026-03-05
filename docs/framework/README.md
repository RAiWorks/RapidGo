# RGo Framework — Architecture Documentation

> Comprehensive architecture documentation for the RGo Go web
> framework, following RFC-style standards.

## Quick Links

- [Getting Started](guides/getting-started.md)
- [Architecture Overview](architecture/overview.md)
- [Project Structure](architecture/project-structure.md)

---

## Documentation Map

### Architecture

| Document | Description |
|----------|-------------|
| [Overview](architecture/overview.md) | Mission, capabilities, tech stack |
| [Design Principles](architecture/design-principles.md) | Laravel-inspired DX, conventions |
| [Project Structure](architecture/project-structure.md) | Directory tree, framework vs user code |
| [Application Lifecycle](architecture/application-lifecycle.md) | Boot sequence, graceful shutdown |

#### Diagrams

| Document | Description |
|----------|-------------|
| [System Overview](architecture/diagrams/system-overview.md) | 7-layer architecture diagram |
| [Request Lifecycle](architecture/diagrams/request-lifecycle.md) | Request → response flow |
| [Service Container](architecture/diagrams/service-container.md) | DI container architecture |
| [Data Flow](architecture/diagrams/data-flow.md) | Data flow through layers |

### Core

| Document | Description |
|----------|-------------|
| [Service Container](core/service-container.md) | Dependency injection container |
| [Service Providers](core/service-providers.md) | Provider pattern, built-in providers |
| [Configuration](core/configuration.md) | `.env` loading, environment detection |
| [Error Handling](core/error-handling.md) | Error middleware, JSON/HTML responses |
| [Logging](core/logging.md) | Structured logging with `log/slog` |

### HTTP

| Document | Description |
|----------|-------------|
| [Routing](http/routing.md) | Routes, groups, resources, named routes |
| [Controllers](http/controllers.md) | MVC controllers, ResourceController |
| [Views](http/views.md) | Templates, static files |
| [Middleware](http/middleware.md) | Registry, aliases, groups |
| [Requests & Validation](http/requests-validation.md) | Built-in and struct-based validation |
| [Responses](http/responses.md) | API response helpers |
| [WebSocket](http/websocket.md) | Real-time connections |

### Data

| Document | Description |
|----------|-------------|
| [Database](data/database.md) | Connection, drivers, pooling |
| [Models](data/models.md) | GORM models, relationships, hooks |
| [Migrations](data/migrations.md) | Schema management |
| [Seeders](data/seeders.md) | Database seeding |
| [Transactions](data/transactions.md) | GORM transaction patterns |
| [Pagination](data/pagination.md) | Paginate helper |

### Security

| Document | Description |
|----------|-------------|
| [Authentication](security/authentication.md) | JWT and session auth |
| [Sessions](security/sessions.md) | Multi-driver sessions, flash messages |
| [CSRF](security/csrf.md) | CSRF protection |
| [CORS](security/cors.md) | Cross-origin configuration |
| [Rate Limiting](security/rate-limiting.md) | Request rate limiting |
| [Crypto](security/crypto.md) | Encryption, hashing, HMAC |
| [Request ID](security/request-id.md) | Unique request identifiers |

### Infrastructure

| Document | Description |
|----------|-------------|
| [Services Layer](infrastructure/services-layer.md) | Business logic services |
| [Caching](infrastructure/caching.md) | Redis and memory cache |
| [Mail](infrastructure/mail.md) | SMTP email sending |
| [File Storage](infrastructure/file-storage.md) | Local and S3 storage |
| [Events](infrastructure/events.md) | Publish-subscribe system |
| [i18n](infrastructure/i18n.md) | Internationalization |

### CLI

| Document | Description |
|----------|-------------|
| [CLI Overview](cli/cli-overview.md) | Available commands |
| [Code Generation](cli/code-generation.md) | Scaffolding templates |

### Guides

| Document | Description |
|----------|-------------|
| [Getting Started](guides/getting-started.md) | First project setup |
| [Creating a CRUD App](guides/creating-crud.md) | Full CRUD walkthrough |
| [Building a REST API](guides/building-api.md) | JSON API with JWT |
| [SSR Forms](guides/ssr-forms.md) | Forms with validation & flash |
| [Custom Service & Provider](guides/custom-service.md) | Extending with DI |
| [Extending the Framework](guides/extending-framework.md) | Customization patterns |

### Testing

| Document | Description |
|----------|-------------|
| [Testing Overview](testing/testing-overview.md) | Strategy and tools |
| [Unit Tests](testing/unit-tests.md) | Service and helper tests |
| [Integration Tests](testing/integration-tests.md) | HTTP handler tests |

### Deployment

| Document | Description |
|----------|-------------|
| [Build and Run](deployment/build-and-run.md) | Entrypoint and graceful shutdown |
| [Caddy Integration](deployment/caddy.md) | Auto-HTTPS reverse proxy |
| [Docker](deployment/docker.md) | Containerization |
| [Health Checks](deployment/health-checks.md) | Liveness and readiness probes |

### Reference

| Document | Description |
|----------|-------------|
| [Environment Variables](reference/env-reference.md) | All `.env` variables |
| [Helpers](reference/helpers-reference.md) | Utility functions |
| [Middleware](reference/middleware-reference.md) | Middleware quick reference |
| [Libraries](reference/libraries.md) | External dependencies |

### Appendix

| Document | Description |
|----------|-------------|
| [Glossary](appendix/glossary.md) | Term definitions |
| [Roadmap](appendix/roadmap.md) | Planned features and timeline |
| [Naming](appendix/naming.md) | Framework naming |

---

## Document Standards

All documents follow RFC-style conventions:

- **YAML frontmatter** with title, version, status, date, authors
- **Abstract** summary
- **Terminology** section with RFC 2119 normative language
- **Security Considerations** where applicable
- **Cross-references** to related documents
- **Revision History** table

See the [Architecture Docs Strategy](../reference/docs/architecture_docs_strategy.md)
for the full documentation plan.
