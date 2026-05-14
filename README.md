# Way Framework

![Way Mascot](./way_mascot.png)

[![Go Reference](https://pkg.go.dev/badge/github.com/swayedev/way.svg)](https://pkg.go.dev/github.com/swayedev/way)

_Version: 1.0.0-rc1_

**Status:** Production-ready release candidate. Way v1.0.0-rc1 is the first stable API release. See [API_FREEZE.md](API_FREEZE.md) for the frozen public API surface.

## Overview

Way is a lightweight, public Go web framework in the Echo/Gin space. It integrates with the [Gorilla Mux](https://github.com/gorilla/mux) router and provides simple request handling with built-in database helpers, session management, and cryptographic wrappers.

Way is designed as a small, efficient foundation for services that need framework ergonomics plus secure defaults. Its crypto helpers delegate to [fcrypt](https://github.com/swayedev/fcrypt) `v1.0.0-rc1` and preserve Way's legacy hex string ciphertext format for compatibility.

## Documentation

- [API Freeze](API_FREEZE.md) – frozen public APIs and stability guarantees
- [Security Policy](SECURITY.md) – security considerations and hardening guidance  
- [Migration Guide](MIGRATION.md) – upgrade from pre-1.0 versions
- [Production Checklist](PRODUCTION_CHECKLIST.md) – release and deployment checks
- [Changelog](CHANGELOG.md) – release history

## Features

- Custom context for HTTP handlers with response helpers (JSON, XML, HTML, String, Data, Images)
- Simplified route declaration (GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD)
- Integrated SQL database operations with optional driver adapter packages for MySQL, PostgreSQL (pgx), SQLite, SQL Server, and Oracle
- Session and cookie management with Gorilla Sessions
- Graceful shutdown and startup management with safe HTTP server timeouts
- Encryption and hashing via [fcrypt](https://github.com/swayedev/fcrypt) integration
- Request logging middleware with method, path, and duration tracking

## Getting Started

### Installation
To start using the Way framework, install it using `go get`:
```bash
go get -u github.com/swayedev/way@v1.0.0-rc1
```

### Basic Usage
Here's a simple example to get you started:

```go
package main

import (
    "github.com/swayedev/way"
)

func main() {
    w := way.New()

    w.GET("/", func(ctx *way.Context) {
        ctx.Response.Write([]byte("Hello, World!"))
    })

    if err := w.Start(":8080"); err != nil {
        panic(err)
    }
}
```

Or with a handler function

```go
package main

import "github.com/swayedev/way"

func main() {
	w := way.New()
	w.GET("/", helloHandler)

	if err := w.Start(":8080"); err != nil {
		panic(err)
	}
}

func helloHandler(c *way.Context) {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.Write([]byte("Hello World"))
}

```

## Routing
Way simplifies route handling with predefined methods for standard HTTP verbs:

```go
w.GET("/path", yourGetHandler)
w.POST("/path", yourPostHandler)
// ... and so on for PUT, DELETE, PATCH, OPTIONS, HEAD
```

## Database Operations

Way keeps database drivers out of the core package. Import the driver adapter your application needs:
```go
import (
	"github.com/swayedev/way"

	_ "github.com/swayedev/way/database/drivers/sqlite"
)
```

Available adapter packages:

- `github.com/swayedev/way/database/drivers/mysql`
- `github.com/swayedev/way/database/drivers/pgx`
- `github.com/swayedev/way/database/drivers/sqlite`
- `github.com/swayedev/way/database/drivers/sqlserver`
- `github.com/swayedev/way/database/drivers/godror`

### Opening a Connection
```go
db := way.NewDB()
err := db.SQLOpen("sqlite3", "app.db")
if err != nil {
    // Handle error
}
w.SetDB(&db)
```

### Executing Queries
```go
result, err := w.SqlExec(context.Background(), "your SQL query here", args...)
if err != nil {
    // Handle error
}
```

### Querying Data
```go
rows, err := w.SqlQuery(context.Background(), "your SQL query here", args...)
if err != nil {
    // Handle error
}
// Remember to close rows
```

## Crypto

Way's `crypto.Encrypt`, `crypto.Decrypt`, and matching `Context` helpers delegate to fcrypt `v1.0.0-rc1`:

```go
ciphertext, err := crypto.Encrypt([]byte("secret"), "passphrase")
plaintext, err := crypto.Decrypt(ciphertext, "passphrase")
```

`Encrypt` returns a hex string for compatibility with existing Way users. New applications that need lower-level key management or key rotation should use fcrypt directly.

## Graceful Shutdown
To gracefully shut down your server:

```go
err := w.Shutdown(context.Background())
if err != nil {
    // Handle error
}
```

## Contributing
Contributions to the Way framework are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License
This project is licensed under the [MIT License](LICENSE).
