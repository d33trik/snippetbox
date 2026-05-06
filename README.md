# Snippetbox

## Description

Snippetbox is a simple web application for sharing text snippets (code fragments, notes, quotes, etc.).

Users can create new snippets with custom titles and expiration times, and view public snippets sorted by creation date.

This project is intended for **learning and practice purposes only**. It was built step-by-step alongside the *[Let's Go](https://lets-go.alexedwards.net/)* book by [Alex Edwards](https://www.alexedwards.net/) to learn core Go web development concepts.

Code quality and patterns prioritize educational clarity over production-grade standards.

## Key Concepts Covered

This project covers core web development concepts in Go, including:

- HTTP server with custom handlers, routing, and static file serving
- Idiomatic project structure and dependency management with Go modules
- Configuration via command-line flags and dependency injection
- Middleware patterns: request logging, panic recovery, security headers, CSRF protection, and session management
- Database design, connection pooling, and safe queries using `database/sql`
- Template rendering, caching, custom functions, and dynamic data handling
- Form validation and user-friendly error reporting
- Session-based authentication and user management
- HTTPS configuration with TLS and secure defaults
- Embedding static assets into the Go binary using `//go:embed`
- Web security fundamentals: SQL injection, CSRF, XSS, clickjacking, and slow-client protection
- Context propagation using `context.Context`
- Testing strategies: unit tests (with mocks), integration tests (with MySQL), and end-to-end tests

## Technologies Used

- **Go 1.26.1**
- **Standard library**: `net/http`, `html/template`, `database/sql`, `embed`, `crypto/tls`, etc.
- **Third-party packages**:
  - MySQL driver: `github.com/go-sql-driver/mysql`
  - Session management: `github.com/alexedwards/scs/v2`
  - Middleware chaining: `github.com/justinas/alice`
  - Form decoding: `github.com/go-playground/form/v4`
  - CSRF protection: `github.com/justinas/nosurf`
- **Tooling**:
  - Podman (MySQL 9.6 container)
  - `modd` (live reload for development)
  - Prettier (HTML template formatting with `prettier-plugin-go-template`)

## Project Structure

```
snippetbox/
├── cmd/              # Application entrypoints
│   └── web/          # Web server, handlers, middleware, routes, templates
├── internal/         # Internal packages
│   ├── assert/       # Test assertion helpers
│   ├── models/       # Database models, integration tests, mocks
│   └── validator/    # Form validation helpers
├── scripts/          # Helper scripts (dev runner, package checker)
├── sql/              # Database bootstrap and migration scripts
├── ui/               # HTML templates and static files (embedded in binary)
│   ├── html/         # Go templates
│   └── static/       # Static assets (CSS, JS, images)
├── compose.yaml      # Podman/Docker configuration for MySQL
├── go.mod            # Defines module path and project dependencies
├── go.sum            # Checksums for Go module dependencies (auto-managed)
├── modd.conf         # Live reload configuration for modd
├── package-lock.json # Lockfile for frontend/tooling dependencies
└── package.json      # Defines frontend/tooling dependencies and scripts
```

## How to Run the Project

### Prerequisites

- Go 1.26+
- Podman (for MySQL container)
- TLS certificates

### Steps

#### Clone the repository and enter the project directory

```sh
git clone https://codeberg.org/d33trik/snippetbox.git && cd snippetbox
```

#### Generate TLS certificates (if not already present)

##### Option A: Using Go's built-in tool

```sh
mkdir -p tls && cd tls
```

```sh
go run "$(go env GOROOT)/src/crypto/tls/generate_cert.go" -rsa-bits=2048 --host=localhost
```

##### Option B: Using `mkcert`

Install `mkcert`: <https://github.com/FiloSottile/mkcert#installation>

```sh
mkdir -p tls && cd tls
```

```sh
mkcert -install
```

```sh
mkcert localhost
```

```sh
mv localhost-key.pem key.pem  
mv localhost.pem cert.pem
```

##### Differences

- `generate_cert.go`
  - Built into Go (no external dependencies)
  - Generates self-signed certificates
  - Browsers will show security warnings
- `mkcert`
  - Requires installation
  - Creates locally trusted certificates
  - No browser warnings

#### Start the database

```sh
podman compose up -d
```

#### Build and run the application

```sh
go build -o web ./cmd/web
```

```sh
./web
```

The application will be available at <https://localhost:4000>

## Development

This section covers tools and workflows used during development. These are **not required** to run the application.

### Prerequisites

- Node.js + npm (for frontend tooling)
- Podman (for development database)
- `modd` (for live reload)

### Steps

#### Install frontend dependencies

```sh
npm install
```

Installs development dependencies used for formatting:

- `prettier`
- `prettier-plugin-go-template`

#### Run in development mode (auto-reload on changes)

```sh
./scripts/run.sh
```

This script:

- Starts the MySQL container (if not already running)
- Waits for the database to be ready
- Runs `modd` to rebuild and restart the app on file changes

#### Check formatting

Verifies if files are properly formatted.

```sh
npm run format:check
```

#### Format files

Formats all project files using Prettier with the Go template plugin.

```sh
npm run format:write
```

## Running Tests

> Requires MySQL running

### Run all tests

```sh
go test ./...
```

## Observations

- This project is **not intended for production use**.
- Some best practices and production-level concerns are intentionally simplified to match the learning goals and the structure presented in the book.
- The test suite requires a running MySQL instance.
- HTML templates are formatted using Prettier with the `go-template` parser.

## Credits

- Thanks to [Alex Edwards](https://www.alexedwards.net/) for the comprehensive guide to Go web development.
- Thanks to third-party package authors for the libraries used in this project.
