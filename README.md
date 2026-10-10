# go-commerce-starter

Open-source e-commerce backend built with Go. A modular starter kit / base engine for storefront APIs — not a turn-key enterprise platform.

Use it to learn, prototype, or as a foundation you extend for your own shop.

[![Go](https://img.shields.io/badge/Go-1.17+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![MongoDB](https://img.shields.io/badge/MongoDB-4.4+-47A248?logo=mongodb&logoColor=white)](https://www.mongodb.com/)
[![Gin](https://img.shields.io/badge/Gin-HTTP-08A4E0)](https://gin-gonic.com/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![CI](https://github.com/abolfazlnorozian/go-commerce-starter/actions/workflows/go.yml/badge.svg)](https://github.com/abolfazlnorozian/go-commerce-starter/actions/workflows/go.yml)

## Tech stack

| Layer | Choice |
| --- | --- |
| Language | Go 1.17+ (Dockerfile uses 1.18; CI uses 1.20) |
| HTTP | [Gin](https://gin-gonic.com/) |
| Database | MongoDB via [mongo-driver](https://pkg.go.dev/go.mongodb.org/mongo-driver) (no separate ORM) |
| Auth | JWT (`dgrijalva/jwt-go`) + bcrypt |
| Payments | [Zarinpal](https://www.zarinpal.com/) gateway helper |
| SMS | IPPanel (OTP / order notifications) |
| Docs | Swagger / OpenAPI (`swaggo`) |
| Containers | Docker + Docker Compose |

## Core features

Key features currently supported out of the box:

**Storefront (`/api`)**
- SMS-based user registration / login and JWT sessions
- Products (list, category filter, slug), categories, brands
- Cart, mix products, coupons, favorites
- Orders + Zarinpal checkout / payment verify callback
- Addresses, product comments, CMS-style pages (`/pages/index`)
- Image serving with on-the-fly WebP conversion
- Province/city JSON endpoint (`/json/state.json` — requires local `state.json`)

**Admin (`/admin/api`)**
- Admin login
- CRUD-style management for products (dimensions / variations), categories, brands, properties
- Orders (list, detail, status updates), coupons, comments moderation
- Users listing / detail, image uploads, pages, broadcast SMS to users

**Ops / DX**
- Swagger UI at `/swagger/index.html`
- Docker Compose for app + MongoDB
- Environment-driven config via `.env`

## Architecture overview

The project uses a **modular layout**: HTTP routers wire into service handlers; entities model MongoDB documents; shared helpers cover auth tokens, SMS, payments, and IDs.

```
  Client / Admin UI
          |
          v
  +------------------+     +------------------+
  | Gin routers      |---->| Services         |
  | /api, /admin/api |     | user/ + admin/   |
  +--------+---------+     +--------+---------+
           |                        |
           v                        v
  +------------------+     +------------------+
  | Auth (JWT)       |     | MongoDB          |
  +------------------+     +------------------+
           |
           v
  +------------------+
  | Helpers          |
  | Zarinpal / SMS   |
  +------------------+
```

```
.
├── main.go              # Entry point, CORS, route registration, Swagger
├── router/              # Public / storefront routes
├── adminroute/          # Admin API routes
├── user/services/       # Storefront business handlers
├── admin/service/       # Admin business handlers
├── entities/            # Domain structs (product, order, user, …)
├── auth/                # JWT issue / validate + middleware
├── database/            # Mongo client + env helpers
├── helpers/             # SMS, Zarinpal, IDs, hashing
├── upload/              # Upload helpers
├── response/            # Shared response shapes
├── docs/                # Generated Swagger artifacts
├── assets/              # Payment result HTML / CSS
├── public/images/       # Uploaded media (gitignored content)
├── Dockerfile
└── docker-compose.yml
```

## Local setup & getting started

### Prerequisites

- Go **1.17+** (1.20 recommended to match CI)
- MongoDB **4.4+** (or use Compose)
- Docker & Docker Compose (optional)
- A `.env` file (see below)
- Optional: Zarinpal merchant ID, IPPanel API key, and a local `state.json` for province/city data

### Environment variables

```bash
cp .env.example .env
```

Edit `.env` and set at least `MONGO_URL`, `DBNAME`, `SECRET_KEY`, `CORS_DOMAIN`, and `PORT`. See [`.env.example`](./.env.example) for the full list and comments.


### Database

There is **no formal migration tool**. Collections are created implicitly when documents are inserted. Seed data (admin users, catalog, etc.) is application-specific — create records via the admin API or insert documents in MongoDB directly.



### Run locally

```bash
# Ensure MongoDB is reachable at MONGO_URL
go run main.go
```

Default listen port is `8000` unless `PORT` is set.

Swagger: `http://localhost:8000/swagger/index.html`

### Run with Docker Compose

```bash
cp .env.example .env
# fill in .env — for Compose, point MONGO_URL at the mongo service, e.g.:
# MONGO_URL=mongodb://mongo:27017
docker-compose up -d --build
```

Compose exposes MongoDB on host port `27018` → container `27017`.

### Tests

Automated unit and integration test suites are planned on the roadmap. Contributions for test coverage are welcome!
```bash
# Run tests once implemented
go test ./...
```

### Build

```bash
go build -o bin/go-commerce-starter .
```

## Known limitations & roadmap

Honest assessment of gaps — this is a starter engine, not production-hardened SaaS:

| Area | Status |
| --- | --- |
| Automated tests | Minimal / missing — high priority |
| Admin route auth consistency | Not all admin handlers are uniformly wrapped with auth middleware; audit before exposing publicly |
| Rate limiting / abuse protection | Not implemented |
| Payment webhook idempotency | Verify callback should be hardened for retries / replay |
| Distributed tracing / metrics | Not implemented |
| Structured logging | Mostly `fmt` / `log` |
| Formal DB migrations / seed CLI | Not implemented |
| Dependency hygiene | Uses `dgrijalva/jwt-go` (unmaintained); consider `golang-jwt/jwt` |
| SMS pattern codes | Must be set via `IPPANEL_PATTERN_*` (no private defaults); missing config skips send instead of failing hard |
| `state.json` | Local file for Iran province/city data; not committed — supply your own |
| Multi-tenant / i18n / inventory reservations | Out of scope for this base kit |

**Suggested roadmap:** auth audit on admin routes → tests for cart/order/payment → migrate JWT library → rate limits → idempotent payment verify → OpenTelemetry or equivalent.

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) for branch naming, PR expectations, and commit style.

## Security

Please report vulnerabilities privately as described in [SECURITY.md](./SECURITY.md). Do not open public issues for sensitive disclosures.

## License

MIT — see [LICENSE](./LICENSE).
